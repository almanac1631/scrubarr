package webserver

import (
	"fmt"
	"net/http"

	"github.com/almanac1631/scrubarr/internal/app/auth"
	"github.com/almanac1631/scrubarr/pkg/common"
	"github.com/almanac1631/scrubarr/pkg/inmemory"
	"github.com/knadh/koanf/v2"
)

type handler struct {
	version       string
	pathPrefix    string
	authProvider  auth.Provider
	templateCache TemplateCache
	manager       common.Manager
	jwtConfig     *JwtConfig
}

func newHandler(config *koanf.Koanf, version, pathPrefix string, authProvider auth.Provider, templateCache TemplateCache,
	mediaManager common.MediaManager, torrentManager common.TorrentClientManager, trackerManager common.TrackerManager) (*handler, error) {
	manager := inmemory.NewManager(mediaManager, torrentManager, trackerManager)

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
		version,
		pathPrefix,
		authProvider,
		templateCache,
		manager,
		jwtConfig,
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
