package webserver

import (
	"net/http"
	"strconv"

	"github.com/almanac1631/scrubarr/internal/utils"
)

type torrentsEndpointData struct {
	Rows      []OrphanedTorrentRow
	NextPage  int
	SortInfo  SortInfo
	Version   string
	DiskQuota DiskQuota
	PageTitle string
}

func (handler *handler) handleTorrentsEndpoint(writer http.ResponseWriter, request *http.Request) {
	logger := getRequestLogger(request)
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	sortInfo := getSortInfoFromUrlQuery(request.URL.Query())
	if utils.IsHTMXRequest(request) {
		if err := handler.ExecuteSubTemplate(writer, "torrents.gohtml", "content", torrentsEndpointData{
			SortInfo:  sortInfo,
			PageTitle: "Torrents",
		}); err != nil {
			logger.Error(err.Error())
		}
		return
	}
	diskQuota, err := handler.quotaService.GetDiskQuota()
	if err != nil {
		logger.Error("could not get disk quota", "err", err)
	}
	if err := handler.ExecuteRootTemplate(writer, "torrents.gohtml", torrentsEndpointData{
		SortInfo:  sortInfo,
		Version:   handler.version,
		DiskQuota: diskQuota,
		PageTitle: "Torrents",
	}); err != nil {
		logger.Error(err.Error())
	}
}

func (handler *handler) handleTorrentEntriesEndpoint(writer http.ResponseWriter, request *http.Request) {
	logger := getRequestLogger(request)
	if !utils.IsHTMXRequest(request) {
		http.Error(writer, "404 Not Found", http.StatusNotFound)
		return
	}
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
		Rows:     rows,
		NextPage: nextPage,
		SortInfo: sortInfo,
	}); err != nil {
		logger.Error(err.Error())
	}
}
