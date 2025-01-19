package common

import (
	"encoding/hex"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval"
	"hash/fnv"
	"log/slog"
	"sync"
)

type RetrieverInfo struct {
	Category     string
	SoftwareName string
	Name         string
}

type RetrieverId string

func (info RetrieverInfo) String() string {
	return fmt.Sprintf("%s#%s#%s", info.Category, info.SoftwareName, info.Name)
}

func (info RetrieverInfo) Id() RetrieverId {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(info.String()))
	hashSum := hash.Sum(make([]byte, 0))
	return RetrieverId(hex.EncodeToString(hashSum))
}

type RetrieverRegistry interface {
	RetrieveEntryMapping() map[retrieval.EntryName]EntryPresencePairs
	GetRetrievers() []RetrieverInfo
}

type MapBasedRetrieverRegistry map[RetrieverInfo]retrieval.EntryRetriever

func (retrieverRegistry MapBasedRetrieverRegistry) RetrieveEntryMapping() map[retrieval.EntryName]EntryPresencePairs {
	retrieverEntriesAll := retrieverRegistry.RetrieveEntriesForRetrievers()
	setOfNames := getSetOfNames(retrieverEntriesAll)
	entryMappings := map[retrieval.EntryName]EntryPresencePairs{}
	for name := range setOfNames {
		presencePair := EntryPresencePairs{}
		for retrieverInfo, retrieverResult := range retrieverEntriesAll {
			entry, ok := retrieverResult[name]
			if !ok {
				continue
			}
			presencePair[retrieverInfo] = &entry
		}
		entryMappings[name] = presencePair
	}
	return entryMappings
}

func (retrieverRegistry MapBasedRetrieverRegistry) RetrieveEntriesForRetrievers() map[RetrieverInfo]map[retrieval.EntryName]retrieval.Entry {
	retrieverEntries := map[RetrieverInfo]map[retrieval.EntryName]retrieval.Entry{}
	type channelEntry struct {
		retrieverInfo RetrieverInfo
		values        map[retrieval.EntryName]retrieval.Entry
	}
	resultsChan := make(chan channelEntry, len(retrieverRegistry))
	syncGroup := &sync.WaitGroup{}
	for retrieverInfo, retriever := range retrieverRegistry {
		syncGroup.Add(1)
		go func() {
			defer syncGroup.Done()
			rSlog := slog.With("retriever", retrieverInfo)
			rSlog.Debug("retrieving entries...")
			entries, err := retriever.RetrieveEntries()
			if err != nil {
				rSlog.Error("could not retrieve entries", "err", err)
			}
			resultsChan <- channelEntry{retrieverInfo, entries}
			rSlog.Debug("done with retrieval")
		}()
	}
	go func() {
		syncGroup.Wait()
		close(resultsChan)
	}()
	for result := range resultsChan {
		retrieverEntries[result.retrieverInfo] = result.values
	}
	return retrieverEntries
}

func (retrieverRegistry MapBasedRetrieverRegistry) GetRetrievers() []RetrieverInfo {
	retrievers := make([]RetrieverInfo, 0)
	for retrieverInfo := range retrieverRegistry {
		retrievers = append(retrievers, retrieverInfo)
	}
	return retrievers
}

func getSetOfNames(retrieverEntries map[RetrieverInfo]map[retrieval.EntryName]retrieval.Entry) map[retrieval.EntryName]bool {
	setOfNames := map[retrieval.EntryName]bool{}
	for _, entries := range retrieverEntries {
		for name, _ := range entries {
			setOfNames[name] = true
		}
	}
	return setOfNames
}
