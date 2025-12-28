package webserver

import (
	"log/slog"
	"net/http"

	"github.com/almanac1631/scrubarr/internal/utils"
)

func (handler *handler) handleMediaEndpoint(writer http.ResponseWriter, request *http.Request) {
	if !utils.IsHTMXRequest(request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		movies, err := handler.radarrRetriever.GetMovies()
		if err != nil {
			slog.Error("failed to get movies from radarr", "err", err)
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte("500 Internal Server Error"))
			return
		}
		mappedMovies := make([]MovieMapped, 0, len(movies))
		for _, movie := range movies {
			exists, err := handler.delugeRetriever.SearchForMovie(movie.OriginalFilePath)
			if err != nil {
				slog.Error("failed to search for movie in deluge", "err", err)
				writer.WriteHeader(http.StatusInternalServerError)
				_, _ = writer.Write([]byte("500 Internal Server Error"))
				return
			}
			if !exists {
				exists, err = handler.rtorrentRetriever.SearchForMovie(movie.OriginalFilePath)
				if err != nil {
					slog.Error("failed to search for movie in rtorrent", "err", err)
					writer.WriteHeader(http.StatusInternalServerError)
					_, _ = writer.Write([]byte("500 Internal Server Error"))
					return
				}
			}
			mappedMovies = append(mappedMovies, MovieMapped{
				Movie:                 movie,
				ExistsInTorrentClient: exists,
			})
		}
		if err := handler.templateCache["media.gohtml"].ExecuteTemplate(writer, "index", mappedMovies); isErrAndNoBrokenPipe(err) {
			slog.Error("failed to execute template", "err", err)
			return
		}
		return
	}
	writer.WriteHeader(http.StatusNotFound)
	_, _ = writer.Write([]byte("404 Not Found"))
}
