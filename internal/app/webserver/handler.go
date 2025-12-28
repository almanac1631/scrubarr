package webserver

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/almanac1631/scrubarr/pkg/common"
	"github.com/almanac1631/scrubarr/pkg/inmemory"
	"github.com/almanac1631/scrubarr/pkg/media"
	"github.com/almanac1631/scrubarr/pkg/torrentclients"
	"github.com/knadh/koanf/v2"
)

type handler struct {
	templateCache TemplateCache

	manager common.Manager

	jwtConfig         *JwtConfig
	username          string
	passwordRetriever func() []byte
	passwordSalt      []byte
}

func newHandler(config *koanf.Koanf, templateCache TemplateCache, radarrRetriever *media.RadarrRetriever, delugeRetriever *torrentclients.DelugeRetriever, rtorrentRetriever *torrentclients.RtorrentRetriever) (*handler, error) {
	manager := inmemory.NewManager(radarrRetriever, delugeRetriever, rtorrentRetriever)

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
	return &handler{
		templateCache,
		manager,
		jwtConfig, username,
		passwordRetriever,
		passwordSalt,
	}, nil
}
