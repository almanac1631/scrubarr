package webserver

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/almanac1631/scrubarr/internal/app/auth"
	"github.com/knadh/koanf/v2"
)

var authProviderRegistry = map[string]func(map[string]string) (auth.Provider, error){
	"passwordhash": loadPasswordHashAuthProvider,
	"jellyfin":     loadJellyfinAuthProvider,
}

func GetAuthProvider(k *koanf.Koanf) (auth.Provider, error) {
	providerToUse := strings.ToLower(k.MustString("general.auth.provider"))
	providerInstantiator, ok := authProviderRegistry[providerToUse]
	if !ok {
		return nil, fmt.Errorf("unknown auth provider %q", providerToUse)
	}
	providerConfig := k.MustStringMap(fmt.Sprintf("general.auth.providers.%s", providerToUse))
	provider, err := providerInstantiator(providerConfig)
	if err != nil {
		return nil, fmt.Errorf("error instantiating auth provider %q: %w", providerToUse, err)
	}
	slog.Info("Successfully set up auth provider.", "auth-provider", providerToUse)
	return provider, nil
}

func loadPasswordHashAuthProvider(config map[string]string) (auth.Provider, error) {
	username := strings.ToLower(config["username"])
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}
	loadByteValue := func(path string) ([]byte, error) {
		value, err := hex.DecodeString(config[path])
		if err != nil {
			return nil, fmt.Errorf("error decoding hex value on path %q: %w", strconv.Quote(path), err)
		}
		if len(value) == 0 {
			return nil, fmt.Errorf("value for %q cannot be empty", path)
		}
		return value, nil
	}
	passwordHash, err := loadByteValue("password_hash")
	if err != nil {
		return nil, err
	}

	passwordSalt, err := loadByteValue("password_salt")
	if err != nil {
		return nil, err
	}
	return auth.NewPasswordBasedProvider(username, passwordHash, passwordSalt), nil
}

func loadJellyfinAuthProvider(config map[string]string) (auth.Provider, error) {
	baseUrl := config["base_url"]
	if baseUrl == "" {
		return nil, fmt.Errorf("base_url is required")
	}
	return auth.NewJellyfinProvider(baseUrl), nil
}
