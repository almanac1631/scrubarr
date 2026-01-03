package media

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"path/filepath"

	"github.com/almanac1631/scrubarr/pkg/common"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

type SonarrRetriever struct {
	seriesCache             []*sonarr.Series
	seriesEpisodeFilesCache map[int64][]*sonarr.EpisodeFile
	client                  *sonarr.Sonarr
	appUrl                  string
}

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
