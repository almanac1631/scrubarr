package webserver

import (
	"net/url"
	"reflect"
	"testing"

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
