package auth

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func hexDecode(strToDecode string) []byte {
	result, err := hex.DecodeString(strToDecode)
	if err != nil {
		panic(err)
	}
	return result
}
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

func Test_GenerateHash(t *testing.T) {
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
			if got := GenerateHash(tt.args.passwordRaw, tt.args.salt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
