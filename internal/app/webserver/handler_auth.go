package webserver

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/almanac1631/scrubarr/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

const (
	sessionCookieName = "session"
	sessionExpiryTime = time.Hour * 12
)

func (handler *handler) handleLogin(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		username := request.PostFormValue("username")
		if username == "" {
			http.Error(writer, "username is required", http.StatusBadRequest)
			return
		}
		slog.Debug("Incoming authentication request", "username", username)
		password := []byte(request.PostFormValue("password"))
		if len(password) == 0 {
			http.Error(writer, "password is required", http.StatusBadRequest)
			return
		}
		passwordHashExpected := handler.passwordRetriever()
		incorrectUsername := handler.username != username
		if incorrectUsername || !checkPassword(passwordHashExpected, password, handler.passwordSalt) {
			writer.Header().Set("Content-Type", "text/html; charset=utf-8")
			writer.WriteHeader(http.StatusUnauthorized)
			if err := handler.templateCache["login.gohtml"].ExecuteTemplate(writer, "login_notification", nil); isErrAndNoBrokenPipe(err) {
				slog.Error("failed to execute template", "err", err)
				return
			}
			return
		}
		slog.Debug("Successfully authenticated user", "username", username)
		jwtStr, err := generateToken(handler.jwtConfig.PrivateKey, username)
		if err != nil {
			slog.Error("Error generating the JWT", "error", err)
			http.Error(writer, "internal server error", http.StatusInternalServerError)
		}
		http.SetCookie(writer, &http.Cookie{
			Name:     sessionCookieName,
			Value:    jwtStr,
			Quoted:   false,
			Expires:  time.Now().Add(sessionExpiryTime),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		writer.Header().Set("Hx-Redirect", "/media")
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	if utils.IsHTMXRequest(request) {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("404 Not Found"))
		return
	}
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := handler.templateCache["login.gohtml"].ExecuteTemplate(writer, "index", nil); isErrAndNoBrokenPipe(err) {
		slog.Error("failed to execute template", "err", err)
		return
	}
}

func (handler *handler) handleLogout(writer http.ResponseWriter, request *http.Request) {
	http.SetCookie(writer, &http.Cookie{Name: sessionCookieName, MaxAge: -1})
	writer.Header().Set("Hx-Redirect", "/login")
	writer.WriteHeader(http.StatusNoContent)
}

func generateToken(key *ecdsa.PrivateKey, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(sessionExpiryTime).Unix(),
	})
	return token.SignedString(key)
}

func checkPassword(passwordHashExpected, passwordRawActual, salt []byte) bool {
	passwordHashActual := generateHash(passwordRawActual, salt)
	return bytes.Equal(passwordHashActual, passwordHashExpected)
}

func generateHash(passwordRaw, salt []byte) []byte {
	result := argon2.IDKey(passwordRaw, salt, 2, 19456, 1, 16)
	return result
}

func validateToken(key *ecdsa.PublicKey, jwtStr string) (bool, error) {
	jwtToken, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodES256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return false, err
	}
	if !jwtToken.Valid {
		return false, nil
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	err = jwt.NewValidator(
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	).Validate(claims)
	return true, nil
}
