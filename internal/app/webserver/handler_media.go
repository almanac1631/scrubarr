package webserver

import (
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/almanac1631/scrubarr/internal/utils"
	"github.com/almanac1631/scrubarr/pkg/common"
)

type TorrentStatus string

const (
	TorrentStatusMissing    TorrentStatus = "missing"
	TorrentStatusPresent    TorrentStatus = "present"
	TorrentStatusIncomplete TorrentStatus = "incomplete"
)

type mediaEndpointData struct {
	MappedMedia []*MappedMedia
	SortInfo    common.SortInfo
	NextPage    int
}

type MappedMedia struct {
	common.MatchedMedia
	TorrentStatus TorrentStatus
	Size          int64
}

func (handler *handler) handleMediaEndpoint(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	sortInfo := getSortInfoFromUrlQuery(request.URL.Query())
	if utils.IsHTMXRequest(request) {
		if err := handler.ExecuteSubTemplate(writer, "media.gohtml", "content", mediaEndpointData{
			SortInfo: sortInfo,
		}); err != nil {
			slog.Error(err.Error())
			return
		}
	} else {
		if err := handler.ExecuteRootTemplate(writer, "media.gohtml", mediaEndpointData{
			SortInfo: sortInfo,
		}); err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

func (handler *handler) handleMediaEntriesEndpoint(writer http.ResponseWriter, request *http.Request) {
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
	matchedMediaList, hasNext, err := handler.manager.GetMatchedMedia(page, sortInfo)
	if err != nil {
		slog.Error("failed to get movie mapping", "err", err)
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("500 Internal Server Error"))
		return
	}
	nextPage := -1
	if hasNext {
		nextPage = page + 1
	}
	mediaEntries := handler.getMatchedMediaList(matchedMediaList)
	for _, mediaEntry := range mediaEntries {
		mediaEntry.Parts = []common.MatchedMediaPart{}
	}
	if err = handler.ExecuteSubTemplate(writer, "media.gohtml", "media_entries", mediaEndpointData{
		MappedMedia: mediaEntries,
		SortInfo:    sortInfo,
		NextPage:    nextPage,
	}); err != nil {
		slog.Error(err.Error())
		return
	}
	return
}

func (handler *handler) getMatchedMediaList(matchedMediaList []common.MatchedMedia) []*MappedMedia {
	mediaEntries := make([]*MappedMedia, 0, len(matchedMediaList))
	for _, matchedMedia := range matchedMediaList {
		totalSize := int64(0)
		missingTorrents := 0
		for _, part := range matchedMedia.Parts {
			totalSize += part.Size
			if part.TorrentFinding == nil {
				missingTorrents++
			}
		}
		var torrentStatus TorrentStatus
		if missingTorrents == 0 {
			torrentStatus = TorrentStatusPresent
		} else if missingTorrents == len(matchedMedia.Parts) {
			torrentStatus = TorrentStatusMissing
		} else {
			torrentStatus = TorrentStatusIncomplete
		}
		mediaEntries = append(mediaEntries, &MappedMedia{
			MatchedMedia:  matchedMedia,
			TorrentStatus: torrentStatus,
			Size:          totalSize,
		})
	}
	return mediaEntries
}

func getSortInfoFromUrlQuery(values url.Values) common.SortInfo {
	sortInfo := common.SortInfo{}
	sortKeyRaw := values.Get("sortKey")
	switch common.SortKey(sortKeyRaw) {
	case common.SortKeyName, common.SortKeySize, common.SortKeyAdded, common.SortKeyTorrentStatus:
		sortInfo.Key = common.SortKey(sortKeyRaw)
	default:
		sortInfo.Key = common.SortKeyName
	}
	sortOrderRaw := values.Get("sortOrder")
	switch common.SortOrder(sortOrderRaw) {
	case common.SortOrderAsc, common.SortOrderDesc:
		sortInfo.Order = common.SortOrder(sortOrderRaw)
	default:
		sortInfo.Order = common.SortOrderAsc
	}
	return sortInfo
}

func (handler *handler) handleMediaSeriesEndpoint(writer http.ResponseWriter, request *http.Request) {
	if !utils.IsHTMXRequest(request) {
		http.Error(writer, "404 Not Found", http.StatusNotFound)
		return
	}
	collapsed := request.URL.Query().Get("collapsed") == "true"
	idString := request.PathValue("id")
	seriesId, _ := strconv.ParseInt(idString, 10, 64)
	media, err := handler.manager.GetMatchedMediaBySeriesId(seriesId)
	if err != nil {
		slog.Error(err.Error())
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	mappedMedia := handler.getMatchedMediaList(media)[0]
	if collapsed {
		mappedMedia.Parts = []common.MatchedMediaPart{}
	}
	if err = handler.ExecuteSubTemplate(writer, "media.gohtml", "media_entry", mappedMedia); err != nil {
		slog.Error(err.Error())
		return
	}
	return
}

func (handler *handler) getMediaDeletionHandler(mediaType common.MediaType) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if !utils.IsHTMXRequest(request) {
			http.Error(writer, "404 Not Found", http.StatusNotFound)
			return
		}
		idStr := request.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(writer, "400 Bad Request", http.StatusBadRequest)
			return
		}
		if err := handler.manager.DeleteMedia(mediaType, id); err != nil {
			slog.Error(err.Error())
			http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		return
	}
}
