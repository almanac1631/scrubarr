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
	entryMappings := []common.EntryMapping{
		{"first", firstEntry},
		{"second", secondEntry},
		{"third", thirdEntry},
	}
	type args struct {
		entryMappings []common.EntryMapping
		filter        common.EntryMappingFilter
	}
	tests := []struct {
		name string
		args args
		want []common.EntryMapping
	}{
		{
			"test no filter", args{entryMappings, common.EntryMappingFilterNoFilter},
			entryMappings,
		},
		{
			"test filter complete", args{entryMappings, common.EntryMappingFilterCompleteEntry},
			[]common.EntryMapping{{"first", firstEntry}},
		},
		{
			"test filter incomplete", args{entryMappings, common.EntryMappingFilterIncompleteEntry},
			[]common.EntryMapping{{"second", secondEntry}, {"third", thirdEntry}},
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
		entryMappings []common.EntryMapping
		page          int
		pageSize      int
	}
	tests := []struct {
		name string
		args args
		want []common.EntryMapping
	}{
		{
			"test page excerpt full response",
			args{[]common.EntryMapping{
				{"first", firstEntry},
			}, 1, 10},
			[]common.EntryMapping{
				{"first", firstEntry},
			},
		},
		{
			"test valid partial page excerpt and respect order",
			args{[]common.EntryMapping{
				{"a", firstEntry},
				{"b", firstEntry},
				{"c", firstEntry},
				{"d", firstEntry},
				{"e", firstEntry},
				{"f", firstEntry},
				{"g", firstEntry},
				{"h", firstEntry},
			}, 2, 2},
			[]common.EntryMapping{
				{"c", firstEntry},
				{"d", firstEntry},
			},
		},
		{
			"test partial page excerpt end out of bounds",
			args{[]common.EntryMapping{
				{"a", firstEntry},
				{"b", firstEntry},
				{"c", firstEntry},
			}, 2, 2},
			[]common.EntryMapping{
				{"c", firstEntry},
			},
		},
		{
			"test partial page excerpt offset out of bounds",
			args{[]common.EntryMapping{
				{"a", firstEntry},
				{"b", firstEntry},
				{"c", firstEntry},
			}, 4, 1},
			[]common.EntryMapping{},
		},
		{
			"test page excerpt empty entry mapping",
			args{[]common.EntryMapping{}, 2, 2},
			[]common.EntryMapping{},
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
