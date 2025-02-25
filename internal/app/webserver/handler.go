package webserver

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/golang-jwt/jwt/v5"
	"github.com/knadh/koanf/v2"
	"os"
	"strconv"
	"strings"
)

type JwtConfig struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

type ApiEndpointHandler struct {
	entryMappingManager common.EntryMappingManager
	jwtConfig           *JwtConfig
	username            string
	passwordRetriever   func() []byte
	passwordSalt        []byte
}

func NewApiEndpointHandler(entryMappingManager common.EntryMappingManager, config *koanf.Koanf) (*ApiEndpointHandler, error) {
	username := strings.ToLower(config.MustString("general.auth.username"))
	loadByteValue := func(path string) ([]byte, error) {
		value, err := hex.DecodeString(config.MustString(path))
		if err != nil {
			return nil, fmt.Errorf("error decoding hex value on path %s: %w", strconv.Quote(path), err)
		}
		return value, nil
	}
	passwordSalt, err := loadByteValue("general.auth.password_salt")
	if err != nil {
		return nil, err
	}
	_, err = loadByteValue("general.auth.password_hash")
	if err != nil {
		return nil, err
	}
	passwordRetriever := func() []byte {
		passwordHash, _ := hex.DecodeString(config.MustString("general.auth.password_hash"))
		return passwordHash
	}
	privateKey, err := loadJwtPrivateKey(config)
	if err != nil {
		return nil, err
	}
	publicKey, err := loadJwtPublicKey(config)
	if err != nil {
		return nil, err
	}
	jwtConfig := &JwtConfig{privateKey, publicKey}
	return &ApiEndpointHandler{entryMappingManager, jwtConfig, username, passwordRetriever, passwordSalt}, nil
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

var _ StrictServerInterface = (*ApiEndpointHandler)(nil)
