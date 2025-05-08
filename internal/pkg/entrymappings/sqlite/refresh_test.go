package sqlite

import (
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/arr_apps"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/torrent_clients"
	"github.com/almanac1631/scrubarr/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_getAddedDateFromEntry(t *testing.T) {
	type args struct {
		entry common.Entry
	}
	type want struct {
		value *time.Time
		err   error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"test get added date from ArrAppEntry", args{common.Entry{
			Name: "test",
			AdditionalData: arr_apps.ArrAppEntry{
				DateAdded: utils.ParseTime("2025-02-08T13:09:41Z"),
			},
		}}, want{utils.ParseTimePtr("2025-02-08T13:09:41Z"), nil}},
		{"test get added date from ArrAppEntry for empty date", args{common.Entry{
			Name: "test",
			AdditionalData: arr_apps.ArrAppEntry{
				DateAdded: time.Time{},
			},
		}}, want{nil, nil}},
		{"test get added date from TorrentClientEntry", args{common.Entry{
			Name: "test",
			AdditionalData: torrent_clients.TorrentClientEntry{
				DownloadedAt: utils.ParseTime("2024-11-06T08:34:38Z"),
			},
		}}, want{utils.ParseTimePtr("2024-11-06T08:34:38Z"), nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := getDateAddedFromEntry(tt.args.entry)
			assert.Equal(t, tt.want.err, err)
			assert.Equalf(t, tt.want.value, value, "getDateAddedFromEntry(%v)", tt.args.entry)
		})
	}
}

func Test_getSizeFromEntry(t *testing.T) {
	type args struct {
		entry common.Entry
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr assert.ErrorAssertionFunc
	}{
		{"test get size from ArrAppEntry", args{common.Entry{
			Name: "test",
			AdditionalData: arr_apps.ArrAppEntry{
				Size: 123456,
			},
		}}, 123456, assert.NoError},
		{"test get size from TorrentClientEntry", args{common.Entry{
			Name: "test",
			AdditionalData: torrent_clients.TorrentClientEntry{
				FileSizeBytes: 987654,
			},
		}}, 987654, assert.NoError},
		{"test get size from unknown entry", args{common.Entry{
			Name:           "test",
			AdditionalData: 42,
		}}, 0, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSizeFromEntry(tt.args.entry)
			if !tt.wantErr(t, err, fmt.Sprintf("getSizeFromEntry(%v)", tt.args.entry)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getSizeFromEntry(%v)", tt.args.entry)
		})
	}
}
