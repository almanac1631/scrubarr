package media

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/almanac1631/scrubarr/pkg/common"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

type SonarrRetriever struct {
	client *sonarr.Sonarr
	appUrl string
}

func NewSonarrRetriever(appUrl string, apiKey string) (*SonarrRetriever, error) {
	config := starr.New(apiKey, appUrl, 0)
	client := sonarr.New(config)
	_, err := client.GetSystemStatus()
	if err != nil {
		return nil, fmt.Errorf("could not get sonarr system status: %w", err)
	}
	return &SonarrRetriever{client, appUrl}, nil
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
				Season:           seriesEpisodeFile.SeasonNumber,
				OriginalFilePath: filepath.Base(seriesEpisodeFile.RelativePath),
				Size:             seriesEpisodeFile.Size,
			})
		}
		media := common.Media{
			MediaMetadata: common.MediaMetadata{
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
