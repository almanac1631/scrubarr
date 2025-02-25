package simple

import (
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"log/slog"
)

func BundledEntryRetriever(entryRetrievers map[common.RetrieverInfo]common.EntryRetriever) (map[common.RetrieverInfo]common.RetrieverEntries, error) {
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
			var err error
			result.entries, err = retriever.RetrieveEntries()
			if err != nil {
				logger.Error("Could not retrieve entries from retriever", "error", err)
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

type retrieverResult struct {
	info    common.RetrieverInfo
	entries common.RetrieverEntries
}
