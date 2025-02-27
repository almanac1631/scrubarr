package sqlite

import (
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/arr_apps"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/folder_scanning"
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
		{"test get added date from FileEntry", args{common.Entry{
			Name: "test",
			AdditionalData: folder_scanning.FileEntry{
				DateModified: utils.ParseTime("2025-01-24T16:58:17Z"),
			},
		}}, want{utils.ParseTimePtr("2025-01-24T16:58:17Z"), nil}},
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
