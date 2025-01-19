package retrieval

import (
	"github.com/almanac1631/scrubarr/internal/pkg/config"
)

type Entry struct {
	Name           EntryName
	AdditionalData any
}

type EntryName string

type EntryRetriever interface {
	RetrieveEntries() (map[EntryName]Entry, error)
}

type EntryRetrieverInitializer func(config.EntryAccessor) (EntryRetriever, error)
