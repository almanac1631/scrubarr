package arr_apps

import (
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval"
	"golift.io/starr/sonarr"
	"reflect"
	"testing"
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
		want   retrieval.Entry
	}{
		{
			"can parse a monitored episode file", fields{nil}, args{
				&sonarr.Series{
					Seasons: []*sonarr.Season{
						{true, 1, nil},
						{true, 2, nil},
						{false, 3, nil},
					},
					Title: "Some series!",
				},
				&sonarr.EpisodeFile{
					SeasonNumber: 2,
					RelativePath: "Season 1/Some Episode.mkv",
					Path:         "/home/myuser/media/downloads/Some Episode.mkv",
				},
			}, retrieval.Entry{
				Name: "Some Episode.mkv",
				AdditionalData: ArrAppEntry{
					Type:          MediaTypeSeries,
					ParentName:    "Some series!",
					Monitored:     true,
					MediaFilePath: "/home/myuser/media/downloads/Some Episode.mkv",
				},
			},
		},
		{
			"can parse an unmonitored episode file", fields{nil}, args{
				&sonarr.Series{
					Seasons: []*sonarr.Season{
						{true, 1, nil},
						{true, 2, nil},
						{false, 3, nil},
					},
					Title: "Some series!",
				},
				&sonarr.EpisodeFile{
					SeasonNumber: 3,
					RelativePath: "Some Episode.mkv",
					Path:         "/home/myuser/media/downloads/Some Episode.mkv",
				},
			}, retrieval.Entry{
				Name: "Some Episode.mkv",
				AdditionalData: ArrAppEntry{
					Type:          MediaTypeSeries,
					ParentName:    "Some series!",
					Monitored:     false,
					MediaFilePath: "/home/myuser/media/downloads/Some Episode.mkv",
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
