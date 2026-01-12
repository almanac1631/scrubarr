package media

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"path"
	"path/filepath"
	"slices"

	"github.com/almanac1631/scrubarr/pkg/common"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

var _ common.MediaRetriever = (*SonarrRetriever)(nil)

type SonarrRetriever struct {
	client *sonarr.Sonarr
	appUrl string
	dryRun bool
}

const (
	sonarrEpisodeFileEndpoint           = sonarr.APIver + "/episodeFile"
	sonarrEpisodeFileBulkDeleteEndpoint = sonarrEpisodeFileEndpoint + "/bulk"
)

func NewSonarrRetriever(appUrl string, apiKey string, dryRun bool) (*SonarrRetriever, error) {
	config := starr.New(apiKey, appUrl, 0)
	client := sonarr.New(config)
	_, err := client.GetSystemStatus()
	if err != nil {
		return nil, fmt.Errorf("could not get sonarr system status: %w", err)
	}
	return &SonarrRetriever{client, appUrl, dryRun}, nil
}

func (r *SonarrRetriever) GetMedia() ([]common.Media, error) {
	seriesList, err := r.client.GetAllSeries()
	if err != nil {
		return nil, fmt.Errorf("could not get sonarr series: %w", err)
	}
	mediaList := make([]common.Media, 0)
	for _, series := range seriesList {
		if series.Statistics.SizeOnDisk == 0 {
			continue
		}
		seriesEpisodeFiles, err := r.client.GetSeriesEpisodeFiles(series.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get series episode files: %w", err)
		}
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

func (r *SonarrRetriever) DeleteMediaFiles(fileIds []int64, stopParentMonitoring bool) error {
	episodeFiles, err := r.getEpisodeFiles(fileIds)
	if err != nil {
		return fmt.Errorf("could not get sonarr episode files: %w", err)
	}
	seriesSeasonMap := make(map[int64][]int)
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
	if err = r.deleteEpisodeFiles(fileIds); err != nil {
		return err
	}
	monitoringUpdateErrors := make([]error, 0)
	if stopParentMonitoring {
		for seriesId, seasonList := range seriesSeasonMap {
			if err = r.stopMonitoringSeasons(seriesId, seasonList); err != nil {
				monitoringUpdateErrors = append(monitoringUpdateErrors, err)
			}
		}
	}
	if len(monitoringUpdateErrors) > 0 {
		return errors.Join(monitoringUpdateErrors...)
	}
	return nil
}

func (r *SonarrRetriever) getEpisodeFiles(fileIds []int64) ([]*sonarr.EpisodeFile, error) {
	req := starr.Request{URI: sonarrEpisodeFileEndpoint, Query: make(url.Values)}
	for _, efID := range fileIds {
		req.Query.Add("episodeFileIds", starr.Str(efID))
	}

	var output []*sonarr.EpisodeFile
	if err := r.client.GetInto(context.Background(), req, &output); err != nil {
		return nil, fmt.Errorf("api.Get(%s): %w", &req, err)
	}

	return output, nil
}

func (r *SonarrRetriever) SupportedMediaType() common.MediaType {
	return common.MediaTypeSeries
}

func (r *SonarrRetriever) deleteEpisodeFiles(fileIds []int64) error {
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
	if r.dryRun {
		slog.Info("[DRY RUN] Skipping delete Sonarr episode files.", "fileIds", fileIds)
		return nil
	}
	if err = r.client.DeleteAny(context.Background(), req); err != nil {
		return fmt.Errorf("could not bulk delete episode files from sonarr: %w",
			fmt.Errorf("api.Delete(%s): %w", &req, err))
	}
	return nil
}

func (r *SonarrRetriever) stopMonitoringSeasons(seriesId int64, seasonNumbers []int) error {
	seasons := make([]*sonarr.Season, len(seasonNumbers))
	for _, season := range seasonNumbers {
		seasons = append(seasons, &sonarr.Season{
			Monitored:    false,
			SeasonNumber: season,
		})
	}
	if r.dryRun {
		slog.Info("[DRY RUN] Skipping Sonarr stop monitoring season.", "seriesId", seriesId, "seasons", seasonNumbers)
		return nil
	}
	_, err := r.client.UpdateSeries(&sonarr.AddSeriesInput{
		ID:      seriesId,
		Seasons: seasons,
	}, false)
	if err != nil {
		return fmt.Errorf("could not update monitoring status of series %d: %w", seriesId, err)
	}
	return nil
}
