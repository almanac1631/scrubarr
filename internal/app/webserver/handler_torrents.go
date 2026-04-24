package webserver

import (
	"errors"
	"net/http"
	"strconv"
)

type torrentsEndpointData struct {
	basePageData
	Rows []OrphanedTorrentRow
}

func (handler *handler) handleTorrentsEndpoint(writer http.ResponseWriter, request *http.Request) {
	handler.renderPage(writer, request, "torrents.gohtml", "Torrents", func(base basePageData) any {
		return torrentsEndpointData{basePageData: base}
	})
}

func (handler *handler) handleTorrentEntriesEndpoint(writer http.ResponseWriter, request *http.Request) {
	logger := getRequestLogger(request)
	sortInfo := getSortInfoFromUrlQuery(request.URL.Query())
	pageRaw := request.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageRaw)
	if page < 1 {
		page = 1
	}
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	rows, hasNext, err := handler.inventoryService.GetOrphanedTorrents(page, sortInfo)
	if err != nil {
		logger.Error("Failed to get orphaned torrents.", "error", err)
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	nextPage := -1
	if hasNext {
		nextPage = page + 1
	}
	if err = handler.ExecuteSubTemplate(writer, "torrents.gohtml", "torrent_entries", torrentsEndpointData{
		basePageData: basePageData{SortInfo: sortInfo, NextPage: nextPage},
		Rows:         rows,
	}); err != nil {
		logger.Error(err.Error())
	}
}

func (handler *handler) handleTorrentDeletionEndpoint(writer http.ResponseWriter, request *http.Request) {
	logger := getRequestLogger(request)
	id := request.PathValue("id")
	logger = logger.With("id", id)
	logger.Debug("Deleting orphaned torrent...")
	if err := handler.inventoryService.DeleteOrphanedTorrent(id); errors.Is(err, ErrMediaNotFound) {
		writer.Header().Set("Hx-Trigger", "diskQuotaUpdate")
		writer.WriteHeader(http.StatusOK)
		return
	} else if err != nil {
		logger.Error("Could not delete orphaned torrent.", "error", err)
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	logger.Info("Successfully deleted orphaned torrent.")
	writer.Header().Set("Hx-Trigger", "diskQuotaUpdate")
	writer.WriteHeader(http.StatusOK)
}
