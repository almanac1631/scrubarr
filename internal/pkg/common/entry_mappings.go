package common

import (
	"fmt"
	"strings"
)

type EntryMapping struct {
	Name  EntryName
	Pairs EntryPresencePairs
}

// EntryMappingManager is used to aggregate the results of single EntryRetriever instances and return the combined results.
type EntryMappingManager interface {
	// RefreshEntryMappings refreshes the entry mappings by querying every registered EntryRetriever and aggregating the results.
	RefreshEntryMappings() error

	// GetEntryMappings returns the filtered entry mapping by applying the given filters.
	GetEntryMappings(page int, pageSize int, filter EntryMappingFilter) ([]EntryMapping, int, error)

	// GetRetrievers returns the information on all registered retrievers.
	GetRetrievers() ([]RetrieverInfo, error)
}

type EntryMappingFilter int

const (
	EntryMappingFilterNoFilter EntryMappingFilter = iota
	EntryMappingFilterIncompleteEntry
	EntryMappingFilterCompleteEntry
)

// EntryPresencePairs is a struct that holds the findings within the retrievers of a given entry.
type EntryPresencePairs struct {
	// Name is the normalized name of the entry.
	Name EntryName
	// RetrieversFound holds a list of retrievers where this entry could be found.
	RetrieversFound []RetrieverInfo
}

func (mapping EntryPresencePairs) String() string {
	stringBuilder := &strings.Builder{}
	stringBuilder.WriteString(fmt.Sprintf("EntryPresencePairs{Name: %q, RetrieversFound: [", mapping.Name))
	for i, retriever := range mapping.RetrieversFound {
		if i > 0 {
			stringBuilder.WriteString(", ")
		}
		stringBuilder.WriteString(retriever.String())
	}
	stringBuilder.WriteString("]}")
	return stringBuilder.String()
}
