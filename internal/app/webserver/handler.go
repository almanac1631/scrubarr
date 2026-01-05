package webserver

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/almanac1631/scrubarr/pkg/common"
	"github.com/almanac1631/scrubarr/pkg/inmemory"
	"github.com/almanac1631/scrubarr/pkg/media"
	"github.com/knadh/koanf/v2"
)

type handler struct {
	templateCache TemplateCache

	manager common.Manager

	pathPrefix string

	jwtConfig         *JwtConfig
	username          string
	passwordRetriever func() []byte
	passwordSalt      []byte
}

func newHandler(config *koanf.Koanf, pathPrefix string, templateCache TemplateCache, radarrRetriever *media.RadarrRetriever, sonarrRetriever *media.SonarrRetriever, torrentManager common.TorrentClientManager) (*handler, error) {
	manager := inmemory.NewManager([]common.MediaRetriever{radarrRetriever, sonarrRetriever}, torrentManager)

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
		pathPrefix,
		jwtConfig,
		username,
		passwordRetriever,
		passwordSalt,
	}, nil
}

func (handler *handler) ExecuteRootTemplate(writer http.ResponseWriter, fileName string, data any) error {
	wrappedData := struct {
		PathPrefix string
		Data       any
	}{
		PathPrefix: handler.pathPrefix,
		Data:       data,
	}
	if err := handler.templateCache[fileName].ExecuteTemplate(writer, "index", wrappedData); isErrAndNoBrokenPipe(err) {
		return fmt.Errorf("failed to execute root template: %w", err)
	}
	return nil
}

func (handler *handler) ExecuteSubTemplate(writer http.ResponseWriter, fileName, templateName string, data any) error {
	if err := handler.templateCache[fileName].ExecuteTemplate(writer, templateName, data); isErrAndNoBrokenPipe(err) {
		return fmt.Errorf("failed to execute sub template: %w", err)
	}
	return nil
}
