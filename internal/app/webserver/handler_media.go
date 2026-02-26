package webserver

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/almanac1631/scrubarr/internal/utils"
)

type mediaEndpointData struct {
	MediaRows []MediaRow
	SortInfo  SortInfo
	NextPage  int
	Version   string
}

func (handler *handler) handleMediaEndpoint(writer http.ResponseWriter, request *http.Request) {
	logger := getRequestLogger(request)
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	sortInfo := getSortInfoFromUrlQuery(request.URL.Query())
	if utils.IsHTMXRequest(request) {
		if err := handler.ExecuteSubTemplate(writer, "media.gohtml", "content", mediaEndpointData{
			SortInfo: sortInfo,
		}); err != nil {
			logger.Error(err.Error())
			return
		}
	} else {
		if err := handler.ExecuteRootTemplate(writer, "media.gohtml", mediaEndpointData{
			SortInfo: sortInfo,
			Version:  handler.version,
		}); err != nil {
			logger.Error(err.Error())
			return
		}
	}
}

func (handler *handler) handleMediaEntriesEndpoint(writer http.ResponseWriter, request *http.Request) {
	logger := getRequestLogger(request)
	sortInfo := getSortInfoFromUrlQuery(request.URL.Query())
	if !utils.IsHTMXRequest(request) {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("404 Not Found"))
		return
	}
	pageRaw := request.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageRaw)
	if page < 1 {
		page = 1
	}
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	mediaRows, hasNext, err := handler.inventoryService.GetMediaInventory(page, sortInfo)
	if err != nil {
		logger.Error("Failed to get media mapping.", "error", err)
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("500 Internal Server Error"))
		return
	}
	nextPage := -1
	if hasNext {
		nextPage = page + 1
	}
	if err = handler.ExecuteSubTemplate(writer, "media.gohtml", "media_entries", mediaEndpointData{
		MediaRows: mediaRows,
		SortInfo:  sortInfo,
		NextPage:  nextPage,
	}); err != nil {
		logger.Error(err.Error())
		return
	}
	return
}

func getSortInfoFromUrlQuery(values url.Values) SortInfo {
	sortInfo := SortInfo{}
	sortKeyRaw := values.Get("sortKey")
	switch SortKey(sortKeyRaw) {
	case SortKeyName, SortKeySize, SortKeyAdded, SortKeyStatus:
		sortInfo.Key = SortKey(sortKeyRaw)
	default:
		sortInfo.Key = SortKeyName
	}
	sortOrderRaw := values.Get("sortOrder")
	switch SortOrder(sortOrderRaw) {
	case SortOrderAsc, SortOrderDesc:
		sortInfo.Order = SortOrder(sortOrderRaw)
	default:
		sortInfo.Order = SortOrderAsc
	}
	return sortInfo
}

func (handler *handler) handleMediaSeriesEndpoint(writer http.ResponseWriter, request *http.Request) {
	if !utils.IsHTMXRequest(request) {
		http.Error(writer, "404 Not Found", http.StatusNotFound)
		return
	}
	collapsed := request.URL.Query().Get("collapsed") == "true"
	id := request.PathValue("id")
	handler.serveMediaSeriesEntry(writer, request, id, collapsed)
}

func (handler *handler) serveMediaSeriesEntry(writer http.ResponseWriter, request *http.Request, id string, collapsed bool) {
	logger := getRequestLogger(request)
	mediaRowExpanded, err := handler.inventoryService.GetExpandedMediaRow(id)
	if errors.Is(err, ErrMediaNotFound) {
		// write 200 because of HTMX request
		writer.WriteHeader(http.StatusOK)
		return
	} else if err != nil {
		logger.Error(err.Error())
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	if collapsed {
		mediaRowExpanded.ChildMediaRows = []MediaRow{}
	}
	if err = handler.ExecuteSubTemplate(writer, "media.gohtml", "media_entry", mediaRowExpanded); err != nil {
		logger.Error(err.Error())
		return
	}
	return
}

func (handler *handler) handleMediaDeletionEndpoint(writer http.ResponseWriter, request *http.Request) {
	logger := getRequestLogger(request)
	if !utils.IsHTMXRequest(request) {
		http.Error(writer, "404 Not Found", http.StatusNotFound)
		return
	}
	id := request.PathValue("id")
	logger = logger.With("id", id)
	logger.Debug("Deleting media...")
	if err := handler.inventoryService.DeleteMedia(id); err != nil {
		logger.Error("Could not delete media.", "error", err)
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	logger.Info("Successfully deleted media.")
	mediaRowExpanded, err := handler.inventoryService.GetExpandedMediaRow(id)
	if errors.Is(err, ErrMediaNotFound) {
		writer.WriteHeader(http.StatusOK)
		return
	} else if err != nil {
		logger.Error(err.Error())
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err = handler.ExecuteSubTemplate(writer, "media.gohtml", "media_entry", mediaRowExpanded); err != nil {
		logger.Error(err.Error())
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

func (handler *handler) handleRefreshEndpoint(writer http.ResponseWriter, request *http.Request) {
	logger := getRequestLogger(request)
	if !utils.IsHTMXRequest(request) {
		http.Error(writer, "404 Not Found", http.StatusNotFound)
		return
	}
	logger.Info("Refreshing media entries cache.")
	if err := handler.inventoryService.RefreshCache(); err != nil {
		logger.Error(err.Error())
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	logger.Info("Successfully refreshed media entries cache.")
	handler.handleMediaEndpoint(writer, request)
}
