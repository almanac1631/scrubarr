package inmemory

import (
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"reflect"
	"testing"
)

var (
	firstEntry = common.EntryPresencePairs{
		&common.RetrieverInfo{Category: "c1", Name: "r1"}: &common.Entry{},
		&common.RetrieverInfo{Category: "c1", Name: "r2"}: nil,
		&common.RetrieverInfo{Category: "c2", Name: "r3"}: &common.Entry{},
	}
	secondEntry = common.EntryPresencePairs{
		&common.RetrieverInfo{Category: "c1", Name: "r1"}: nil,
		&common.RetrieverInfo{Category: "c1", Name: "r2"}: nil,
		&common.RetrieverInfo{Category: "c2", Name: "r3"}: &common.Entry{},
	}
	thirdEntry = common.EntryPresencePairs{
		&common.RetrieverInfo{Category: "c1", Name: "r1"}: &common.Entry{},
		&common.RetrieverInfo{Category: "c1", Name: "r2"}: nil,
		&common.RetrieverInfo{Category: "c2", Name: "r3"}: nil,
	}
)

func Test_applyFilter(t *testing.T) {
	entryMappings := map[common.EntryName]common.EntryPresencePairs{
		"first":  firstEntry,
		"second": secondEntry,
		"third":  thirdEntry,
	}
	type args struct {
		entryMappings map[common.EntryName]common.EntryPresencePairs
		filter        common.EntryMappingFilter
	}
	tests := []struct {
		name string
		args args
		want map[common.EntryName]common.EntryPresencePairs
	}{
		{
			"test no filter", args{entryMappings, common.EntryMappingFilterNoFilter},
			entryMappings,
		},
		{
			"test filter complete", args{entryMappings, common.EntryMappingFilterCompleteEntry},
			map[common.EntryName]common.EntryPresencePairs{"first": firstEntry},
		},
		{
			"test filter incomplete", args{entryMappings, common.EntryMappingFilterIncompleteEntry},
			map[common.EntryName]common.EntryPresencePairs{"second": secondEntry, "third": thirdEntry},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyFilter(tt.args.entryMappings, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_areEntryPresencePairsComplete(t *testing.T) {
	type args struct {
		pairs common.EntryPresencePairs
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test complete presence pairs with only one retriever",
			args{common.EntryPresencePairs{
				&common.RetrieverInfo{Category: "c1", Name: "r1"}: &common.Entry{},
			}},
			true,
		},
		{
			"test incomplete presence pairs",
			args{common.EntryPresencePairs{
				&common.RetrieverInfo{Category: "c1", Name: "r1"}: &common.Entry{},
				&common.RetrieverInfo{Category: "c2", Name: "r2"}: nil,
			}},
			false,
		},
		{
			"test complete presence pairs based on categories",
			args{common.EntryPresencePairs{
				&common.RetrieverInfo{Category: "c1", Name: "r1"}: &common.Entry{},
				&common.RetrieverInfo{Category: "c1", Name: "r1"}: nil,
				&common.RetrieverInfo{Category: "c2", Name: "r2"}: &common.Entry{},
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := areEntryPresencePairsComplete(tt.args.pairs); got != tt.want {
				t.Errorf("areEntryPresencePairsComplete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPageExcerpt(t *testing.T) {
	type args struct {
		entryMappings map[common.EntryName]common.EntryPresencePairs
		page          int
		pageSize      int
	}
	tests := []struct {
		name string
		args args
		want map[common.EntryName]common.EntryPresencePairs
	}{
		{
			"test page excerpt full response",
			args{map[common.EntryName]common.EntryPresencePairs{
				"first": firstEntry,
			}, 1, 10},
			map[common.EntryName]common.EntryPresencePairs{
				"first": firstEntry,
			},
		},
		{
			"test valid partial page excerpt and respect order",
			args{map[common.EntryName]common.EntryPresencePairs{
				"a": firstEntry,
				"b": firstEntry,
				"c": firstEntry,
				"d": firstEntry,
				"e": firstEntry,
				"f": firstEntry,
				"g": firstEntry,
				"h": firstEntry,
			}, 2, 2},
			map[common.EntryName]common.EntryPresencePairs{
				"c": firstEntry,
				"d": firstEntry,
			},
		},
		{
			"test partial page excerpt end out of bounds",
			args{map[common.EntryName]common.EntryPresencePairs{
				"a": firstEntry,
				"b": firstEntry,
				"c": firstEntry,
			}, 2, 2},
			map[common.EntryName]common.EntryPresencePairs{
				"c": firstEntry,
			},
		},
		{
			"test partial page excerpt offset out of bounds",
			args{map[common.EntryName]common.EntryPresencePairs{
				"a": firstEntry,
				"b": firstEntry,
				"c": firstEntry,
			}, 4, 1},
			map[common.EntryName]common.EntryPresencePairs{},
		},
		{
			"test page excerpt empty entry mapping",
			args{map[common.EntryName]common.EntryPresencePairs{}, 2, 2},
			map[common.EntryName]common.EntryPresencePairs{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPageExcerpt(tt.args.entryMappings, tt.args.page, tt.args.pageSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPageExcerpt() = %v, want %v", got, tt.want)
			}
		})
	}
}
