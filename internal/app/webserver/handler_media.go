package webserver

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/almanac1631/scrubarr/internal/utils"
	"github.com/almanac1631/scrubarr/pkg/common"
)

type MappedMovie struct {
	common.MovieInfo
	NextPage int
}

func (handler *handler) handleMediaEndpoint(writer http.ResponseWriter, request *http.Request) {
	if !utils.IsHTMXRequest(request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := handler.templateCache["media.gohtml"].ExecuteTemplate(writer, "index", nil); isErrAndNoBrokenPipe(err) {
			slog.Error("failed to execute template", "err", err)
			return
		}
		return
	}
	writer.WriteHeader(http.StatusNotFound)
	_, _ = writer.Write([]byte("404 Not Found"))
}

func (handler *handler) handleMediaEntriesEndpoint(writer http.ResponseWriter, request *http.Request) {
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
	movies, hasNext, err := handler.manager.GetMovieInfos(page)
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
	if err = handler.templateCache["media.gohtml"].ExecuteTemplate(writer, "media_entries", mediaEntries); isErrAndNoBrokenPipe(err) {
		slog.Error("failed to execute template", "err", err)
		return
	}
	return
}
