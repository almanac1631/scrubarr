package webserver

import (
	"log/slog"
	"net/http"
	"net/url"
	"slices"
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

type MappedMediaSeason struct {
	Season        int
	Size          int64
	TorrentStatus TorrentStatus
	Parts         []common.MatchedMediaPart
}

type MappedMediaSeries struct {
	*MappedMedia
	Seasons []*MappedMediaSeason
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
		totalSize, torrentStatus := getStatusFromParts(matchedMedia.Parts)
		mediaEntries = append(mediaEntries, &MappedMedia{
			MatchedMedia:  matchedMedia,
			TorrentStatus: torrentStatus,
			Size:          totalSize,
		})
	}
	return mediaEntries
}

func getStatusFromParts(parts []common.MatchedMediaPart) (totalSize int64, torrentStatus TorrentStatus) {
	missingTorrents := 0
	for _, part := range parts {
		totalSize += part.Size
		if part.TorrentFinding == nil {
			missingTorrents++
		}
	}
	if missingTorrents == 0 {
		torrentStatus = TorrentStatusPresent
	} else if missingTorrents == len(parts) {
		torrentStatus = TorrentStatusMissing
	} else {
		torrentStatus = TorrentStatusIncomplete
	}
	return totalSize, torrentStatus
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

func getSeasonGroupedParts(mappedMedia *MappedMedia) MappedMediaSeries {
	mappedSeries := MappedMediaSeries{
		MappedMedia: mappedMedia,
		Seasons:     make([]*MappedMediaSeason, 0),
	}
	partsNoSeason := make([]common.MatchedMediaPart, 0)
	for _, part := range mappedMedia.Parts {
		seasonNumber := part.Season
		if seasonNumber == 0 {
			partsNoSeason = append(partsNoSeason, part)
			continue
		}
		index := slices.IndexFunc(mappedSeries.Seasons, func(season *MappedMediaSeason) bool {
			return season.Season == seasonNumber
		})
		var seasonObj *MappedMediaSeason
		if index == -1 {
			seasonObj = &MappedMediaSeason{
				Season: seasonNumber,
				Parts:  []common.MatchedMediaPart{part},
			}
			mappedSeries.Seasons = append(mappedSeries.Seasons, seasonObj)
		} else {
			seasonObj = mappedSeries.Seasons[index]
			seasonObj.Parts = append(seasonObj.Parts, part)
		}
	}
	for _, season := range mappedSeries.Seasons {
		season.Size, season.TorrentStatus = getStatusFromParts(season.Parts)
	}
	mappedMedia.Parts = partsNoSeason
	return mappedSeries
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
