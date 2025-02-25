package arr_apps

import (
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"golift.io/starr"
	"golift.io/starr/sonarr"
	"path"
	"slices"
)

type SonarrMediaRetriever struct {
	client             *sonarr.Sonarr
	allowedFileEndings []string
}

func NewSonarrMediaRetriever(allowedFileEndings []string, hostname string, apiKey string) (*SonarrMediaRetriever, error) {
	starrConfig := starr.New(apiKey, hostname, 0)
	client := sonarr.New(starrConfig)
	_, err := client.GetSystemStatus()
	if err != nil {
		return nil, fmt.Errorf("could not get sonarr system status: %w", err)
	}
	return &SonarrMediaRetriever{client, allowedFileEndings}, nil
}

func (s *SonarrMediaRetriever) RetrieveEntries() (common.RetrieverEntries, error) {
	seriesList, err := s.client.GetAllSeries()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve sonarr series list: %w", err)
	}
	mediaEntryList := common.RetrieverEntries{}
	for _, series := range seriesList {
		fileList, err := s.client.GetSeriesEpisodeFiles(series.ID)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve files for series (id: %d): %w", series.ID, err)
		}
		for _, file := range fileList {
			if !slices.Contains(s.allowedFileEndings, path.Ext(file.Path)) {
				continue
			}
			mediaEntry := s.parseSeriesEpisodeFile(series, file)
			mediaEntryList[mediaEntry.Name] = mediaEntry
		}
	}
	return mediaEntryList, nil
}

func (s *SonarrMediaRetriever) parseSeriesEpisodeFile(series *sonarr.Series, episodeFile *sonarr.EpisodeFile) common.Entry {
	monitored := s.isSeasonMonitored(series, episodeFile.SeasonNumber)
	name := path.Base(episodeFile.Path)
	return common.Entry{
		Name: common.EntryName(name),
		AdditionalData: ArrAppEntry{
			Type:          MediaTypeSeries,
			ParentName:    series.Title,
			Monitored:     monitored,
			MediaFilePath: episodeFile.Path,
		},
	}
}

func (s *SonarrMediaRetriever) isSeasonMonitored(series *sonarr.Series, seasonNumber int) bool {
	for _, season := range series.Seasons {
		if season.SeasonNumber == seasonNumber {
			return season.Monitored
		}
	}
	return false
}
