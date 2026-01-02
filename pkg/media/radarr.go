package media

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/almanac1631/scrubarr/pkg/common"
	"golift.io/starr"
	"golift.io/starr/radarr"
)

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

func (r *RadarrRetriever) GetMovies() ([]common.Media, error) {
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
				Type:  common.MediaTypeMovie,
				Title: movie.Title,
				Url:   path.Join(r.appUrl, fmt.Sprintf("/movie/%d", movie.TmdbID)),
				Added: movie.Added,
			},
			Parts: []common.MediaPart{
				{
					OriginalFilePath: filepath.Base(movie.MovieFile.OriginalFilePath),
					Size:             movie.SizeOnDisk,
				},
			},
		})
	}
	return mappedMovies, nil
}
