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
	entryMappings     map[common.EntryName]common.EntryPresencePairs
}

func NewEntryMappingManager(entryRetrievers map[common.RetrieverInfo]common.EntryRetriever, bundledEntryRetriever common.BundledEntryRetriever) *EntryMappingManager {
	return &EntryMappingManager{
		entryRetrievers:       entryRetrievers,
		bundledEntryRetriever: bundledEntryRetriever,

		entryMappingsLock: &sync.Mutex{},
		entryMappings:     make(map[common.EntryName]common.EntryPresencePairs),
	}
}

func (e *EntryMappingManager) RefreshEntryMappings() error {
	e.entryMappingsLock.Lock()
	defer e.entryMappingsLock.Unlock()
	e.entryMappings = make(map[common.EntryName]common.EntryPresencePairs)
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
		e.entryMappings[name] = presencePair
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

func (e *EntryMappingManager) GetEntryMappings(page int, pageSize int, filter common.EntryMappingFilter) (map[common.EntryName]common.EntryPresencePairs, int, error) {
	filteredEntryMappings := applyFilter(e.entryMappings, filter)
	totalAmount := len(filteredEntryMappings)
	return getPageExcerpt(filteredEntryMappings, page, pageSize), totalAmount, nil
}

func getPageExcerpt(entryMappings map[common.EntryName]common.EntryPresencePairs, page int, pageSize int) map[common.EntryName]common.EntryPresencePairs {
	offset := (page - 1) * pageSize
	if offset > len(entryMappings) {
		return map[common.EntryName]common.EntryPresencePairs{}
	}
	entryNames := slices.Collect(maps.Keys(entryMappings))
	sort.SliceStable(entryNames, func(i, j int) bool {
		return strings.Compare(string(entryNames[i]), string(entryNames[j])) == -1
	})
	end := offset + pageSize
	if end > len(entryNames) {
		end = len(entryNames)
	}
	selectedEntryNames := entryNames[offset:end]
	entryMappingsExcerpt := map[common.EntryName]common.EntryPresencePairs{}
	for _, entryName := range selectedEntryNames {
		entryMappingsExcerpt[entryName] = entryMappings[entryName]
	}
	return entryMappingsExcerpt
}

func applyFilter(entryMappings map[common.EntryName]common.EntryPresencePairs, filter common.EntryMappingFilter) map[common.EntryName]common.EntryPresencePairs {
	if filter == common.EntryMappingFilterNoFilter {
		return entryMappings
	}
	filteredEntryMappings := map[common.EntryName]common.EntryPresencePairs{}
	for entryName, entry := range entryMappings {
		entryPresencePairsComplete := areEntryPresencePairsComplete(entry)
		if entryPresencePairsComplete && filter == common.EntryMappingFilterCompleteEntry {
			filteredEntryMappings[entryName] = entry
		} else if !entryPresencePairsComplete && filter == common.EntryMappingFilterIncompleteEntry {
			filteredEntryMappings[entryName] = entry
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
