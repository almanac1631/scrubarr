package media

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/almanac1631/scrubarr/pkg/domain"
)

var _ domain.MediaSourceManager = (*DefaultMediaManager)(nil)

type DefaultMediaManager struct {
	Entries    map[domain.MediaType][]domain.MediaEntry
	entryLock  *sync.Mutex
	retrievers map[domain.MediaType]domain.MediaSource
}

func NewDefaultMediaManager(retrievers ...domain.MediaSource) *DefaultMediaManager {
	manager := &DefaultMediaManager{nil, &sync.Mutex{}, make(map[domain.MediaType]domain.MediaSource)}
	for _, retriever := range retrievers {
		manager.retrievers[retriever.SupportedMediaType()] = retriever
	}
	return manager
}

func (manager *DefaultMediaManager) GetMedia() ([]*domain.MediaEntry, error) {
	if manager.Entries == nil {
		if err := manager.RefreshCache(); err != nil {
			return nil, err
		}
	}

	mediaList := make([]*domain.MediaEntry, 0)
	for _, mediaEntries := range manager.Entries {
		for _, mediaEntry := range mediaEntries {
			mediaList = append(mediaList, &mediaEntry)
		}
	}
	return mediaList, nil
}

func (manager *DefaultMediaManager) DeleteMediaFiles(mediaType domain.MediaType, fileIds []int64, stopParentMonitoring bool) error {
	manager.entryLock.Lock()
	defer manager.entryLock.Unlock()
	retriever, ok := manager.retrievers[mediaType]
	if !ok {
		return fmt.Errorf("could not find retriever for media type %q", mediaType)
	}
	return retriever.DeleteMediaFiles(fileIds, stopParentMonitoring)
}

func (manager *DefaultMediaManager) RefreshCache() error {
	manager.entryLock.Lock()
	defer manager.entryLock.Unlock()
	manager.Entries = make(map[domain.MediaType][]domain.MediaEntry)
	for mediaType, retriever := range manager.retrievers {
		var err error
		manager.Entries[mediaType], err = retriever.GetMedia()
		if err != nil {
			return err
		}
	}
	return nil
}

func (manager *DefaultMediaManager) SaveCache(writer io.Writer) error {
	manager.entryLock.Lock()
	defer manager.entryLock.Unlock()
	return json.NewEncoder(writer).Encode(manager.Entries)
}

func (manager *DefaultMediaManager) LoadCache(reader io.ReadSeeker) error {
	manager.entryLock.Lock()
	defer manager.entryLock.Unlock()
	manager.Entries = make(map[domain.MediaType][]domain.MediaEntry)
	return json.NewDecoder(reader).Decode(&manager.Entries)
}
