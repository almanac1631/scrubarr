package arr_apps

import (
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval"
	"golift.io/starr/radarr"
	"reflect"
	"testing"
)

func TestRadarrMediaRetriever_getEntriesFromMovieFileList(t *testing.T) {
	type fields struct {
		allowedFileEndings []string
	}
	type args struct {
		movie     *radarr.Movie
		movieList []*radarr.MovieFile
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[retrieval.EntryName]retrieval.Entry
	}{
		{
			"can parse a single movie file",
			fields{[]string{".mkv"}},
			args{
				&radarr.Movie{
					ID:        1337,
					Title:     "Some Cool Film yo",
					Monitored: true,
				},
				[]*radarr.MovieFile{
					{
						ID:   13371,
						Path: "some/film/dir/Some Cool Film Name.mkv",
					},
				}},
			map[retrieval.EntryName]retrieval.Entry{
				"Some Cool Film Name.mkv": {
					Name: "Some Cool Film Name.mkv",
					AdditionalData: ArrAppEntry{
						Type:          MediaTypeMovie,
						ParentName:    "Some Cool Film yo",
						Monitored:     true,
						MediaFilePath: "some/film/dir/Some Cool Film Name.mkv",
					},
				},
			},
		},
		{
			"can filter invalid file extensions",
			fields{[]string{".mp4"}},
			args{
				&radarr.Movie{
					ID:        1337,
					Title:     "Some Cool Film yo",
					Monitored: true,
				},
				[]*radarr.MovieFile{
					{
						ID:   13371,
						Path: "some/film/dir/Some Cool Film Name.mkv",
					},
				}},
			map[retrieval.EntryName]retrieval.Entry{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RadarrMediaRetriever{
				allowedFileEndings: tt.fields.allowedFileEndings,
			}
			if got := r.getEntriesFromMovieFileList(tt.args.movie, tt.args.movieList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEntriesFromMovieFileList() = %v, want %v", got, tt.want)
			}
		})
	}
}
