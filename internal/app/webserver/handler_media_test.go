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

func Test_getTorrentInformationFromParts(t *testing.T) {
	now = func() time.Time {
		return time.Date(2026, time.January, 19, 22, 0, 0, 0, time.UTC)
	}
	type args struct {
		parts []common.MatchedMediaPart
	}
	torrentComplete := common.MatchedMediaPart{
		MediaPart: common.MediaPart{Size: 382842},
		Tracker: common.Tracker{
			Name:     "test",
			MinRatio: 0.5,
			MinAge:   time.Hour,
		},
		TorrentFinding: &common.TorrentEntry{
			Added: now().Add(-2 * time.Hour),
			Ratio: 1.2,
		},
	}

	torrentIncompleteRatio := torrentComplete
	torrentFindingIncompleteRatio := *torrentComplete.TorrentFinding
	torrentFindingIncompleteRatio.Ratio = 0.1
	torrentIncompleteRatio.TorrentFinding = &torrentFindingIncompleteRatio

	torrentIncompleteAge := torrentComplete
	torrentFindingIncompleteAge := *torrentComplete.TorrentFinding
	torrentFindingIncompleteAge.Added = now().Add(2 * time.Hour)
	torrentIncompleteAge.TorrentFinding = &torrentFindingIncompleteAge

	torrentUnknownTracker := torrentComplete
	torrentUnknownTracker.Tracker = common.Tracker{}

	torrentMissing := common.MatchedMediaPart{
		MediaPart:      common.MediaPart{Size: 382842},
		TorrentFinding: nil,
	}
	tests := []struct {
		name                   string
		args                   args
		wantTotalSize          int64
		wantTorrentInformation TorrentInformation
	}{
		{
			name: "handle missing torrent entries",
			args: args{
				parts: []common.MatchedMediaPart{
					torrentMissing,
					torrentMissing,
				},
			},
			wantTotalSize: 765684,
			wantTorrentInformation: TorrentInformation{
				Status:      TorrentStatusMissing,
				RatioStatus: TorrentAttributeStatusUnknown,
				Ratio:       -1,
				MinRatio:    -1,
				AgeStatus:   TorrentAttributeStatusUnknown,
				Age:         -1,
				MinAge:      -1,
			},
		},
		{
			name: "handle complete and missing torrent entries",
			args: args{
				parts: []common.MatchedMediaPart{
					torrentComplete,
					torrentMissing,
				},
			},
			wantTotalSize: 765684,
			wantTorrentInformation: TorrentInformation{
				Status:      TorrentStatusIncomplete,
				RatioStatus: TorrentAttributeStatusUnknown,
				Ratio:       -1,
				MinRatio:    -1,
				AgeStatus:   TorrentAttributeStatusUnknown,
				Age:         -1,
				MinAge:      -1,
			},
		},
		{
			name: "handle single complete torrent entry",
			args: args{
				parts: []common.MatchedMediaPart{
					torrentComplete,
				},
			},
			wantTotalSize: 382842,
			wantTorrentInformation: TorrentInformation{
				Status:      TorrentStatusPresent,
				RatioStatus: TorrentAttributeStatusFulfilled,
				Ratio:       1.2,
				MinRatio:    0.5,
				AgeStatus:   TorrentAttributeStatusFulfilled,
				Age:         time.Hour * 2,
				MinAge:      time.Hour,
			},
		},
		{
			name: "handle multiple complete torrent entries",
			args: args{
				parts: []common.MatchedMediaPart{
					torrentComplete,
					torrentComplete,
				},
			},
			wantTotalSize: 765684,
			wantTorrentInformation: TorrentInformation{
				Status:      TorrentStatusPresent,
				RatioStatus: TorrentAttributeStatusFulfilled,
				Ratio:       -1,
				MinRatio:    -1,
				AgeStatus:   TorrentAttributeStatusFulfilled,
				Age:         -1,
				MinAge:      -1,
			},
		},
		{
			name: "handle multiple complete and incomplete ratio torrent entries",
			args: args{
				parts: []common.MatchedMediaPart{
					torrentComplete,
					torrentIncompleteRatio,
				},
			},
			wantTotalSize: 765684,
			wantTorrentInformation: TorrentInformation{
				Status:      TorrentStatusPresent,
				RatioStatus: TorrentAttributeStatusPending,
				Ratio:       -1,
				MinRatio:    -1,
				AgeStatus:   TorrentAttributeStatusFulfilled,
				Age:         -1,
				MinAge:      -1,
			},
		},
		{
			name: "handle multiple complete and incomplete ratio torrent entries",
			args: args{
				parts: []common.MatchedMediaPart{
					torrentComplete,
					torrentIncompleteAge,
				},
			},
			wantTotalSize: 765684,
			wantTorrentInformation: TorrentInformation{
				Status:      TorrentStatusPresent,
				RatioStatus: TorrentAttributeStatusFulfilled,
				Ratio:       -1,
				MinRatio:    -1,
				AgeStatus:   TorrentAttributeStatusPending,
				Age:         -1,
				MinAge:      -1,
			},
		},
		{
			name: "handle complete torrent entry with unknown tracker",
			args: args{
				parts: []common.MatchedMediaPart{
					torrentUnknownTracker,
				},
			},
			wantTotalSize: 382842,
			wantTorrentInformation: TorrentInformation{
				Status:      TorrentStatusPresent,
				RatioStatus: TorrentAttributeStatusUnknown,
				Ratio:       1.2,
				MinRatio:    -1,
				AgeStatus:   TorrentAttributeStatusUnknown,
				Age:         time.Hour * 2,
				MinAge:      -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTotalSize, gotTorrentInformation := getTorrentInformationFromParts(tt.args.parts)
			if gotTotalSize != tt.wantTotalSize {
				t.Errorf("getTorrentInformationFromParts() gotTotalSize = %v, want %v", gotTotalSize, tt.wantTotalSize)
			}
			if !reflect.DeepEqual(gotTorrentInformation, tt.wantTorrentInformation) {
				t.Errorf("getTorrentInformationFromParts() gotTorrentInformation = %v, want %v", gotTorrentInformation, tt.wantTorrentInformation)
			}
		})
	}
}
