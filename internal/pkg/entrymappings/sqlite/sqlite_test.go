package sqlite

import (
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntryMappingManager_parseEntryMapping(t *testing.T) {
	// id: 69914146
	retriever1 := common.RetrieverInfo{Category: "c1", SoftwareName: "s1", Name: "r1"}
	// id: 7b74ad25
	retriever2 := common.RetrieverInfo{Category: "c2", SoftwareName: "s2", Name: "r2"}
	// id: 49192bf4
	retriever3 := common.RetrieverInfo{Category: "c3", SoftwareName: "s3", Name: "r3"}

	instancePreset := EntryMappingManager{
		entryRetrievers: map[common.RetrieverInfo]common.EntryRetriever{
			retriever1: nil,
			retriever2: nil,
			retriever3: nil,
		},
	}

	t.Run("parses simple single-retriever entry mappings", func(t *testing.T) {
		instanceCopy := instancePreset
		instance := &instanceCopy
		var entryMappings []*common.EntryMapping
		var err error
		entryMappings, err = instance.parseEntryMapping("Some Film", "69914146", entryMappings)
		entryMappings, err = instance.parseEntryMapping("Another film", "7b74ad25", entryMappings)
		assert.Nil(t, err)
		assert.Equal(t, []*common.EntryMapping{
			{"Some Film", []common.RetrieverInfo{retriever1}},
			{"Another film", []common.RetrieverInfo{retriever2}},
		}, entryMappings)
	})

	t.Run("parses multi-retriever entry mappings", func(t *testing.T) {
		instanceCopy := instancePreset
		instance := &instanceCopy
		var entryMappings []*common.EntryMapping
		var err error
		entryMappings, err = instance.parseEntryMapping("Some Film", "69914146", entryMappings)
		entryMappings, err = instance.parseEntryMapping("Some Film", "49192bf4", entryMappings)
		entryMappings, err = instance.parseEntryMapping("Another film", "7b74ad25", entryMappings)
		assert.Nil(t, err)
		assert.Equal(t, []*common.EntryMapping{
			{"Some Film", []common.RetrieverInfo{retriever1, retriever3}},
			{"Another film", []common.RetrieverInfo{retriever2}},
		}, entryMappings)
	})
}
