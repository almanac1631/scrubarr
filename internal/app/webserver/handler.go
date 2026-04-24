package webserver

import (
	"fmt"
	"net/http"

	"github.com/almanac1631/scrubarr/internal/app/auth"
	"github.com/almanac1631/scrubarr/internal/utils"
	"github.com/knadh/koanf/v2"
)

type handler struct {
	version          string
	pathPrefix       string
	authProvider     auth.Provider
	templateCache    TemplateCache
	inventoryService InventoryService
	quotaService     QuotaService
	jwtConfig        *JwtConfig
}

func newHandler(config *koanf.Koanf, version, pathPrefix string, authProvider auth.Provider, templateCache TemplateCache, inventoryService InventoryService, quotaService QuotaService) (*handler, error) {
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
		inventoryService,
		quotaService,
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

// renderPage handles the HTMX/full-page split common to all page endpoints.
// For HTMX requests it renders the "content" sub-template; for full-page requests
// it fetches the disk quota and renders the root template.
// buildData receives the populated basePageData and returns the full endpoint-specific data struct.
func (handler *handler) renderPage(
	writer http.ResponseWriter,
	request *http.Request,
	templateFile, pageTitle string,
	buildData func(base basePageData) any,
) {
	logger := getRequestLogger(request)
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	sortInfo := getSortInfoFromUrlQuery(request.URL.Query())
	base := basePageData{SortInfo: sortInfo, PageTitle: pageTitle}
	if utils.IsHTMXRequest(request) {
		if err := handler.ExecuteSubTemplate(writer, templateFile, "content", buildData(base)); err != nil {
			logger.Error(err.Error())
		}
		return
	}
	diskQuota, err := handler.quotaService.GetDiskQuota()
	if err != nil {
		logger.Error("could not get disk quota", "err", err)
	}
	base.Version = handler.version
	base.DiskQuota = diskQuota
	if err := handler.ExecuteRootTemplate(writer, templateFile, buildData(base)); err != nil {
		logger.Error(err.Error())
	}
}

func (handler *handler) ExecuteSubTemplate(writer http.ResponseWriter, fileName, templateName string, data any) error {
	if err := handler.templateCache[fileName].ExecuteTemplate(writer, templateName, data); isErrAndNoBrokenPipe(err) {
		return fmt.Errorf("failed to execute sub template: %w", err)
	}
	return nil
}
