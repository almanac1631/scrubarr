package webserver

import (
	"crypto/ecdsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/knadh/koanf/v2"
)

type JwtConfig struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func loadJwtPrivateKey(config *koanf.Koanf) (*ecdsa.PrivateKey, error) {
	privateKeyPath := config.MustString("general.auth.jwt.private_key_path")
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error opening jwt private key file %s: %w", privateKeyPath, err)
	}
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing jwt private key file %s: %w", privateKeyPath, err)
	}
	return privateKey, nil
}

func loadJwtPublicKey(config *koanf.Koanf) (*ecdsa.PublicKey, error) {
	publicKeyPath := config.MustString("general.auth.jwt.public_key_path")
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error opening jwt public key file %s: %w", publicKeyPath, err)
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing jwt public key file %s: %w", publicKeyPath, err)
	}
	return publicKey, nil
}
