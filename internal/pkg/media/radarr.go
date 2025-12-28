package media

import (
	"fmt"
	"time"

	"golift.io/starr"
	"golift.io/starr/radarr"
)

type RadarrRetriever struct {
	client *radarr.Radarr
}

type Movie struct {
	Title            string
	Size             int64
	Added            time.Time
	OriginalFilePath string
}

func NewRadarrRetriever(hostname string, apiKey string) (*RadarrRetriever, error) {
	starrConfig := starr.New(apiKey, hostname, 0)
	client := radarr.New(starrConfig)
	_, err := client.GetSystemStatus()
	if err != nil {
		return nil, fmt.Errorf("could not get radarr system status: %w", err)
	}
	return &RadarrRetriever{client}, nil
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
		})
	}
	return mappedMovies, nil
}
