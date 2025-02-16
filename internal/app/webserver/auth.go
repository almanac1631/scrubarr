package webserver

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
	"log/slog"
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
