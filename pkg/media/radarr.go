package media

import (
	"fmt"
	"path"
	"time"

	"golift.io/starr"
	"golift.io/starr/radarr"
)

type RadarrRetriever struct {
	client *radarr.Radarr
	appUrl string
}

type Movie struct {
	Title            string
	Size             int64
	Added            time.Time
	OriginalFilePath string
	Url              string
}

func NewRadarrRetriever(appUrl string, apiKey string) (*RadarrRetriever, error) {
	starrConfig := starr.New(apiKey, appUrl, 0)
	client := radarr.New(starrConfig)
	_, err := client.GetSystemStatus()
	if err != nil {
		return nil, fmt.Errorf("could not get radarr system status: %w", err)
	}
	return &RadarrRetriever{client, appUrl}, nil
}

func (r RadarrRetriever) GetMovies() ([]Movie, error) {
	movies, err := r.client.GetMovie(&radarr.GetMovie{
		TMDBID:             0,
		ExcludeLocalCovers: true,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get movie list: %w", err)
	}
	var mappedMovies []Movie
	for _, movie := range movies {
		if !movie.HasFile {
			continue
		}
		mappedMovies = append(mappedMovies, Movie{
			Title:            movie.Title,
			Size:             movie.SizeOnDisk,
			Added:            movie.Added,
			OriginalFilePath: movie.MovieFile.OriginalFilePath,
			Url:              path.Join(r.appUrl, fmt.Sprintf("/movie/%d", movie.TmdbID)),
		})
	}
	return mappedMovies, nil
}
