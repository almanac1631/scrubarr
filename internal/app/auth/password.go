package auth

import (
	"bytes"

	"golang.org/x/crypto/argon2"
)

type PasswordBasedProvider struct {
	username                   string
	passwordHash, passwordSalt []byte
}

func NewPasswordBasedProvider(username string, passwordHash, passwordSalt []byte) *PasswordBasedProvider {
	return &PasswordBasedProvider{username: username, passwordHash: passwordHash, passwordSalt: passwordSalt}
}

func (provider PasswordBasedProvider) CheckCredentials(username string, password []byte) (bool, error) {
	passwordHashExpected := provider.passwordHash
	incorrectUsername := provider.username != username
	if incorrectUsername || !checkPassword(passwordHashExpected, password, provider.passwordSalt) {
		return false, nil
	}
	return true, nil
}

func checkPassword(passwordHashExpected, passwordRawActual, salt []byte) bool {
	passwordHashActual := GenerateHash(passwordRawActual, salt)
	return bytes.Equal(passwordHashActual, passwordHashExpected)
}

func GenerateHash(passwordRaw, salt []byte) []byte {
	result := argon2.IDKey(passwordRaw, salt, 2, 19456, 1, 16)
	return result
}
