package common

import (
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval"
	"strings"
)

type EntryPresencePairs map[RetrieverInfo]*retrieval.Entry

func (mapping EntryPresencePairs) String() string {
	stringBuilder := &strings.Builder{}
	stringBuilder.WriteString("[")
	for retrieverInfo, entry := range mapping {
		stateString := "/"
		if entry != nil {
			stateString = "+"
		}
		stringBuilder.WriteString(fmt.Sprintf("%s=%s", retrieverInfo.Id(), stateString))
	}
	stringBuilder.WriteString("]")
	return stringBuilder.String()
}

func (mapping EntryPresencePairs) IsComplete() bool {
	for _, entry := range mapping {
		if entry == nil {
			return false
		}
	}
	return true
}

func (mapping EntryPresencePairs) Name() retrieval.EntryName {
	var entryPresent *retrieval.Entry
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
