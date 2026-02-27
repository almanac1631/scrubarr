package torrentclients

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"sync"

	"github.com/almanac1631/scrubarr/pkg/domain"
)

var _ domain.TorrentSourceManager = (*DefaultTorrentManager)(nil)

type DefaultTorrentManager struct {
	Entries    map[string][]*domain.TorrentEntry
	entryLock  *sync.Mutex
	retrievers map[string]domain.TorrentSource
}

func NewDefaultTorrentManager(retrievers ...domain.TorrentSource) *DefaultTorrentManager {
	manager := &DefaultTorrentManager{
		retrievers: make(map[string]domain.TorrentSource),
		entryLock:  new(sync.Mutex),
	}
	for _, retriever := range retrievers {
		manager.retrievers[retriever.Name()] = retriever
	}
	return manager
}

func (manager *DefaultTorrentManager) GetTorrents() ([]*domain.TorrentEntry, error) {
	if manager.Entries == nil {
		if err := manager.RefreshCache(); err != nil {
			return nil, err
		}
	}

	torrentList := make([]*domain.TorrentEntry, 0)
	for _, torrentEntries := range manager.Entries {
		for _, torrentEntry := range torrentEntries {
			torrentList = append(torrentList, torrentEntry)
		}
	}
	return torrentList, nil
}

func (manager *DefaultTorrentManager) DeleteTorrent(client, id string) error {
	manager.entryLock.Lock()
	defer manager.entryLock.Unlock()
	retriever, ok := manager.retrievers[client]
	if !ok {
		return fmt.Errorf("could not find retriever for %q", client)
	}
	if err := retriever.DeleteTorrent(id); err != nil {
		return fmt.Errorf("could not delete torrent %q from client %q: %w", id, client, err)
	}
	manager.Entries[client] = slices.DeleteFunc(manager.Entries[client], func(entrySearch *domain.TorrentEntry) bool {
		return entrySearch.Id == id
	})
	return nil
}

func (manager *DefaultTorrentManager) RefreshCache() error {
	manager.entryLock.Lock()
	defer manager.entryLock.Unlock()
	entryLock := &sync.Mutex{}
	manager.Entries = make(map[string][]*domain.TorrentEntry)
	errChan := make(chan error)
	defer close(errChan)
	for name, retriever := range manager.retrievers {
		go func() {
			slog.Debug("Refreshing torrent cache...", "client", name)
			retrieverEntries, err := retriever.GetTorrentEntries()
			if err == nil {
				entryLock.Lock()
				defer entryLock.Unlock()
				manager.Entries[name] = retrieverEntries
			}
			if err != nil {
				err = fmt.Errorf("could not get torrent entries for client %q: %w", name, err)
			}
			errChan <- err
			slog.Debug("Refreshed torrent cache", "client", name)
		}()
	}
	var err error
	for i := 0; i < len(manager.retrievers); i++ {
		err = errors.Join(err, <-errChan)
	}
	return err
}

func (manager *DefaultTorrentManager) SaveCache(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(manager.Entries)
}

func (manager *DefaultTorrentManager) LoadCache(reader io.ReadSeeker) error {
	manager.Entries = make(map[string][]*domain.TorrentEntry)
	return json.NewDecoder(reader).Decode(&manager.Entries)
}
