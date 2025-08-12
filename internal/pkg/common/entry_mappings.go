package common

import (
	"fmt"
	"time"
)

var (
	ErrEntryMappingNotFound = fmt.Errorf("entry mapping not found")
)

type EntryMapping struct {
	// Id is the unique identifier of the entry.
	Id string
	// Name is the normalized name of the entry.
	Name EntryName
	// DateAdded is the date when the entry was added.
	DateAdded time.Time
	// Size is the size of the entry in bytes.
	Size int64
	// RetrieversFound holds a list of retrievers where this entry could be found.
	RetrieversFound []RetrieverInfo
}

// EntryMappingDetails holds the details of an entry mapping, which includes the retriever responses.
type EntryMappingDetails map[RetrieverId]string

func (e EntryMapping) String() string {
	return fmt.Sprintf("EntryMapping{Id: %q, Name: %q, RetrieversFound: %v}", e.Id, e.Name, e.RetrieversFound)
}

// EntryMappingManager is used to aggregate the results of single EntryRetriever instances and return the combined results.
type EntryMappingManager interface {
	// RefreshEntryMappings refreshes the entry mappings by querying every registered EntryRetriever and aggregating the results.
	RefreshEntryMappings() error

	// GetEntryMappings returns the filtered entry mapping by applying the given filters.
	GetEntryMappings(page int, pageSize int, filter EntryMappingFilter, sortBy EntryMappingSortBy, name string) ([]*EntryMapping, int, error)

	// GetEntryMappingById returns the entry mapping by its unique identifier.
	GetEntryMappingById(id string) (*EntryMapping, error)

	// GetRetrievers returns the information on all registered retrievers.
	GetRetrievers() ([]RetrieverInfo, error)

	// GetRetrieverById returns the retriever by its unique identifier.
	GetRetrieverById(id RetrieverId) (RetrieverInfo, EntryRetriever, error)

	// DeleteEntryMappingById deletes the entry mapping by its unique identifier.
	DeleteEntryMappingById(id string) error

	// GetEntryMappingDetails returns the details (retriever responses) of the entry mapping by its unique identifier.
	GetEntryMappingDetails(id string) (EntryMappingDetails, error)
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
	EntryMappingSortByNameAsc
	EntryMappingSortByNameDesc
)
