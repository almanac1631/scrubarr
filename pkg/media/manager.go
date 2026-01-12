package media

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/almanac1631/scrubarr/pkg/common"
)

var _ common.MediaManager = (*DefaultMediaManager)(nil)

type DefaultMediaManager struct {
	Entries    map[common.MediaType][]common.Media
	entryLock  *sync.Mutex
	retrievers map[common.MediaType]common.MediaRetriever
}

func NewDefaultMediaManager(retrievers ...common.MediaRetriever) *DefaultMediaManager {
	manager := &DefaultMediaManager{nil, &sync.Mutex{}, make(map[common.MediaType]common.MediaRetriever)}
	for _, retriever := range retrievers {
		manager.retrievers[retriever.SupportedMediaType()] = retriever
	}
	return manager
}

func (manager *DefaultMediaManager) GetMedia() ([]common.Media, error) {
	if manager.Entries == nil {
		if err := manager.RefreshCache(); err != nil {
			return nil, err
		}
	}
	mediaList := make([]common.Media, 0)
	for _, mediaEntries := range manager.Entries {
		mediaList = append(mediaList, mediaEntries...)
	}
	return mediaList, nil
}

func (manager *DefaultMediaManager) DeleteMediaFiles(mediaType common.MediaType, fileIds []int64, stopParentMonitoring bool) error {
	retriever, ok := manager.retrievers[mediaType]
	if !ok {
		return fmt.Errorf("could not find retriever for media type %q", mediaType)
	}
	return retriever.DeleteMediaFiles(fileIds, stopParentMonitoring)
}

func (manager *DefaultMediaManager) RefreshCache() error {
	manager.entryLock.Lock()
	defer manager.entryLock.Unlock()
	manager.Entries = make(map[common.MediaType][]common.Media)
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
	manager.Entries = make(map[common.MediaType][]common.Media)
	return json.NewDecoder(reader).Decode(&manager.Entries)
}
