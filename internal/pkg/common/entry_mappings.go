package common

import (
	"fmt"
	"strings"
)

// EntryMappingManager is used to aggregate the results of single EntryRetriever instances and return the combined results.
type EntryMappingManager interface {
	// RefreshEntryMappings refreshes the entry mappings by querying every registered EntryRetriever and aggregating the results.
	RefreshEntryMappings() error

	// GetEntryMappings returns the filtered entry mapping by applying the given filters.
	GetEntryMappings(page int, pageSize int, filter EntryMappingFilter) (map[EntryName]EntryPresencePairs, int, error)

	// GetRetrievers returns the information on all registered retrievers.
	GetRetrievers() ([]RetrieverInfo, error)
}

type EntryMappingFilter int

const (
	EntryMappingFilterNoFilter EntryMappingFilter = iota
	EntryMappingFilterIncompleteEntry
	EntryMappingFilterCompleteEntry
)

// EntryPresencePairs is a map that	holds the findings within the retrievers of a given entry.
type EntryPresencePairs map[*RetrieverInfo]*Entry

func (mapping EntryPresencePairs) String() string {
	stringBuilder := &strings.Builder{}
	stringBuilder.WriteString("[")
	for retrieverId, entry := range mapping {
		stateString := "/"
		if entry != nil {
			stateString = "+"
		}
		stringBuilder.WriteString(fmt.Sprintf("%s=%s", retrieverId, stateString))
	}
	stringBuilder.WriteString("]")
	return stringBuilder.String()
}

func (mapping EntryPresencePairs) Name() EntryName {
	var entryPresent *Entry
	for _, entry := range mapping {
		if entry == nil {
			continue
		}
		entryPresent = entry
		break
	}
	if entryPresent == nil {
		panic("entry presence mapping has to contain at least one non-nil entry")
	}
	return entryPresent.Name
}
