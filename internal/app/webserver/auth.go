package webserver

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"
)

func (a ApiEndpointHandler) Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error) {
	username := request.Body.Username
	slog.Debug("incoming login attempt", "username", username)
	passwordRawActual := []byte(request.Body.Password)
	passwordHashExpected := a.passwordRetriever()
	incorrectUsername := a.username != username
	if incorrectUsername || !checkPassword(passwordHashExpected, passwordRawActual, a.passwordSalt) {
		return Login401JSONResponse{"login failed", "username and password combination does not match"}, nil
	}
	slog.Debug("successful log in", "username", username)
	jwtStr, err := generateToken(a.jwtConfig.PrivateKey, username)
	if err != nil {
		return nil, err
	}
	return Login200JSONResponse{"Ok", jwtStr}, nil
}

const headerName = "Authorization"

func (a ApiEndpointHandler) AuthenticationMiddleware(routesWithoutAuth []string) MiddlewareFunc {
	sendUnauthorizedStatus := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusUnauthorized)
		bodyBytes, err := json.Marshal(ErrorResponseBody{
			Error:  http.StatusText(http.StatusUnauthorized),
			Detail: "the provided credentials are not valid",
		})
		if err != nil {
			slog.Error("error marshaling error body for unauthorized request", "error", err)
		}
		_, err = w.Write(bodyBytes)
		if err != nil {
			slog.Error("error writing error body for unauthorized request", "error", err)
		}
	}

	checkAuthorization := func(next http.Handler, w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(headerName)
		if authHeader == "" {
			sendUnauthorizedStatus(w)
			return
		}
		authHeaderSplit := strings.Split(authHeader, " ")
		if len(authHeaderSplit) != 2 || authHeaderSplit[0] != "Bearer" {
			sendUnauthorizedStatus(w)
			return
		}
		tokenOk, err := validateToken(a.jwtConfig.PublicKey, authHeaderSplit[1])
		if !tokenOk {
			if err != nil {
				slog.Debug("could not validate jwt token", "error", err, "token", authHeaderSplit[1])
			}
			sendUnauthorizedStatus(w)
			return
		}
		next.ServeHTTP(w, r)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if slices.Contains(routesWithoutAuth, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			checkAuthorization(next, w, r)
		})
	}
}

func generateToken(key *ecdsa.PrivateKey, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
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
