package webserver

import (
	"crypto/ecdsa"
	"encoding/hex"
	"reflect"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
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

func failSafePEMDecodeStr(pemStr string) {

}

// openssl ecparam -name prime256v1 -genkey -noout
var testPrivateCert, _ = jwt.ParseECPrivateKeyFromPEM([]byte("-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIHdvea3pJSgZ5NuoVPEh2AkM1LHDMQtn91bYU0Gdpx8eoAoGCCqGSM49\nAwEHoUQDQgAE9cma9SWlH+JDigPxhw45Xhu+Gjsj9CQ8ybZ1yLEH5vZhsCe7RCDE\nbUlD3gAHjL5oY6+vQQcGGcUfmmx8EdZI4Q==\n-----END EC PRIVATE KEY-----\n")) //gitleaks:allow

// openssl ec -pubout
var testPublicCert, _ = jwt.ParseECPublicKeyFromPEM([]byte("-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE9cma9SWlH+JDigPxhw45Xhu+Gjsj\n9CQ8ybZ1yLEH5vZhsCe7RCDEbUlD3gAHjL5oY6+vQQcGGcUfmmx8EdZI4Q==\n-----END PUBLIC KEY-----\n"))

func Test_generateToken(t *testing.T) {
	username := "blinks"
	jwtTokenStr, err := generateToken(testPrivateCert, username)
	jwtToken, err := jwt.Parse(jwtTokenStr, func(token *jwt.Token) (interface{}, error) {
		return testPublicCert, nil
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

func Test_validateToken(t *testing.T) {
	validToken, _ := generateToken(testPrivateCert, "blinks")
	expiredToken := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk4MTU5MzYsImlhdCI6MTczOTgxMjMzNiwic3ViIjoiYmxpbmtzIn0.m5_fHTf2CEPpzEDCBfwdCR109b4gvoyIzxZ7U4zYNCgPFeuSsG4BMqQrnNfGbb7b9OSLT4ST2_irXEgRqHEM9g" //gitleaks:allow
	testToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"                                        //gitleaks:allow
	malformedToken := "sometokenthat.is.notjwt"

	type args struct {
		key    *ecdsa.PublicKey
		jwtStr string
	}
	type want struct {
		ok       bool
		username string
		errStr   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"validates valid token", args{testPublicCert, validToken}, want{true, "blinks", ""}},
		{"does not validate expired token", args{testPublicCert, expiredToken}, want{false, "", "expired"}},
		{"does not validate test token with invalid mechanism", args{testPublicCert, testToken}, want{false, "", "unexpected signing method"}},
		{"does not validate malformed token", args{testPublicCert, malformedToken}, want{false, "", "token is malformed"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotUsername, gotErr := validateToken(tt.args.key, tt.args.jwtStr)
			if gotOk != tt.want.ok || gotUsername != tt.want.username ||
				(tt.want.errStr != "" && !strings.Contains(gotErr.Error(), tt.want.errStr)) {
				t.Errorf("validateToken() = (%v, %s, %s), want (%v, %s, error string: %v)", gotOk, gotUsername, gotErr, tt.want.ok, tt.want.username, tt.want.errStr)
			}
		})
	}
}
