package webserver

import (
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/almanac1631/scrubarr/internal/utils"
	"github.com/almanac1631/scrubarr/pkg/common"
)

type TorrentStatus string

const (
	TorrentStatusMissing    TorrentStatus = "missing"
	TorrentStatusPresent    TorrentStatus = "present"
	TorrentStatusIncomplete TorrentStatus = "incomplete"
)

type TorrentAttributeStatus string

const (
	TorrentAttributeStatusFulfilled TorrentAttributeStatus = "fulfilled"
	TorrentAttributeStatusPending   TorrentAttributeStatus = "pending"
	TorrentAttributeStatusUnknown   TorrentAttributeStatus = "unknown"
)

type TorrentInformation struct {
	Status                 TorrentStatus
	RatioStatus, AgeStatus TorrentAttributeStatus
	Ratio, MinRatio        float64
	Age, MinAge            time.Duration
}

type mediaEndpointData struct {
	MappedMedia []*MappedMedia
	SortInfo    common.SortInfo
	NextPage    int
	Version     string
}

type MappedMedia struct {
	common.MatchedMedia
	TorrentInformation TorrentInformation
	Size               int64
}

type MappedMediaSeason struct {
	Season             int
	Size               int64
	TorrentInformation TorrentInformation
	Parts              []common.MatchedMediaPart
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
		mediaEntry.Parts = []common.MatchedMediaPart{}
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

func (handler *handler) getMatchedMediaList(matchedMediaList []common.MatchedMedia) []*MappedMedia {
	mediaEntries := make([]*MappedMedia, 0, len(matchedMediaList))
	for _, matchedMedia := range matchedMediaList {
		totalSize, torrentInformation := getTorrentInformationFromParts(matchedMedia.Parts)
		mediaEntries = append(mediaEntries, &MappedMedia{
			MatchedMedia:       matchedMedia,
			TorrentInformation: torrentInformation,
			Size:               totalSize,
		})
	}
	return mediaEntries
}

func getTorrentInformationFromParts(parts []common.MatchedMediaPart) (int64, TorrentInformation) {
	torrentInformation := &TorrentInformation{
		Status:      TorrentStatusPresent,
		RatioStatus: "",
		Ratio:       -1,
		MinRatio:    -1,
		AgeStatus:   "",
		Age:         -1,
		MinAge:      -1,
	}
	totalSize := int64(0)
	missingTorrents := 0
	for _, part := range parts {
		totalSize += part.Size

		ratioStatus, ageStatus := getTrackerRequirementsStatus(part.TorrentFinding, part.Tracker)
		torrentInformation.RatioStatus = getNewTorrentAttributeStatus(torrentInformation.RatioStatus, ratioStatus)
		torrentInformation.AgeStatus = getNewTorrentAttributeStatus(torrentInformation.AgeStatus, ageStatus)

		if part.TorrentFinding == nil {
			missingTorrents++
		}

		if len(parts) == 1 {
			if part.TorrentFinding != nil {
				torrentInformation.Ratio = part.TorrentFinding.Ratio
				torrentInformation.Age = now().Sub(part.TorrentFinding.Added)
			}

			if part.Tracker.IsValid() {
				torrentInformation.MinRatio = part.Tracker.MinRatio
				torrentInformation.MinAge = part.Tracker.MinAge
			}
		}
	}
	if missingTorrents == len(parts) {
		torrentInformation.Status = TorrentStatusMissing
	} else if missingTorrents != 0 {
		torrentInformation.Status = TorrentStatusIncomplete
	}
	return totalSize, *torrentInformation
}

func getNewTorrentAttributeStatus(torrentInformationStatus TorrentAttributeStatus, torrentStatus TorrentAttributeStatus) TorrentAttributeStatus {
	if torrentInformationStatus == TorrentAttributeStatusUnknown {
		return torrentInformationStatus
	} else if torrentInformationStatus == TorrentAttributeStatusFulfilled && torrentStatus != torrentInformationStatus {
		return torrentStatus
	}
	return torrentStatus
}

func getTrackerRequirementsStatus(torrentFinding *common.TorrentEntry, tracker common.Tracker) (ratioStatus TorrentAttributeStatus, ageStatus TorrentAttributeStatus) {
	if torrentFinding == nil || !tracker.IsValid() {
		return TorrentAttributeStatusUnknown, TorrentAttributeStatusUnknown
	}
	if torrentFinding.Ratio >= tracker.MinRatio {
		ratioStatus = TorrentAttributeStatusFulfilled
	} else {
		ratioStatus = TorrentAttributeStatusPending
	}
	if now().After(torrentFinding.Added.Add(tracker.MinAge)) {
		ageStatus = TorrentAttributeStatusFulfilled
	} else {
		ageStatus = TorrentAttributeStatusPending
	}
	return ratioStatus, ageStatus
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
	handler.serveMediaSeriesEntry(writer, request, seriesId, collapsed)
}

func (handler *handler) serveMediaSeriesEntry(writer http.ResponseWriter, request *http.Request, seriesId int64, collapsed bool) {
	logger := getRequestLogger(request)
	media, err := handler.manager.GetMatchedMediaBySeriesId(seriesId)
	if errors.Is(err, common.ErrMediaNotFound) {
		// write 200 because of HTMX request
		writer.WriteHeader(http.StatusOK)
		return
	} else if err != nil {
		logger.Error(err.Error())
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	mappedMedias := handler.getMatchedMediaList([]common.MatchedMedia{media})
	if len(mappedMedias) == 0 {
		http.NotFound(writer, request)
		return
	}
	mappedMedia := mappedMedias[0]
	if collapsed {
		mappedMedia.Parts = []common.MatchedMediaPart{}
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
		season.Size, season.TorrentInformation = getTorrentInformationFromParts(season.Parts)
	}
	mappedMedia.Parts = partsNoSeason
	return mappedSeries
}

func (handler *handler) getMediaDeletionHandler(mediaType common.MediaType) http.HandlerFunc {
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
		logger = logger.With("mediaType", common.MediaTypeSeries, "season", season, "id", id)
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
