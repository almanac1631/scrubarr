package webserver

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"log/slog"
	"net/http"
	"path"
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
		logger := slog.With("remote", request.RemoteAddr)
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
			logger.Warn("Failed login attempt with incorrect credentials.", "username", username)
			writer.Header().Set("Content-Type", "text/html; charset=utf-8")
			writer.WriteHeader(http.StatusUnauthorized)
			if err := handler.ExecuteSubTemplate(writer, "login.gohtml", "login_notification", nil); err != nil {
				logger.Error(err.Error())
				return
			}
			return
		}
		logger.Info("Successful user login.", "username", username)
		jwtStr, err := generateToken(handler.jwtConfig.PrivateKey, username)
		if err != nil {
			logger.Error("Error generating the JWT", "error", err)
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
		writer.Header().Set("Hx-Redirect", path.Join(handler.pathPrefix, "/media"))
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	if utils.IsHTMXRequest(request) {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("404 Not Found"))
		return
	}
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := handler.ExecuteRootTemplate(writer, "login.gohtml", nil); err != nil {
		slog.Error(err.Error())
		return
	}
}

func (handler *handler) handleLogout(writer http.ResponseWriter, request *http.Request) {
	http.SetCookie(writer, &http.Cookie{Name: sessionCookieName, MaxAge: -1})
	writer.Header().Set("Hx-Redirect", path.Join(handler.pathPrefix, "/login"))
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
	passwordHashActual := GenerateHash(passwordRawActual, salt)
	return bytes.Equal(passwordHashActual, passwordHashExpected)
}

func GenerateHash(passwordRaw, salt []byte) []byte {
	result := argon2.IDKey(passwordRaw, salt, 2, 19456, 1, 16)
	return result
}

func validateToken(key *ecdsa.PublicKey, jwtStr string) (bool, string, error) {
	jwtToken, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodES256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return false, "", err
	}
	if !jwtToken.Valid {
		return false, "", nil
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	err = jwt.NewValidator(
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	).Validate(claims)
	username, err := claims.GetSubject()
	if err != nil {
		return false, "", err
	}
	return true, username, nil
}
