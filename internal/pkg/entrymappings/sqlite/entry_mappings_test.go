package sqlite

import (
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/almanac1631/scrubarr/internal/pkg/utils"
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
		entryMappings, err = instance.parseEntryMapping("Some Film", "69914146", utils.ParseTime("2025-02-17T19:32:07Z"), 1337, entryMappings)
		entryMappings, err = instance.parseEntryMapping("Another film", "7b74ad25", utils.ParseTime("2025-03-18T09:32:01Z"), 19232, entryMappings)
		assert.Nil(t, err)
		assert.Equal(t, []*common.EntryMapping{
			{"Some Film", utils.ParseTime("2025-02-17T19:32:07Z"), 1337, []common.RetrieverInfo{retriever1}},
			{"Another film", utils.ParseTime("2025-03-18T09:32:01Z"), 19232, []common.RetrieverInfo{retriever2}},
		}, entryMappings)
	})

	t.Run("parses multi-retriever entry mappings", func(t *testing.T) {
		instanceCopy := instancePreset
		instance := &instanceCopy
		var entryMappings []*common.EntryMapping
		var err error
		entryMappings, err = instance.parseEntryMapping("Some Film", "69914146", utils.ParseTime("2025-02-17T22:00:00Z"), 1337, entryMappings)
		entryMappings, err = instance.parseEntryMapping("Some Film", "7b74ad25", utils.ParseTime("2025-02-17T19:00:00Z"), 101337, entryMappings)
		entryMappings, err = instance.parseEntryMapping("Some Film", "49192bf4", utils.ParseTime("2025-02-17T21:00:00Z"), 9337, entryMappings)
		entryMappings, err = instance.parseEntryMapping("Another film", "7b74ad25", utils.ParseTime("2025-03-18T09:32:01Z"), 19232, entryMappings)
		assert.Nil(t, err)
		assert.Equal(t, []*common.EntryMapping{
			{"Some Film", utils.ParseTime("2025-02-17T19:00:00Z"), 101337, []common.RetrieverInfo{retriever1, retriever2, retriever3}},
			{"Another film", utils.ParseTime("2025-03-18T09:32:01Z"), 19232, []common.RetrieverInfo{retriever2}},
		}, entryMappings)
	})
}
