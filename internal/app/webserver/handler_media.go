package webserver

import (
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"github.com/almanac1631/scrubarr/internal/utils"
	"github.com/almanac1631/scrubarr/pkg/domain"
)

type mediaEndpointData struct {
	MappedMedia []*MappedMedia
	SortInfo    domain.SortInfo
	NextPage    int
	Version     string
}

type MappedMedia struct {
	domain.MatchedMedia
	TorrentInformation domain.TorrentInformation
	Size               int64
}

type MappedMediaSeason struct {
	Season             int
	Size               int64
	TorrentInformation domain.TorrentInformation
	Parts              []domain.MatchedMediaPart
}

type MappedMediaSeries struct {
	*MappedMedia
	Seasons []*MappedMediaSeason
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
	matchedMediaList, hasNext, err := handler.manager.GetMatchedMedia(page, sortInfo)
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
	mediaEntries := handler.getMatchedMediaList(matchedMediaList)
	for _, mediaEntry := range mediaEntries {
		mediaEntry.Parts = []domain.MatchedMediaPart{}
	}
	if err = handler.ExecuteSubTemplate(writer, "media.gohtml", "media_entries", mediaEndpointData{
		MappedMedia: mediaEntries,
		SortInfo:    sortInfo,
		NextPage:    nextPage,
	}); err != nil {
		logger.Error(err.Error())
		return
	}
	return
}

func (handler *handler) getMatchedMediaList(matchedMediaList []domain.MatchedMedia) []*MappedMedia {
	mediaEntries := make([]*MappedMedia, 0, len(matchedMediaList))
	for _, matchedMedia := range matchedMediaList {
		torrentInformation := getBundledTorrentInformationFromParts(matchedMedia.Parts)
		mediaEntries = append(mediaEntries, &MappedMedia{
			MatchedMedia:       matchedMedia,
			TorrentInformation: torrentInformation,
			Size:               matchedMedia.Size,
		})
	}
	return mediaEntries
}

func getBundledTorrentInformationFromParts(parts []domain.MatchedMediaPart) domain.TorrentInformation {
	torrentInformation := &domain.TorrentInformation{
		Status:      domain.TorrentStatusPresent,
		Tracker:     domain.Tracker{},
		RatioStatus: domain.TorrentAttributeStatusFulfilled,
		Ratio:       -1,
		AgeStatus:   domain.TorrentAttributeStatusFulfilled,
		Age:         -1,
	}
	missingTorrents := 0
	for _, part := range parts {
		torrentInformation.RatioStatus = getNewTorrentAttributeStatus(torrentInformation.RatioStatus, part.TorrentInformation.RatioStatus)
		torrentInformation.AgeStatus = getNewTorrentAttributeStatus(torrentInformation.AgeStatus, part.TorrentInformation.AgeStatus)

		if part.TorrentInformation.Status == domain.TorrentStatusMissing {
			missingTorrents++
		}

		if len(parts) == 1 {
			torrentInformation.Tracker = part.TorrentInformation.Tracker
			if part.TorrentInformation.Status == domain.TorrentStatusPresent {
				torrentInformation.Ratio = part.TorrentInformation.Ratio
				torrentInformation.Age = part.TorrentInformation.Age
			}

			if part.TorrentInformation.Tracker.IsValid() {
				torrentInformation.Tracker.MinRatio = part.TorrentInformation.Tracker.MinRatio
				torrentInformation.Tracker.MinAge = part.TorrentInformation.Tracker.MinAge
			}
		}
	}
	if missingTorrents == len(parts) {
		torrentInformation.Status = domain.TorrentStatusMissing
	} else if missingTorrents != 0 {
		torrentInformation.Status = domain.TorrentStatusIncomplete
	}
	return *torrentInformation
}

func getNewTorrentAttributeStatus(torrentInformationStatus, torrentStatus domain.TorrentAttributeStatus) domain.TorrentAttributeStatus {
	if torrentInformationStatus == domain.TorrentAttributeStatusUnknown || torrentStatus == domain.TorrentAttributeStatusUnknown {
		return domain.TorrentAttributeStatusUnknown
	} else if torrentInformationStatus == domain.TorrentAttributeStatusFulfilled && torrentStatus == domain.TorrentAttributeStatusFulfilled {
		return domain.TorrentAttributeStatusFulfilled
	}
	return domain.TorrentAttributeStatusPending
}

func getSortInfoFromUrlQuery(values url.Values) domain.SortInfo {
	sortInfo := domain.SortInfo{}
	sortKeyRaw := values.Get("sortKey")
	switch domain.SortKey(sortKeyRaw) {
	case domain.SortKeyName, domain.SortKeySize, domain.SortKeyAdded, domain.SortKeyTorrentStatus:
		sortInfo.Key = domain.SortKey(sortKeyRaw)
	default:
		sortInfo.Key = domain.SortKeyName
	}
	sortOrderRaw := values.Get("sortOrder")
	switch domain.SortOrder(sortOrderRaw) {
	case domain.SortOrderAsc, domain.SortOrderDesc:
		sortInfo.Order = domain.SortOrder(sortOrderRaw)
	default:
		sortInfo.Order = domain.SortOrderAsc
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
	handler.serveMediaSeriesEntry(writer, request, seriesId, collapsed)
}

func (handler *handler) serveMediaSeriesEntry(writer http.ResponseWriter, request *http.Request, seriesId int64, collapsed bool) {
	logger := getRequestLogger(request)
	media, err := handler.manager.GetMatchedMediaBySeriesId(seriesId)
	if errors.Is(err, domain.ErrMediaNotFound) {
		// write 200 because of HTMX request
		writer.WriteHeader(http.StatusOK)
		return
	} else if err != nil {
		logger.Error(err.Error())
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	mappedMedias := handler.getMatchedMediaList([]domain.MatchedMedia{media})
	if len(mappedMedias) == 0 {
		http.NotFound(writer, request)
		return
	}
	mappedMedia := mappedMedias[0]
	if collapsed {
		mappedMedia.Parts = []domain.MatchedMediaPart{}
	}
	if err = handler.ExecuteSubTemplate(writer, "media.gohtml", "media_entry", mappedMedia); err != nil {
		logger.Error(err.Error())
		return
	}
	return
}

func getSeasonGroupedParts(mappedMedia *MappedMedia) MappedMediaSeries {
	mappedSeries := MappedMediaSeries{
		MappedMedia: mappedMedia,
		Seasons:     make([]*MappedMediaSeason, 0),
	}
	partsNoSeason := make([]domain.MatchedMediaPart, 0)
	for _, part := range mappedMedia.Parts {
		seasonNumber := part.MediaPart.Season
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
				Parts:  []domain.MatchedMediaPart{part},
				Size:   part.MediaPart.Size,
			}
			mappedSeries.Seasons = append(mappedSeries.Seasons, seasonObj)
		} else {
			seasonObj = mappedSeries.Seasons[index]
			seasonObj.Parts = append(seasonObj.Parts, part)
			seasonObj.Size += part.MediaPart.Size
		}
	}
	for _, season := range mappedSeries.Seasons {
		season.TorrentInformation = getBundledTorrentInformationFromParts(season.Parts)
	}
	mappedMedia.Parts = partsNoSeason
	return mappedSeries
}

func (handler *handler) getMediaDeletionHandler(mediaType domain.MediaType) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger := getRequestLogger(request)
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
		logger = logger.With("mediaType", mediaType, "id", id)
		logger.Debug("Deleting media...")
		if err := handler.manager.DeleteMedia(mediaType, id); err != nil {
			logger.Error("Could not delete media.", "error", err)
			http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		logger.Info("Successfully deleted media.")
		writer.WriteHeader(http.StatusOK)
		return
	}
}

func (handler *handler) getMediaSeasonDeletionHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger := getRequestLogger(request)
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
		seasonStr := request.PathValue("season")
		season, err := strconv.Atoi(seasonStr)
		if err != nil {
			http.Error(writer, "400 Bad Request", http.StatusBadRequest)
			return
		}
		logger = logger.With("mediaType", domain.MediaTypeSeries, "season", season, "id", id)
		logger.Debug("Deleting series season...")
		if err = handler.manager.DeleteSeason(id, season); err != nil {
			logger.Error("Could not delete series season.", "error", err)
			http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		logger.Info("Successfully series season.")
		handler.serveMediaSeriesEntry(writer, request, id, false)
	}
}
