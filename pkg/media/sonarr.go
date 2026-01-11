package media

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"slices"

	"github.com/almanac1631/scrubarr/pkg/common"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

var _ common.MediaRetriever = (*SonarrRetriever)(nil)

type SonarrRetriever struct {
	seriesCache             []*sonarr.Series
	seriesEpisodeFilesCache map[int64][]*sonarr.EpisodeFile
	client                  *sonarr.Sonarr
	appUrl                  string
}

const sonarrEpisodeFileBulkDeleteEndpoint = sonarr.APIver + "/episodeFile/bulk"

func NewSonarrRetriever(appUrl string, apiKey string) (*SonarrRetriever, error) {
	config := starr.New(apiKey, appUrl, 0)
	client := sonarr.New(config)
	_, err := client.GetSystemStatus()
	if err != nil {
		return nil, fmt.Errorf("could not get sonarr system status: %w", err)
	}
	return &SonarrRetriever{nil, nil, client, appUrl}, nil
}

func (r *SonarrRetriever) RefreshCache() error {
	var err error
	r.seriesCache, err = r.client.GetAllSeries()
	if err != nil {
		return fmt.Errorf("could not get sonarr series: %w", err)
	}
	r.seriesEpisodeFilesCache = make(map[int64][]*sonarr.EpisodeFile)
	for _, series := range r.seriesCache {
		if series.Statistics.SizeOnDisk == 0 {
			continue
		}
		r.seriesEpisodeFilesCache[series.ID], err = r.client.GetSeriesEpisodeFiles(series.ID)
		if err != nil {
			return fmt.Errorf("could not get series episode files: %w", err)
		}
	}
	return nil
}

type cacheWrapper struct {
	Series             []*sonarr.Series
	SeriesEpisodeFiles map[int64][]*sonarr.EpisodeFile
}

func (r *SonarrRetriever) SaveCache(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(cacheWrapper{r.seriesCache, r.seriesEpisodeFilesCache})
}

func (r *SonarrRetriever) LoadCache(reader io.ReadSeeker) error {
	wrapper := cacheWrapper{}
	if err := json.NewDecoder(reader).Decode(&wrapper); err != nil {
		return err
	}
	r.seriesCache = wrapper.Series
	r.seriesEpisodeFilesCache = wrapper.SeriesEpisodeFiles
	return nil
}

func (r *SonarrRetriever) GetMedia() ([]common.Media, error) {
	if r.seriesCache == nil {
		if err := r.RefreshCache(); err != nil {
			return nil, err
		}
	}
	mediaList := make([]common.Media, 0)
	for _, series := range r.seriesCache {
		if series.Statistics.SizeOnDisk == 0 {
			continue
		}
		seriesEpisodeFiles := r.seriesEpisodeFilesCache[series.ID]
		parts := make([]common.MediaPart, 0, len(seriesEpisodeFiles))
		for _, seriesEpisodeFile := range seriesEpisodeFiles {
			parts = append(parts, common.MediaPart{
				Id:               seriesEpisodeFile.ID,
				Season:           seriesEpisodeFile.SeasonNumber,
				OriginalFilePath: filepath.Base(seriesEpisodeFile.RelativePath),
				Size:             seriesEpisodeFile.Size,
			})
		}
		media := common.Media{
			MediaMetadata: common.MediaMetadata{
				Id:    series.ID,
				Type:  common.MediaTypeSeries,
				Title: series.Title,
				Url:   path.Join(r.appUrl, fmt.Sprintf("series/%s", series.TitleSlug)),
				Added: series.Added,
			},
			Parts: parts,
		}
		mediaList = append(mediaList, media)
	}
	return mediaList, nil
}

func (r *SonarrRetriever) DeleteMedia(id int64) error {
	if err := r.client.DeleteSeries(int(id), true, false); err != nil {
		return fmt.Errorf("could not delete series: %d from sonarr %w", id, err)
	}
	return nil
}

func (r *SonarrRetriever) DeleteMediaFiles(fileIds []int64, stopParentMonitoring bool) error {
	episodeFiles, err := r.client.GetEpisodeFiles(fileIds...)
	if err != nil {
		return fmt.Errorf("could not get sonarr episode files: %w", err)
	}
	var seriesSeasonMap map[int64][]int
	if stopParentMonitoring {
		for _, episodeFile := range episodeFiles {
			seasonList, ok := seriesSeasonMap[episodeFile.SeriesID]
			if !ok {
				seasonList = []int{episodeFile.SeasonNumber}
			} else if !slices.Contains(seasonList, episodeFile.SeasonNumber) {
				seasonList = append(seasonList, episodeFile.SeasonNumber)
			}
			seriesSeasonMap[episodeFile.SeriesID] = seasonList
		}
	}
	payload := struct {
		EpisodeFileIds []int64 `json:"episodeFileIds"`
	}{
		EpisodeFileIds: fileIds,
	}
	payloadEncoded, err := json.Marshal(&payload)
	if err != nil {
		return fmt.Errorf("could not encode sonarr episode file bulk delete payload: %w", err)
	}
	req := starr.Request{URI: sonarrEpisodeFileBulkDeleteEndpoint, Body: bytes.NewReader(payloadEncoded)}
	if err = r.client.DeleteAny(context.Background(), req); err != nil {
		return fmt.Errorf("could not bulkd delete episode files from sonarr: %w",
			fmt.Errorf("api.Delete(%s): %w", &req, err))
	}
	monitoringUpdateErrors := make([]error, 0)
	if stopParentMonitoring {
		for seriesId, seasonList := range seriesSeasonMap {
			seasons := make([]*sonarr.Season, len(seasonList))
			for _, season := range seasonList {
				seasons = append(seasons, &sonarr.Season{
					Monitored:    false,
					SeasonNumber: season,
				})
			}
			_, err = r.client.UpdateSeries(&sonarr.AddSeriesInput{
				ID:      seriesId,
				Seasons: seasons,
			}, false)
			if err != nil {
				monitoringUpdateErrors = append(monitoringUpdateErrors,
					fmt.Errorf("could not update monitoring status of series %d: %w", seriesId, err))
			}
		}
	}
	if len(monitoringUpdateErrors) > 0 {
		return errors.Join(monitoringUpdateErrors...)
	}
	return nil
}

func (r *SonarrRetriever) SupportedMediaType() common.MediaType {
	return common.MediaTypeSeries
}
