package arr_apps

import (
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"golift.io/starr"
	"golift.io/starr/radarr"
	"maps"
	"path"
	"slices"
)

var _ common.EntryRetriever = (*RadarrMediaRetriever)(nil)

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

func (r RadarrMediaRetriever) RetrieveEntries() (common.RetrieverEntries, error) {
	movieList, err := r.client.GetMovie(0)
	if err != nil {
		return nil, fmt.Errorf("could not get movie list: %w", err)
	}
	mediaEntryList := common.RetrieverEntries{}
	for _, movie := range movieList {
		movieFileList, err := r.client.GetMovieFile(movie.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get movie file list for movie %d: %w", movie.ID, err)
		}
		maps.Copy(mediaEntryList, r.getEntriesFromMovieFileList(movie, movieFileList))
	}
	return mediaEntryList, nil
}

func (r RadarrMediaRetriever) getEntriesFromMovieFileList(movie *radarr.Movie, movieList []*radarr.MovieFile) map[common.EntryName]common.Entry {
	mediaEntryList := map[common.EntryName]common.Entry{}
	for _, movieFile := range movieList {
		fileExtensions := path.Ext(movieFile.Path)
		if !slices.Contains(r.allowedFileEndings, fileExtensions) {
			continue
		}
		name := common.EntryName(path.Base(movieFile.Path))
		mediaEntryList[name] = common.Entry{
			Name: name,
			AdditionalData: ArrAppEntry{
				ID:            movieFile.ID,
				Type:          MediaTypeMovie,
				ParentName:    movie.Title,
				Monitored:     movie.Monitored,
				MediaFilePath: movieFile.Path,
				DateAdded:     movieFile.DateAdded,
				Size:          movieFile.Size,
			},
		}
	}
	return mediaEntryList
}

func (r RadarrMediaRetriever) DeleteEntry(id any) error {
	movieFileID, ok := id.(int64)
	if !ok {
		return fmt.Errorf("could not convert id to int64: %v", id)
	}
	err := r.client.DeleteMovieFiles(movieFileID)
	if err != nil {
		return fmt.Errorf("could not delete movie file with id %d: %w", movieFileID, err)
	}
	return nil
}
