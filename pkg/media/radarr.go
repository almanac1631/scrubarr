package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"

	"github.com/almanac1631/scrubarr/pkg/common"
	"golift.io/starr"
	"golift.io/starr/radarr"
)

var _ common.MediaRetriever = (*RadarrRetriever)(nil)

type RadarrRetriever struct {
	moviesCache []*radarr.Movie
	client      *radarr.Radarr
	appUrl      string
}

func NewRadarrRetriever(appUrl string, apiKey string) (*RadarrRetriever, error) {
	starrConfig := starr.New(apiKey, appUrl, 0)
	client := radarr.New(starrConfig)
	_, err := client.GetSystemStatus()
	if err != nil {
		return nil, fmt.Errorf("could not get radarr system status: %w", err)
	}
	return &RadarrRetriever{nil, client, appUrl}, nil
}

func (r *RadarrRetriever) RefreshCache() error {
	var err error
	r.moviesCache, err = r.client.GetMovie(&radarr.GetMovie{
		TMDBID:             0,
		ExcludeLocalCovers: true,
	})
	if err != nil {
		return fmt.Errorf("could not get radarr movies: %w", err)
	}
	return nil
}

func (r *RadarrRetriever) SaveCache(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(r.moviesCache)
}

func (r *RadarrRetriever) LoadCache(reader io.ReadSeeker) error {
	r.moviesCache = []*radarr.Movie{}
	return json.NewDecoder(reader).Decode(&r.moviesCache)
}

func (r *RadarrRetriever) GetMedia() ([]common.Media, error) {
	if r.moviesCache == nil {
		if err := r.RefreshCache(); err != nil {
			return nil, err
		}
	}
	var mappedMovies []common.Media
	for _, movie := range r.moviesCache {
		if !movie.HasFile {
			continue
		}
		mappedMovies = append(mappedMovies, common.Media{
			MediaMetadata: common.MediaMetadata{
				Id:    movie.ID,
				Type:  common.MediaTypeMovie,
				Title: movie.Title,
				Url:   path.Join(r.appUrl, fmt.Sprintf("/movie/%d", movie.TmdbID)),
				Added: movie.Added,
			},
			Parts: []common.MediaPart{
				{
					Id:               movie.MovieFile.ID,
					OriginalFilePath: filepath.Base(movie.MovieFile.OriginalFilePath),
					Size:             movie.SizeOnDisk,
				},
			},
		})
	}
	return mappedMovies, nil
}
func (r *RadarrRetriever) DeleteMedia(id int64) error {
	if err := r.client.DeleteMovie(id, true, false); err != nil {
		return fmt.Errorf("could not delete movie %d from radarr: %w", id, err)
	}
	return nil
}

func (r *RadarrRetriever) DeleteMediaFiles(fileIds []int64, stopParentMonitoring bool) error {
	movieFiles, err := r.client.GetMovieFiles(fileIds)
	if err != nil {
		return fmt.Errorf("could not get radarr movie files: %w", err)
	}
	var movies map[int64]struct{}
	if stopParentMonitoring {
		for _, movieFile := range movieFiles {
			if _, ok := movies[movieFile.MovieID]; ok {
				continue
			}
			movies[movieFile.MovieID] = struct{}{}
		}
	}
	if err = r.client.DeleteMovieFiles(fileIds...); err != nil {
		return fmt.Errorf("could not bulk delete movie files: %w", err)
	}
	monitoringUpdateErrors := make([]error, 0)
	if stopParentMonitoring {
		for movieId, _ := range movies {
			if _, err = r.client.UpdateMovie(movieId, &radarr.Movie{
				ID:        movieId,
				Monitored: false,
			}, false); err != nil {
				monitoringUpdateErrors = append(monitoringUpdateErrors, err)
			}
		}
	}
	if len(monitoringUpdateErrors) > 0 {
		return errors.Join(monitoringUpdateErrors...)
	}
	return nil
}

func (r *RadarrRetriever) SupportedMediaType() common.MediaType {
	return common.MediaTypeMovie
}
