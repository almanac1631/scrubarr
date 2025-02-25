package inmemory

import (
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"maps"
	"slices"
	"sort"
	"strings"
	"sync"
)

type EntryMappingManager struct {
	entryRetrievers       map[common.RetrieverInfo]common.EntryRetriever
	bundledEntryRetriever common.BundledEntryRetriever

	entryMappingsLock *sync.Mutex
	rawEntryMappings  map[common.EntryName]common.EntryPresencePairs
}

func NewEntryMappingManager(entryRetrievers map[common.RetrieverInfo]common.EntryRetriever, bundledEntryRetriever common.BundledEntryRetriever) common.EntryMappingManager {
	return &EntryMappingManager{
		entryRetrievers:       entryRetrievers,
		bundledEntryRetriever: bundledEntryRetriever,

		entryMappingsLock: &sync.Mutex{},
		rawEntryMappings:  make(map[common.EntryName]common.EntryPresencePairs),
	}
}

func (e *EntryMappingManager) RefreshEntryMappings() error {
	e.entryMappingsLock.Lock()
	defer e.entryMappingsLock.Unlock()
	e.rawEntryMappings = make(map[common.EntryName]common.EntryPresencePairs)
	rawEntries, err := e.bundledEntryRetriever(e.entryRetrievers)
	if err != nil {
		return fmt.Errorf("could not query entries using given entry retriever: %w", err)
	}
	setOfNames := getSetOfNames(rawEntries)
	for name := range setOfNames {
		presencePair := common.EntryPresencePairs{}
		for retrieverInfo, retrieverResult := range rawEntries {
			entry, ok := retrieverResult[name]
			if !ok {
				continue
			}
			presencePair[&retrieverInfo] = &entry
		}
		e.rawEntryMappings[name] = presencePair
	}
	return nil
}

func getSetOfNames(retrieverEntries map[common.RetrieverInfo]common.RetrieverEntries) map[common.EntryName]any {
	setOfNames := map[common.EntryName]any{}
	for _, entries := range retrieverEntries {
		for name, _ := range entries {
			setOfNames[name] = true
		}
	}
	return setOfNames
}

func (e *EntryMappingManager) GetEntryMappings(page int, pageSize int, filter common.EntryMappingFilter) ([]common.EntryMapping, int, error) {
	entryMappings := make([]common.EntryMapping, len(e.rawEntryMappings))
	i := 0
	for entryName, pairs := range e.rawEntryMappings {
		entryMappings[i] = common.EntryMapping{Name: entryName, Pairs: pairs}
		i++
	}
	filteredEntryMappings := applyFilter(entryMappings, filter)
	totalAmount := len(filteredEntryMappings)
	return getPageExcerpt(filteredEntryMappings, page, pageSize), totalAmount, nil
}

func getPageExcerpt(entryMappings []common.EntryMapping, page int, pageSize int) []common.EntryMapping {
	offset := (page - 1) * pageSize
	if offset > len(entryMappings) {
		return []common.EntryMapping{}
	}
	sort.SliceStable(entryMappings, func(i, j int) bool {
		return strings.Compare(string(entryMappings[i].Name), string(entryMappings[j].Name)) == -1
	})
	end := offset + pageSize
	if end > len(entryMappings) {
		end = len(entryMappings)
	}
	return entryMappings[offset:end]
}

func applyFilter(entryMappings []common.EntryMapping, filter common.EntryMappingFilter) []common.EntryMapping {
	if filter == common.EntryMappingFilterNoFilter {
		return entryMappings
	}
	var filteredEntryMappings []common.EntryMapping
	for _, entryMapping := range entryMappings {
		entryPresencePairsComplete := areEntryPresencePairsComplete(entryMapping.Pairs)
		if (entryPresencePairsComplete && filter == common.EntryMappingFilterCompleteEntry) ||
			(!entryPresencePairsComplete && filter == common.EntryMappingFilterIncompleteEntry) {
			filteredEntryMappings = append(filteredEntryMappings, entryMapping)
		}
	}
	return filteredEntryMappings
}

func areEntryPresencePairsComplete(pairs common.EntryPresencePairs) bool {
	categories := make(map[string]bool, len(pairs))
	for retrieverInfo := range pairs {
		categories[retrieverInfo.Category] = false
	}
	expectedCategoryCount := len(categories)
	categories = make(map[string]bool, expectedCategoryCount)

	for retrieverInfo, pair := range pairs {
		if pair == nil {
			continue
		}
		categories[retrieverInfo.Category] = true
	}
	return len(categories) == expectedCategoryCount
}

func (e *EntryMappingManager) GetRetrievers() ([]common.RetrieverInfo, error) {
	return slices.Collect(maps.Keys(e.entryRetrievers)), nil
}
