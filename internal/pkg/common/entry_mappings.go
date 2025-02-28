package common

import (
	"fmt"
	"time"
)

type EntryMapping struct {
	// Name is the normalized name of the entry.
	Name EntryName
	// DateAdded is the date when the entry was added.
	DateAdded time.Time
	// Size is the size of the entry in bytes.
	Size int64
	// RetrieversFound holds a list of retrievers where this entry could be found.
	RetrieversFound []RetrieverInfo
}

func (e EntryMapping) String() string {
	return fmt.Sprintf("EntryMapping{Name: %q, RetrieversFound: %v}", e.Name, e.RetrieversFound)
}

// EntryMappingManager is used to aggregate the results of single EntryRetriever instances and return the combined results.
type EntryMappingManager interface {
	// RefreshEntryMappings refreshes the entry mappings by querying every registered EntryRetriever and aggregating the results.
	RefreshEntryMappings() error

	// GetEntryMappings returns the filtered entry mapping by applying the given filters.
	GetEntryMappings(page int, pageSize int, filter EntryMappingFilter, sortBy EntryMappingSortBy) ([]*EntryMapping, int, error)

	// GetRetrievers returns the information on all registered retrievers.
	GetRetrievers() ([]RetrieverInfo, error)
}

type EntryMappingFilter int

const (
	EntryMappingFilterNoFilter EntryMappingFilter = iota
	EntryMappingFilterIncompleteEntry
	EntryMappingFilterCompleteEntry
)

type EntryMappingSortBy int

const (
	EntryMappingSortByNoSort EntryMappingSortBy = iota
	EntryMappingSortByDateAsc
	EntryMappingSortByDateDesc
	EntryMappingSortBySizeAsc
	EntryMappingSortBySizeDesc
)
