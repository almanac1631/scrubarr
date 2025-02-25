package simple

import (
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"log/slog"
	"strings"
)

func BundledEntryRetriever(fileEndings []string) func(entryRetrievers map[common.RetrieverInfo]common.EntryRetriever) (map[common.RetrieverInfo]common.RetrieverEntries, error) {
	fn := func(entryRetrievers map[common.RetrieverInfo]common.EntryRetriever) (map[common.RetrieverInfo]common.RetrieverEntries, error) {
		entriesCombined := make(map[common.RetrieverInfo]common.RetrieverEntries)
		resultChan := make(chan *retrieverResult, len(entryRetrievers))
		for info, retriever := range entryRetrievers {
			go func() {
				result := &retrieverResult{
					info:    info,
					entries: nil,
				}
				defer func() {
					resultChan <- result
				}()
				logger := slog.With("retrieverInfo", info, "retrieverId", info.Id())
				logger.Debug("Retrieving entries...")
				respEntries, err := retriever.RetrieveEntries()
				if err != nil {
					logger.Error("Could not retrieve entries from retriever", "error", err)
					return
				}
				result.entries = map[common.EntryName]common.Entry{}
				for entryName, entry := range respEntries {
					result.entries[convertEntryName(entryName, fileEndings)] = entry
				}
			}()
		}
		for i := 0; i < len(entryRetrievers); i++ {
			result := <-resultChan
			entriesCombined[result.info] = result.entries
		}
		close(resultChan)
		return entriesCombined, nil
	}
	return fn
}

func convertEntryName(rawEntryName common.EntryName, fileEndings []string) common.EntryName {
	lowerEntryName := strings.ToLower(string(rawEntryName))
	for _, fileEnding := range fileEndings {
		index := strings.LastIndex(lowerEntryName, fileEnding)
		if index == len(lowerEntryName)-len(fileEnding) {
			return rawEntryName[:index]
		}
	}
	return rawEntryName
}

type retrieverResult struct {
	info    common.RetrieverInfo
	entries common.RetrieverEntries
}
