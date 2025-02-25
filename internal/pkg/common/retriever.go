package common

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
)

type RetrieverId string

type RetrieverInfo struct {
	Category     string
	SoftwareName string
	Name         string
}

func (info RetrieverInfo) String() string {
	return fmt.Sprintf("%s#%s#%s", info.Category, info.SoftwareName, info.Name)
}

func (info RetrieverInfo) Id() RetrieverId {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(info.String()))
	hashSum := hash.Sum(make([]byte, 0))
	return RetrieverId(hex.EncodeToString(hashSum))
}

// RetrieverEntries is the return type of an EntryRetriever instance that successfully retrieved all available entries in the data source.
type RetrieverEntries map[EntryName]Entry

// EntryRetriever is the interface used to implement new entry retrieves (e.g. from folders, *arr apps or torrent clients)
type EntryRetriever interface {
	// RetrieveEntries retrieves a mapping of entries consisting of  unique entry names and the actual entry as value.
	RetrieveEntries() (RetrieverEntries, error)
}

// BundledEntryRetriever queries all of the EntryRetriever instances, returns the responses and blocks until completion.
type BundledEntryRetriever func(map[RetrieverInfo]EntryRetriever) (map[RetrieverInfo]RetrieverEntries, error)
