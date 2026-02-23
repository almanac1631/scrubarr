package inventory

import (
	"reflect"
	"testing"

	"github.com/almanac1631/scrubarr/internal/app/webserver"
	"github.com/almanac1631/scrubarr/pkg/domain"
	"github.com/stretchr/testify/require"
)

var wantErrMalformedMediaId = func(t require.TestingT, err error, i ...interface{}) {
	require.ErrorIs(t, webserver.ErrMalformedMediaId, err)
}

func Test_parseMediaId(t *testing.T) {
	type args struct {
		rawId string
	}
	tests := []struct {
		name    string
		args    args
		want    mediaId
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "movie id",
			args: args{
				rawId: "movie-1337",
			},
			want: mediaId{
				MediaType: domain.MediaTypeMovie,
				Id:        1337,
			},
			wantErr: require.NoError,
		},
		{
			name: "series id",
			args: args{
				rawId: "series-10",
			},
			want: mediaId{
				MediaType: domain.MediaTypeSeries,
				Id:        10,
			},
			wantErr: require.NoError,
		},
		{
			name: "error on invalid media type",
			args: args{
				rawId: "film-1337",
			},
			wantErr: wantErrMalformedMediaId,
		},
		{
			name: "error on invalid id",
			args: args{
				rawId: "movie-10a",
			},
			wantErr: wantErrMalformedMediaId,
		},
		{
			name: "specific file id of movie",
			args: args{
				rawId: "movie-10-7",
			},
			want: mediaId{
				MediaType: domain.MediaTypeMovie,
				Id:        10,
				FileId:    7,
			},
			wantErr: require.NoError,
		},
		{
			name: "error on invalid file id",
			args: args{
				rawId: "movie-10-7a",
			},
			wantErr: wantErrMalformedMediaId,
		},
		{
			name: "specific season of series",
			args: args{
				rawId: "series-1337-s-2",
			},
			want: mediaId{
				MediaType: domain.MediaTypeSeries,
				Id:        1337,
				Season:    2,
			},
			wantErr: require.NoError,
		},
		{
			name: "error on invalid season",
			args: args{
				rawId: "series-1337-s-2a",
			},
			wantErr: wantErrMalformedMediaId,
		},
		{
			name: "error on invalid season prefix",
			args: args{
				rawId: "series-1337-a-2",
			},
			wantErr: wantErrMalformedMediaId,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMediaId(tt.args.rawId)
			tt.wantErr(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_mediaId_String(t *testing.T) {
	type fields struct {
		MediaType domain.MediaType
		Id        int64
		FileId    int64
		Season    int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "generate basic id",
			fields: fields{
				MediaType: domain.MediaTypeMovie,
				Id:        1337,
			},
			want: "movie-1337",
		},
		{
			name: "generate file id",
			fields: fields{
				MediaType: domain.MediaTypeSeries,
				Id:        10,
				FileId:    7,
			},
			want: "series-10-7",
		},
		{
			name: "generate season id",
			fields: fields{
				MediaType: domain.MediaTypeSeries,
				Id:        10,
				Season:    2,
			},
			want: "series-10-s-2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mediaId{
				MediaType: tt.fields.MediaType,
				Id:        tt.fields.Id,
				FileId:    tt.fields.FileId,
				Season:    tt.fields.Season,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mediaId_getMatchingLinkedMediaIndexes(t *testing.T) {
	file1 := LinkedMediaFile{
		MediaFile: domain.MediaFile{
			Id:     1,
			Season: 1,
		},
	}
	file2 := LinkedMediaFile{
		MediaFile: domain.MediaFile{
			Id:     2,
			Season: 1,
		},
	}
	file3 := LinkedMediaFile{
		MediaFile: domain.MediaFile{
			Id: 3,
		},
	}
	type fields struct {
		MediaType domain.MediaType
		Id        int64
		FileId    int64
		Season    int
	}
	type args struct {
		linkedMediaFiles []LinkedMediaFile
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		{
			name: "whole entry",
			fields: fields{
				MediaType: domain.MediaTypeMovie,
				Id:        1337,
			},
			args: args{
				linkedMediaFiles: []LinkedMediaFile{
					file1,
					file2,
					file3,
				},
			},
			want: []int{
				0, 1, 2,
			},
		},
		{
			name: "specific file id",
			fields: fields{
				MediaType: domain.MediaTypeMovie,
				Id:        1337,
				FileId:    2,
			},
			args: args{
				linkedMediaFiles: []LinkedMediaFile{
					file1,
					file2,
					file3,
				},
			},
			want: []int{
				1,
			},
		},
		{
			name: "specific season",
			fields: fields{
				MediaType: domain.MediaTypeSeries,
				Id:        1337,
				Season:    1,
			},
			args: args{
				linkedMediaFiles: []LinkedMediaFile{
					file1,
					file2,
					file3,
				},
			},
			want: []int{
				0,
				1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mediaId{
				MediaType: tt.fields.MediaType,
				Id:        tt.fields.Id,
				FileId:    tt.fields.FileId,
				Season:    tt.fields.Season,
			}
			if got := m.getMatchingLinkedMediaIndexes(tt.args.linkedMediaFiles); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMatchingLinkedMediaIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}
