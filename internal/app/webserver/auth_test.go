package webserver

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"reflect"
	"testing"
)

func Test_checkPassword(t *testing.T) {
	type args struct {
		passwordHashExpected []byte
		passwordRawActual    []byte
		salt                 []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test hash matches", args{hexDecode("0d71f11951721c1d4cd4273a696eefc0"), []byte("SomeSecurePassword"), []byte("WFqP9t2QwwUjwiOu")}, true},
		{"test hash no match", args{hexDecode("0d71f11951721c1d4cd4273a696eefc0"), []byte("AnotherSecurePassword"), []byte("WFqP9t2QwwUjwiOu")}, false},
		{"test hash no match different salt", args{hexDecode("0d71f11951721c1d4cd4273a696eefc0"), []byte("SomeSecurePassword"), []byte("SomeOtherSalt")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkPassword(tt.args.passwordHashExpected, tt.args.passwordRawActual, tt.args.salt); got != tt.want {
				t.Errorf("checkPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateHash(t *testing.T) {
	type args struct {
		passwordRaw []byte
		salt        []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"generate one argon2id hash", args{[]byte("Hello"), []byte("RNaMJfQ1owJktbnj")}, hexDecode("8065c3b981f5f3cfdd7c6309d0dbdc6a")},
		{"generate another argon2id hash", args{[]byte("Bye"), []byte("CqC6mbILITnHwLUD")}, hexDecode("feea14dee4899829af6bff741f85fcb0")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateHash(tt.args.passwordRaw, tt.args.salt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func hexDecode(strToDecode string) []byte {
	result, err := hex.DecodeString(strToDecode)
	if err != nil {
		panic(err)
	}
	return result
}

// openssl ecparam -name prime256v1 -genkey -noout
const testPrivateCert = "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIHdvea3pJSgZ5NuoVPEh2AkM1LHDMQtn91bYU0Gdpx8eoAoGCCqGSM49\nAwEHoUQDQgAE9cma9SWlH+JDigPxhw45Xhu+Gjsj9CQ8ybZ1yLEH5vZhsCe7RCDE\nbUlD3gAHjL5oY6+vQQcGGcUfmmx8EdZI4Q==\n-----END EC PRIVATE KEY-----\n" //gitleaks:allow

// openssl ec -pubout
const testPublicCert = "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE9cma9SWlH+JDigPxhw45Xhu+Gjsj\n9CQ8ybZ1yLEH5vZhsCe7RCDEbUlD3gAHjL5oY6+vQQcGGcUfmmx8EdZI4Q==\n-----END PUBLIC KEY-----"

func Test_generateToken(t *testing.T) {
	block, _ := pem.Decode([]byte(testPrivateCert))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		panic(fmt.Errorf("failed to decode PEM block containing EC private key"))
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	username := "blinks"
	jwtTokenStr, err := generateToken(key, username)
	jwtToken, err := jwt.Parse(jwtTokenStr, func(token *jwt.Token) (interface{}, error) {
		blockPub, _ := pem.Decode([]byte(testPublicCert))
		if blockPub == nil || blockPub.Type != "PUBLIC KEY" {
			panic(fmt.Errorf("failed to decode PEM block containing public key"))
		}

		keyPub, err := x509.ParsePKIXPublicKey(blockPub.Bytes)
		if err != nil {
			panic(err)
		}
		return keyPub, nil
	})
	if err != nil {
		panic(err)
	}
	if !jwtToken.Valid {
		t.Errorf("JWT token is not valid")
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	err = jwt.NewValidator(
		jwt.WithSubject(username),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	).Validate(claims)
	if err != nil {
		t.Errorf("JWT token is not valid: %e", err)
	}
}
