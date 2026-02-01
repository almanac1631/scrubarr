package webserver

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/almanac1631/scrubarr/pkg/domain"
)

func Test_getSortInfoFromUrlQuery(t *testing.T) {
	type args struct {
		values url.Values
	}
	tests := []struct {
		name string
		args args
		want domain.SortInfo
	}{
		{
			name: "test valid parse from url query",
			args: args{
				values: map[string][]string{
					"sortKey":   {"size"},
					"sortOrder": {"asc"},
				},
			},
			want: domain.SortInfo{
				Key:   domain.SortKeySize,
				Order: domain.SortOrderAsc,
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
			want: domain.SortInfo{
				Key:   domain.SortKeyName,
				Order: domain.SortOrderAsc,
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
		parts []domain.MatchedMediaPart
	}
	torrentComplete := domain.MatchedMediaPart{
		MediaPart: domain.MediaPart{Size: 382842},
		TorrentInformation: domain.TorrentInformation{
			Client: "sonarr",
			Id:     "abc",
			Status: domain.TorrentStatusPresent,
			Tracker: domain.Tracker{
				Name:     "test",
				MinRatio: 0.5,
				MinAge:   time.Hour,
			},
			Age:         time.Hour * 2,
			AgeStatus:   domain.TorrentAttributeStatusFulfilled,
			Ratio:       1.2,
			RatioStatus: domain.TorrentAttributeStatusFulfilled,
		},
	}

	torrentIncompleteRatio := torrentComplete
	torrentIncompleteRatio.TorrentInformation.Ratio = 0.1
	torrentIncompleteRatio.TorrentInformation.RatioStatus = domain.TorrentAttributeStatusPending

	torrentIncompleteAge := torrentComplete
	torrentIncompleteAge.TorrentInformation.Age = time.Minute * 30
	torrentIncompleteAge.TorrentInformation.AgeStatus = domain.TorrentAttributeStatusPending

	torrentUnknownTracker := torrentComplete
	torrentUnknownTracker.TorrentInformation.Tracker = domain.Tracker{}
	torrentUnknownTracker.TorrentInformation.AgeStatus = domain.TorrentAttributeStatusUnknown
	torrentUnknownTracker.TorrentInformation.RatioStatus = domain.TorrentAttributeStatusUnknown

	torrentMissing := domain.MatchedMediaPart{
		MediaPart: domain.MediaPart{Size: 382842},
		TorrentInformation: domain.TorrentInformation{
			Status:      domain.TorrentStatusMissing,
			AgeStatus:   domain.TorrentAttributeStatusUnknown,
			RatioStatus: domain.TorrentAttributeStatusUnknown,
		},
	}
	tests := []struct {
		name                   string
		args                   args
		wantTorrentInformation domain.TorrentInformation
	}{
		{
			name: "handle missing torrent entries",
			args: args{
				parts: []domain.MatchedMediaPart{
					torrentMissing,
					torrentMissing,
				},
			},
			wantTorrentInformation: domain.TorrentInformation{
				Status:      domain.TorrentStatusMissing,
				RatioStatus: domain.TorrentAttributeStatusUnknown,
				Ratio:       -1,
				AgeStatus:   domain.TorrentAttributeStatusUnknown,
				Age:         -1,
			},
		},
		{
			name: "handle complete and missing torrent entries",
			args: args{
				parts: []domain.MatchedMediaPart{
					torrentComplete,
					torrentMissing,
				},
			},
			wantTorrentInformation: domain.TorrentInformation{
				Status:      domain.TorrentStatusIncomplete,
				RatioStatus: domain.TorrentAttributeStatusUnknown,
				Ratio:       -1,
				AgeStatus:   domain.TorrentAttributeStatusUnknown,
				Age:         -1,
			},
		},
		{
			name: "handle single complete torrent entry",
			args: args{
				parts: []domain.MatchedMediaPart{
					torrentComplete,
				},
			},
			wantTorrentInformation: domain.TorrentInformation{
				Status: domain.TorrentStatusPresent,
				Tracker: domain.Tracker{
					Name:     "test",
					MinRatio: 0.5,
					MinAge:   time.Hour,
				},
				RatioStatus: domain.TorrentAttributeStatusFulfilled,
				Ratio:       1.2,
				AgeStatus:   domain.TorrentAttributeStatusFulfilled,
				Age:         time.Hour * 2,
			},
		},
		{
			name: "handle multiple complete torrent entries",
			args: args{
				parts: []domain.MatchedMediaPart{
					torrentComplete,
					torrentComplete,
				},
			},
			wantTorrentInformation: domain.TorrentInformation{
				Status:      domain.TorrentStatusPresent,
				RatioStatus: domain.TorrentAttributeStatusFulfilled,
				Ratio:       -1,
				AgeStatus:   domain.TorrentAttributeStatusFulfilled,
				Age:         -1,
			},
		},
		{
			name: "handle multiple complete and incomplete ratio torrent entries",
			args: args{
				parts: []domain.MatchedMediaPart{
					torrentComplete,
					torrentIncompleteRatio,
				},
			},
			wantTorrentInformation: domain.TorrentInformation{
				Status:      domain.TorrentStatusPresent,
				RatioStatus: domain.TorrentAttributeStatusPending,
				Ratio:       -1,
				AgeStatus:   domain.TorrentAttributeStatusFulfilled,
				Age:         -1,
			},
		},
		{
			name: "handle multiple complete and incomplete ratio torrent entries",
			args: args{
				parts: []domain.MatchedMediaPart{
					torrentComplete,
					torrentIncompleteAge,
					torrentComplete,
				},
			},
			wantTorrentInformation: domain.TorrentInformation{
				Status:      domain.TorrentStatusPresent,
				RatioStatus: domain.TorrentAttributeStatusFulfilled,
				Ratio:       -1,
				AgeStatus:   domain.TorrentAttributeStatusPending,
				Age:         -1,
			},
		},
		{
			name: "handle complete torrent entry with unknown tracker",
			args: args{
				parts: []domain.MatchedMediaPart{
					torrentUnknownTracker,
				},
			},
			wantTorrentInformation: domain.TorrentInformation{
				Status:      domain.TorrentStatusPresent,
				RatioStatus: domain.TorrentAttributeStatusUnknown,
				Ratio:       1.2,
				AgeStatus:   domain.TorrentAttributeStatusUnknown,
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
