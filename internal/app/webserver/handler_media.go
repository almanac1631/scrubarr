package webserver

import (
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/almanac1631/scrubarr/internal/utils"
	"github.com/almanac1631/scrubarr/pkg/common"
)

type mediaEndpointData struct {
	Movies   []MappedMovie
	SortInfo common.SortInfo
}

type MappedMovie struct {
	common.MovieInfo
	NextPage int
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
	movies, hasNext, err := handler.manager.GetMovieInfos(page, sortInfo)
	if err != nil {
		slog.Error("failed to get movie mapping", "err", err)
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("500 Internal Server Error"))
		return
	}
	mediaEntries := make([]MappedMovie, 0, len(movies))
	for i, movie := range movies {
		mappedMovie := &MappedMovie{
			MovieInfo: movie,
			NextPage:  -1,
		}
		if hasNext && i == len(movies)-1 {
			mappedMovie.NextPage = page + 1
		}
		mediaEntries = append(mediaEntries, *mappedMovie)
	}
	if err = handler.ExecuteSubTemplate(writer, "media.gohtml", "media_entries", mediaEndpointData{
		Movies:   mediaEntries,
		SortInfo: sortInfo,
	}); err != nil {
		slog.Error(err.Error())
		return
	}
	return
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
