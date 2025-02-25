package simple

import (
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"reflect"
	"testing"
)

type testEntryRetriever struct {
	entries common.RetrieverEntries
	err     error
}

func (t testEntryRetriever) RetrieveEntries() (common.RetrieverEntries, error) {
	return t.entries, t.err
}

func TestBundledEntryRetriever(t *testing.T) {
	retrieverInfo1 := common.RetrieverInfo{
		Category:     "somecategory",
		SoftwareName: "somesoftware",
		Name:         "retr1",
	}
	retrieverEntries1 := common.RetrieverEntries{
		common.EntryName("some name"): common.Entry{},
	}
	retriever1 := testEntryRetriever{
		entries: retrieverEntries1,
		err:     nil,
	}
	retrieverInfo2 := common.RetrieverInfo{
		Category:     "anothercategory",
		SoftwareName: "someothersoftware",
		Name:         "retr2",
	}
	retrieverEntries2 := common.RetrieverEntries{
		common.EntryName("another name"): common.Entry{},
	}
	retriever2 := testEntryRetriever{
		entries: retrieverEntries2,
		err:     nil,
	}

	type args struct {
		entryRetrievers map[common.RetrieverInfo]common.EntryRetriever
	}
	tests := []struct {
		name    string
		args    args
		want    map[common.RetrieverInfo]common.RetrieverEntries
		wantErr bool
	}{
		{
			"does not return entries when no retriever is registered",
			args{map[common.RetrieverInfo]common.EntryRetriever{}},
			map[common.RetrieverInfo]common.RetrieverEntries{},
			false,
		},
		{
			"can retrieve entries from single entry retriever",
			args{map[common.RetrieverInfo]common.EntryRetriever{retrieverInfo1: retriever1}},
			map[common.RetrieverInfo]common.RetrieverEntries{
				retrieverInfo1: retrieverEntries1,
			},
			false,
		},
		{
			"can retrieve entries from multiple entry retrievers",
			args{map[common.RetrieverInfo]common.EntryRetriever{retrieverInfo1: retriever1, retrieverInfo2: retriever2}},
			map[common.RetrieverInfo]common.RetrieverEntries{
				retrieverInfo1: retrieverEntries1,
				retrieverInfo2: retrieverEntries2,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BundledEntryRetriever(tt.args.entryRetrievers)
			if (err != nil) != tt.wantErr {
				t.Errorf("BundledEntryRetriever() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BundledEntryRetriever() got = %v, want %v", got, tt.want)
			}
		})
	}
}
