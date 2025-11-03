package arr_apps

import (
	"reflect"
	"testing"

	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/almanac1631/scrubarr/internal/pkg/utils"
	"golift.io/starr/sonarr"
)

func TestSonarrMediaRetriever_parseSeriesEpisodeFile(t *testing.T) {
	type fields struct {
		client *sonarr.Sonarr
	}
	type args struct {
		series      *sonarr.Series
		episodeFile *sonarr.EpisodeFile
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   common.Entry
	}{
		{
			"can parse a monitored episode file", fields{nil}, args{
				&sonarr.Series{
					ID: 2199,
					Seasons: []*sonarr.Season{
						{true, 1, nil},
						{true, 2, nil},
						{false, 3, nil},
					},
					Title: "Some series!",
				},
				&sonarr.EpisodeFile{
					ID:           21991,
					SeasonNumber: 2,
					RelativePath: "Season 1/Some Episode.mkv",
					Path:         "/home/myuser/media/downloads/Some Episode.mkv",
					DateAdded:    utils.ParseTime("2025-02-18T13:29:48Z"),
					Size:         21849284229329,
				},
			}, common.Entry{
				Name: "Some Episode.mkv",
				AdditionalData: ArrAppEntry{
					ID:            21991,
					Type:          MediaTypeSeries,
					ParentName:    "Some series!",
					ParentId:      "2199/2",
					Monitored:     true,
					MediaFilePath: "/home/myuser/media/downloads/Some Episode.mkv",
					DateAdded:     utils.ParseTime("2025-02-18T13:29:48Z"),
					Size:          21849284229329,
				},
			},
		},
		{
			"can parse an unmonitored episode file", fields{nil}, args{
				&sonarr.Series{
					ID: 8429,
					Seasons: []*sonarr.Season{
						{true, 1, nil},
						{true, 2, nil},
						{false, 3, nil},
					},
					Title: "Some series!",
				},
				&sonarr.EpisodeFile{
					ID:           84291,
					SeasonNumber: 3,
					RelativePath: "Some Episode.mkv",
					Path:         "/home/myuser/media/downloads/Some Episode.mkv",
					DateAdded:    utils.ParseTime("2025-02-09T08:19:04Z"),
					Size:         21849284229329,
				},
			}, common.Entry{
				Name: "Some Episode.mkv",
				AdditionalData: ArrAppEntry{
					ID:            84291,
					Type:          MediaTypeSeries,
					ParentName:    "Some series!",
					ParentId:      "8429/3",
					Monitored:     false,
					MediaFilePath: "/home/myuser/media/downloads/Some Episode.mkv",
					DateAdded:     utils.ParseTime("2025-02-09T08:19:04Z"),
					Size:          21849284229329,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SonarrMediaRetriever{
				client: tt.fields.client,
			}
			got := s.parseSeriesEpisodeFile(tt.args.series, tt.args.episodeFile)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSeriesEpisodeFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
