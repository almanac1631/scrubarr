package arr_apps

import (
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval"
	"golift.io/starr"
	"golift.io/starr/radarr"
	"maps"
	"path"
	"slices"
)

type RadarrMediaRetriever struct {
	client             *radarr.Radarr
	allowedFileEndings []string
}

func NewRadarrMediaRetriever(allowedFileEndings []string, hostname string, apiKey string) (*RadarrMediaRetriever, error) {
	starrConfig := starr.New(apiKey, hostname, 0)
	client := radarr.New(starrConfig)
	_, err := client.GetSystemStatus()
	if err != nil {
		return nil, fmt.Errorf("could not get radarr system status: %w", err)
	}
	return &RadarrMediaRetriever{client, allowedFileEndings}, nil
}

func (r RadarrMediaRetriever) RetrieveEntries() (map[retrieval.EntryName]retrieval.Entry, error) {
	movieList, err := r.client.GetMovie(0)
	if err != nil {
		return nil, fmt.Errorf("could not get movie list: %w", err)
	}
	mediaEntryList := map[retrieval.EntryName]retrieval.Entry{}
	for _, movie := range movieList {
		movieFileList, err := r.client.GetMovieFile(movie.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get movie file list for movie %d: %w", movie.ID, err)
		}
		maps.Copy(mediaEntryList, r.getEntriesFromMovieFileList(movie, movieFileList))
	}
	return mediaEntryList, nil
}

func (r RadarrMediaRetriever) getEntriesFromMovieFileList(movie *radarr.Movie, movieList []*radarr.MovieFile) map[retrieval.EntryName]retrieval.Entry {
	mediaEntryList := map[retrieval.EntryName]retrieval.Entry{}
	for _, movieFile := range movieList {
		fileExtensions := path.Ext(movieFile.Path)
		if !slices.Contains(r.allowedFileEndings, fileExtensions) {
			continue
		}
		name := retrieval.EntryName(path.Base(movieFile.Path))
		mediaEntryList[name] = retrieval.Entry{
			Name: name,
			AdditionalData: ArrAppEntry{
				Type:          MediaTypeMovie,
				ParentName:    movie.Title,
				Monitored:     movie.Monitored,
				MediaFilePath: movieFile.Path,
			},
		}
	}
	return mediaEntryList
}
