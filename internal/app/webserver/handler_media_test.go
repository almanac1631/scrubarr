package webserver

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/almanac1631/scrubarr/pkg/common"
)

func Test_getSortInfoFromUrlQuery(t *testing.T) {
	type args struct {
		values url.Values
	}
	tests := []struct {
		name string
		args args
		want common.SortInfo
	}{
		{
			name: "test valid parse from url query",
			args: args{
				values: map[string][]string{
					"sortKey":   {"size"},
					"sortOrder": {"asc"},
				},
			},
			want: common.SortInfo{
				Key:   common.SortKeySize,
				Order: common.SortOrderAsc,
			},
		},
		{
			name: "test invalid parse from url query",
			args: args{
				values: map[string][]string{
					"sortKey":   {"test"},
					"sortOrder": {"abc"},
				},
			},
			want: common.SortInfo{
				Key:   common.SortKeyName,
				Order: common.SortOrderAsc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSortInfoFromUrlQuery(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSortInfoFromUrlQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBundledTorrentInformationFromParts(t *testing.T) {
	now = func() time.Time {
		return time.Date(2026, time.January, 19, 22, 0, 0, 0, time.UTC)
	}
	type args struct {
		parts []common.MatchedEntryPart
	}
	torrentComplete := common.MatchedEntryPart{
		MediaPart: common.MediaPart{Size: 382842},
		TorrentInformation: common.TorrentInformation{
			Client: "sonarr",
			Id:     "abc",
			Status: common.TorrentStatusPresent,
			Tracker: common.Tracker{
				Name:     "test",
				MinRatio: 0.5,
				MinAge:   time.Hour,
			},
			Age:         time.Hour * 2,
			AgeStatus:   common.TorrentAttributeStatusFulfilled,
			Ratio:       1.2,
			RatioStatus: common.TorrentAttributeStatusFulfilled,
		},
	}

	torrentIncompleteRatio := torrentComplete
	torrentIncompleteRatio.TorrentInformation.Ratio = 0.1
	torrentIncompleteRatio.TorrentInformation.RatioStatus = common.TorrentAttributeStatusPending

	torrentIncompleteAge := torrentComplete
	torrentIncompleteAge.TorrentInformation.Age = time.Minute * 30
	torrentIncompleteAge.TorrentInformation.AgeStatus = common.TorrentAttributeStatusPending

	torrentUnknownTracker := torrentComplete
	torrentUnknownTracker.TorrentInformation.Tracker = common.Tracker{}
	torrentUnknownTracker.TorrentInformation.AgeStatus = common.TorrentAttributeStatusUnknown
	torrentUnknownTracker.TorrentInformation.RatioStatus = common.TorrentAttributeStatusUnknown

	torrentMissing := common.MatchedEntryPart{
		MediaPart: common.MediaPart{Size: 382842},
		TorrentInformation: common.TorrentInformation{
			Status:      common.TorrentStatusMissing,
			AgeStatus:   common.TorrentAttributeStatusUnknown,
			RatioStatus: common.TorrentAttributeStatusUnknown,
		},
	}
	tests := []struct {
		name                   string
		args                   args
		wantTorrentInformation common.TorrentInformation
	}{
		{
			name: "handle missing torrent entries",
			args: args{
				parts: []common.MatchedEntryPart{
					torrentMissing,
					torrentMissing,
				},
			},
			wantTorrentInformation: common.TorrentInformation{
				Status:      common.TorrentStatusMissing,
				RatioStatus: common.TorrentAttributeStatusUnknown,
				Ratio:       -1,
				AgeStatus:   common.TorrentAttributeStatusUnknown,
				Age:         -1,
			},
		},
		{
			name: "handle complete and missing torrent entries",
			args: args{
				parts: []common.MatchedEntryPart{
					torrentComplete,
					torrentMissing,
				},
			},
			wantTorrentInformation: common.TorrentInformation{
				Status:      common.TorrentStatusIncomplete,
				RatioStatus: common.TorrentAttributeStatusUnknown,
				Ratio:       -1,
				AgeStatus:   common.TorrentAttributeStatusUnknown,
				Age:         -1,
			},
		},
		{
			name: "handle single complete torrent entry",
			args: args{
				parts: []common.MatchedEntryPart{
					torrentComplete,
				},
			},
			wantTorrentInformation: common.TorrentInformation{
				Status: common.TorrentStatusPresent,
				Tracker: common.Tracker{
					Name:     "test",
					MinRatio: 0.5,
					MinAge:   time.Hour,
				},
				RatioStatus: common.TorrentAttributeStatusFulfilled,
				Ratio:       1.2,
				AgeStatus:   common.TorrentAttributeStatusFulfilled,
				Age:         time.Hour * 2,
			},
		},
		{
			name: "handle multiple complete torrent entries",
			args: args{
				parts: []common.MatchedEntryPart{
					torrentComplete,
					torrentComplete,
				},
			},
			wantTorrentInformation: common.TorrentInformation{
				Status:      common.TorrentStatusPresent,
				RatioStatus: common.TorrentAttributeStatusFulfilled,
				Ratio:       -1,
				AgeStatus:   common.TorrentAttributeStatusFulfilled,
				Age:         -1,
			},
		},
		{
			name: "handle multiple complete and incomplete ratio torrent entries",
			args: args{
				parts: []common.MatchedEntryPart{
					torrentComplete,
					torrentIncompleteRatio,
				},
			},
			wantTorrentInformation: common.TorrentInformation{
				Status:      common.TorrentStatusPresent,
				RatioStatus: common.TorrentAttributeStatusPending,
				Ratio:       -1,
				AgeStatus:   common.TorrentAttributeStatusFulfilled,
				Age:         -1,
			},
		},
		{
			name: "handle multiple complete and incomplete ratio torrent entries",
			args: args{
				parts: []common.MatchedEntryPart{
					torrentComplete,
					torrentIncompleteAge,
				},
			},
			wantTorrentInformation: common.TorrentInformation{
				Status:      common.TorrentStatusPresent,
				RatioStatus: common.TorrentAttributeStatusFulfilled,
				Ratio:       -1,
				AgeStatus:   common.TorrentAttributeStatusPending,
				Age:         -1,
			},
		},
		{
			name: "handle complete torrent entry with unknown tracker",
			args: args{
				parts: []common.MatchedEntryPart{
					torrentUnknownTracker,
				},
			},
			wantTorrentInformation: common.TorrentInformation{
				Status:      common.TorrentStatusPresent,
				RatioStatus: common.TorrentAttributeStatusUnknown,
				Ratio:       1.2,
				AgeStatus:   common.TorrentAttributeStatusUnknown,
				Age:         time.Hour * 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTorrentInformation := getBundledTorrentInformationFromParts(tt.args.parts)
			if !reflect.DeepEqual(gotTorrentInformation, tt.wantTorrentInformation) {
				t.Errorf("getBundledTorrentInformationFromParts() gotTorrentInformation = %v, want %v", gotTorrentInformation, tt.wantTorrentInformation)
			}
		})
	}
}
