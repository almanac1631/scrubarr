package scrubarr

import (
	"github.com/almanac1631/scrubarr/internal/app/common"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval"
)

type cachedRetrieverRegistry struct {
	cachedEntryMappings map[retrieval.EntryName]common.EntryPresencePairs
	retrieverRegistry   common.MapBasedRetrieverRegistry
}

func (c *cachedRetrieverRegistry) RetrieveEntryMapping() map[retrieval.EntryName]common.EntryPresencePairs {
	return c.cachedEntryMappings
}

func (c *cachedRetrieverRegistry) RefreshCachedEntryMapping() {
	c.cachedEntryMappings = c.retrieverRegistry.RetrieveEntryMapping()
}

func (c *cachedRetrieverRegistry) GetRetrievers() []common.RetrieverInfo {
	return c.retrieverRegistry.GetRetrievers()
}
