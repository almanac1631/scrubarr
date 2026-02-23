package linker

import (
	"path/filepath"

	"github.com/almanac1631/scrubarr/pkg/domain"
	"github.com/almanac1631/scrubarr/pkg/inventory"
)

var _ inventory.Linker = (*Service)(nil)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s Service) LinkMedia(mediaEntries []*domain.MediaEntry, torrentEntries []*domain.TorrentEntry) ([]inventory.LinkedMedia, error) {
	linkedMedias := make([]inventory.LinkedMedia, 0)
	for _, mediaEntry := range mediaEntries {
		linkedMedia := &inventory.LinkedMedia{
			MediaMetadata: mediaEntry.MediaMetadata,
		}
		for _, mediaFile := range mediaEntry.Files {
			linkedMediaFile, found := searchLinkedTorrentEntry(mediaFile, torrentEntries)
			if !found {
				linkedMediaFile = inventory.LinkedMediaFile{MediaFile: mediaFile}
			}
			if linkedMedia.Files == nil {
				linkedMedia.Files = []inventory.LinkedMediaFile{linkedMediaFile}
			} else {
				linkedMedia.Files = append(linkedMedia.Files, linkedMediaFile)
			}
		}
		if linkedMedia.Files != nil {
			linkedMedias = append(linkedMedias, *linkedMedia)
		}
	}
	return linkedMedias, nil
}

func searchLinkedTorrentEntry(mediaFile domain.MediaFile, torrentEntries []*domain.TorrentEntry) (inventory.LinkedMediaFile, bool) {
	for _, torrentEntry := range torrentEntries {
		linkedMediaFile := inventory.LinkedMediaFile{
			MediaFile:    mediaFile,
			TorrentEntry: torrentEntry,
		}

		if torrentEntry.Name == mediaFile.OriginalFilePath {
			return linkedMediaFile, true
		}

		torrentEntryNameWithExt := torrentEntry.Name + filepath.Ext(mediaFile.OriginalFilePath)
		if torrentEntryNameWithExt == mediaFile.OriginalFilePath {
			return linkedMediaFile, true
		}

		for _, torrentEntryFile := range torrentEntry.Files {
			if torrentEntryFile.Size != mediaFile.Size {
				continue
			}
			if torrentEntryFile.Path == mediaFile.OriginalFilePath {
				return linkedMediaFile, true
			}
			torrentFileBase := filepath.Base(torrentEntryFile.Path)
			if torrentFileBase == mediaFile.OriginalFilePath {
				return linkedMediaFile, true
			}
		}
	}
	return inventory.LinkedMediaFile{}, false
}
