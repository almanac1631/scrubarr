package torrentclients

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"slices"
	"sync"

	"github.com/almanac1631/scrubarr/pkg/common"
)

var _ common.TorrentClientManager = (*DefaultTorrentManager)(nil)
var _ common.CachedRetriever = (*DefaultTorrentManager)(nil)

type DefaultTorrentManager struct {
	Entries    map[string][]*common.TorrentEntry
	entryLock  *sync.Mutex
	retrievers map[string]common.TorrentClientRetriever
}

func NewDefaultTorrentManager(retrievers ...common.TorrentClientRetriever) *DefaultTorrentManager {
	manager := &DefaultTorrentManager{
		retrievers: make(map[string]common.TorrentClientRetriever),
		entryLock:  new(sync.Mutex),
	}
	for _, retriever := range retrievers {
		manager.retrievers[retriever.Name()] = retriever
	}
	return manager
}

func (manager *DefaultTorrentManager) SearchForMedia(originalFilePath string, size int64) (finding *common.TorrentEntry, err error) {
	for _, entries := range manager.Entries {
		for _, entry := range entries {
			if matches(entry, originalFilePath, size) {
				return entry, nil
			}
		}
	}
	return nil, err
}

func matches(entry *common.TorrentEntry, originalFilePath string, size int64) bool {
	if entry.Name == originalFilePath {
		return true
	}
	torrentNameWithExt := entry.Name + filepath.Ext(originalFilePath)
	if torrentNameWithExt == originalFilePath {
		return true
	}
	for _, file := range entry.Files {
		if file.Size != size {
			continue
		}
		if file.Path == originalFilePath {
			return true
		}
		fileNameCmp := filepath.Base(file.Path)
		if fileNameCmp == originalFilePath {
			return true
		}
	}
	return false
}

func (manager *DefaultTorrentManager) DeleteFinding(client, id string) error {
	manager.entryLock.Lock()
	defer manager.entryLock.Unlock()
	retriever, ok := manager.retrievers[client]
	if !ok {
		return fmt.Errorf("could not find retriever for %q", client)
	}
	if err := retriever.DeleteTorrent(id); err != nil {
		return fmt.Errorf("could not delete torrent %q from client %q: %w", id, client, err)
	}
	manager.Entries[client] = slices.DeleteFunc(manager.Entries[client], func(entrySearch *common.TorrentEntry) bool {
		return entrySearch.Id == id
	})
	return nil
}

func (manager *DefaultTorrentManager) RefreshCache() error {
	manager.entryLock.Lock()
	defer manager.entryLock.Unlock()
	manager.Entries = make(map[string][]*common.TorrentEntry)
	for name, retriever := range manager.retrievers {
		slog.Debug("Refreshing torrent cache...", "client", name)
		retrieverEntries, err := retriever.GetTorrentEntries()
		if err != nil {
			return fmt.Errorf("could not get torrent entries for client %q: %w", name, err)
		}
		manager.Entries[name] = retrieverEntries
	}
	return nil
}

func (manager *DefaultTorrentManager) SaveCache(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(manager.Entries)
}

func (manager *DefaultTorrentManager) LoadCache(reader io.ReadSeeker) error {
	manager.Entries = make(map[string][]*common.TorrentEntry)
	return json.NewDecoder(reader).Decode(&manager.Entries)
}
