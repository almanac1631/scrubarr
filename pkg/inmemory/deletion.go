package inmemory

import (
	"fmt"
	"log/slog"
	"maps"
	"slices"

	"github.com/almanac1631/scrubarr/pkg/common"
)

func (m *Manager) DeleteMedia(mediaType common.MediaType, id int64) error {
	matchedMedia, err := m.getSingleMatchedMediaEntry(mediaType, id)
	if err != nil {
		return err
	}
	return m.deleteMediaParts(id, mediaType, matchedMedia.Parts)
}

func (m *Manager) DeleteSeason(id int64, season int) error {
	filteredMatchedMedia, err := m.GetMatchedMediaBySeriesId(id)
	if err != nil {
		return err
	}
	seasonParts := make([]common.MatchedEntryPart, 0)
	for _, part := range filteredMatchedMedia.Parts {
		if part.MediaPart.Season != season {
			continue
		}
		seasonParts = append(seasonParts, part)
	}
	return m.deleteMediaParts(id, common.MediaTypeSeries, seasonParts)
}

func (m *Manager) deleteMediaParts(mediaId int64, mediaType common.MediaType, parts []common.MatchedEntryPart) error {
	type torrent struct {
		client string
		id     string
	}
	torrentsToDelete := make([]torrent, 0)
	fileIdsToDeleteMap := make(map[int64]struct{})
	for _, part := range parts {
		if _, ok := fileIdsToDeleteMap[part.MediaPart.Id]; !ok {
			fileIdsToDeleteMap[part.MediaPart.Id] = struct{}{}
		}
		if part.TorrentInformation.Status == common.TorrentStatusMissing {
			continue
		}
		partTorrent := torrent{
			client: part.TorrentInformation.Client,
			id:     part.TorrentInformation.Id,
		}
		if !slices.ContainsFunc(torrentsToDelete, func(compareTorrent torrent) bool {
			if part.TorrentInformation.Status == common.TorrentStatusMissing {
				return false
			}
			return partTorrent.client == compareTorrent.client && partTorrent.id == compareTorrent.id
		}) {
			torrentsToDelete = append(torrentsToDelete, partTorrent)
		}
	}
	for _, torrentToDelete := range torrentsToDelete {
		slog.Debug("Deleting torrent...", "client", torrentToDelete.client, "torrentId", torrentToDelete.id)
		if err := m.torrentManager.DeleteFinding(torrentToDelete.client, torrentToDelete.id); err != nil {
			return fmt.Errorf("could not delete %s (id: %d) from torrent client %q (torrent id: %q): %w",
				mediaType, mediaId, torrentToDelete.client, torrentToDelete.id, err)
		}
	}
	fileIdsToDelete := slices.Collect(maps.Keys(fileIdsToDeleteMap))
	if err := m.mediaManager.DeleteMediaFiles(mediaType, fileIdsToDelete, true); err != nil {
		return err
	}
	slog.Debug("Successfully deleted delete media files from media manager.", "mediaId", mediaId)
	mediaIndex := slices.IndexFunc(m.matchedEntriesCache, func(media common.MatchedEntry) bool {
		return media.Type == mediaType && media.Id == mediaId
	})
	if mediaIndex == -1 {
		slog.Warn("No matched media found for this media.", "mediaType", mediaType, "mediaId", mediaId)
		return nil
	}
	mediaEntry := m.matchedEntriesCache[mediaIndex]
	for partId := range fileIdsToDeleteMap {
		index := slices.IndexFunc(mediaEntry.Parts, func(mediaPart common.MatchedEntryPart) bool {
			return mediaPart.MediaPart.Id == partId
		})
		if index == -1 {
			continue
		}
		mediaEntry.Parts = append(mediaEntry.Parts[:index], mediaEntry.Parts[index+1:]...)
	}
	if len(mediaEntry.Parts) > 0 {
		m.matchedEntriesCache[mediaIndex] = mediaEntry
	} else {
		m.matchedEntriesCache = append(m.matchedEntriesCache[:mediaIndex], m.matchedEntriesCache[mediaIndex+1:]...)
	}
	return nil
}
