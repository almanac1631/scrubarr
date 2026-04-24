package webserver

import (
	"errors"
	"net/http"
	"strconv"
)

type mediaEndpointData struct {
	basePageData
	Rows []MediaRow
}

func (handler *handler) handleMediaEndpoint(writer http.ResponseWriter, request *http.Request) {
	handler.renderPage(writer, request, "media.gohtml", "Media", func(base basePageData) any {
		return mediaEndpointData{basePageData: base}
	})
}

func (handler *handler) handleMediaEntriesEndpoint(writer http.ResponseWriter, request *http.Request) {
	logger := getRequestLogger(request)
	sortInfo := getSortInfoFromUrlQuery(request.URL.Query())
	pageRaw := request.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageRaw)
	if page < 1 {
		page = 1
	}
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	mediaRows, hasNext, err := handler.inventoryService.GetMediaInventory(page, sortInfo)
	if err != nil {
		logger.Error("Failed to get media mapping.", "error", err)
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	nextPage := -1
	if hasNext {
		nextPage = page + 1
	}
	if err = handler.ExecuteSubTemplate(writer, "media.gohtml", "media_entries", mediaEndpointData{
		basePageData: basePageData{SortInfo: sortInfo, NextPage: nextPage},
		Rows:         mediaRows,
	}); err != nil {
		logger.Error(err.Error())
	}
}

func (handler *handler) handleMediaSeriesEndpoint(writer http.ResponseWriter, request *http.Request) {
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
	}
}

func (handler *handler) handleMediaDeletionEndpoint(writer http.ResponseWriter, request *http.Request) {
	logger := getRequestLogger(request)
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
	writer.Header().Set("Hx-Trigger", "diskQuotaUpdate")
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
	logger.Info("Refreshing media entries cache.")
	if err := handler.inventoryService.RefreshCache(); err != nil {
		logger.Error(err.Error())
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	logger.Info("Successfully refreshed media entries cache.")
	writer.Header().Set("Hx-Trigger", "diskQuotaUpdate")
	handler.handleMediaEndpoint(writer, request)
}
