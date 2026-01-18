package inmemory

import (
	"cmp"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strings"

	"github.com/almanac1631/scrubarr/pkg/common"
)

var _ common.Manager = (*Manager)(nil)

type Manager struct {
	matchedMediaCache []common.MatchedMedia

	mediaManager common.MediaManager

	torrentManager common.TorrentClientManager

	trackerManager common.TrackerManager
}

func NewManager(mediaManager common.MediaManager, torrentManager common.TorrentClientManager, trackerManager common.TrackerManager) *Manager {
	return &Manager{
		nil, mediaManager, torrentManager, trackerManager,
	}
}

const pageSize = 10

func (m *Manager) refreshCache() error {
	m.matchedMediaCache = make([]common.MatchedMedia, 0)
	medias, err := m.mediaManager.GetMedia()
	if err != nil {
		return fmt.Errorf("could not retrieve media from media manager: %w", err)
	}
	for _, mediaEntry := range medias {
		parts := make([]common.MatchedMediaPart, 0, len(mediaEntry.Parts))
		added := mediaEntry.Added
		for _, part := range mediaEntry.Parts {
			finding, err := m.torrentManager.SearchForMedia(part.OriginalFilePath, part.Size)
			if err != nil {
				return err
			}
			if finding != nil && !finding.Added.IsZero() && finding.Added.After(added) {
				added = finding.Added
			}
			var trackerName string
			if finding != nil {
				trackerName, err = m.trackerManager.GetTrackerName(finding.Trackers)
				if err != nil {
					if errors.Is(err, common.ErrTrackerNotFound) {
						slog.Warn("Could not find tracker name for media entry.",
							"mediaType", mediaEntry.Type, "mediaId", mediaEntry.Id, "part", finding.Name,
							"trackers", finding.Trackers, "findingId", finding.Id, "findingClient", finding.Client)
					} else {
						return err
					}
				}
			}
			mediaPart := common.MatchedMediaPart{
				MediaPart:      part,
				TrackerName:    trackerName,
				TorrentFinding: finding,
			}
			parts = append(parts, mediaPart)
		}
		matchedMedia := common.MatchedMedia{
			MediaMetadata: mediaEntry.MediaMetadata,
			Parts:         parts,
		}
		matchedMedia.Added = added
		m.matchedMediaCache = append(m.matchedMediaCache, matchedMedia)
	}
	return nil
}

func (m *Manager) GetMatchedMedia(page int, sortInfo common.SortInfo) ([]common.MatchedMedia, bool, error) {
	if m.matchedMediaCache == nil {
		if err := m.refreshCache(); err != nil {
			return nil, false, err
		}
	}
	hasNext := false
	movies := make([]common.MatchedMedia, len(m.matchedMediaCache))
	copy(movies, m.matchedMediaCache)

	totalSizeCache := map[string]int64{}
	totalSize := func(matchedMedia common.MatchedMedia) int64 {
		if _, ok := totalSizeCache[matchedMedia.Url]; !ok {
			totalSizeCalc := int64(0)
			for _, part := range matchedMedia.Parts {
				totalSizeCalc += part.Size
			}
			totalSizeCache[matchedMedia.Url] = totalSizeCalc
		}
		return totalSizeCache[matchedMedia.Url]
	}

	existsInTorrentClientCache := map[string]bool{}
	existsInTorrentClient := func(matchedMedia common.MatchedMedia) bool {
		if _, ok := existsInTorrentClientCache[matchedMedia.Url]; !ok {
			existsInTorrentClientCache[matchedMedia.Url] = !slices.ContainsFunc(matchedMedia.Parts, func(part common.MatchedMediaPart) bool {
				return part.TorrentFinding == nil
			})
		}
		return existsInTorrentClientCache[matchedMedia.Url]
	}

	slices.SortFunc(movies, func(a, b common.MatchedMedia) int {
		var result int
		switch sortInfo.Key {
		case common.SortKeyName:
			result = strings.Compare(strings.ToLower(a.Title), strings.ToLower(b.Title))
			break
		case common.SortKeySize:
			result = cmp.Compare(totalSize(a), totalSize(b))
			break
		case common.SortKeyAdded:
			result = cmp.Compare(a.Added.Unix(), b.Added.Unix())
			break
		case common.SortKeyTorrentStatus:
			result = compareBool(existsInTorrentClient(a), existsInTorrentClient(b))
			break
		default:
			slog.Error("Received unknown sort key.", "sortKey", sortInfo.Key)
			result = 0 // mark as incomparable
		}
		if sortInfo.Order == common.SortOrderDesc {
			result = -result
		}
		return result
	})
	if pageSize*page < len(movies) {
		hasNext = true
		movies = movies[pageSize*(page-1) : pageSize*page]
	} else {
		movies = movies[pageSize*(page-1):]
	}
	return movies, hasNext, nil
}

func (m *Manager) GetMatchedMediaBySeriesId(seriesId int64) (media common.MatchedMedia, err error) {
	return m.getSingleMatchedMediaEntry(common.MediaTypeSeries, seriesId)
}

func (m *Manager) getSingleMatchedMediaEntry(mediaType common.MediaType, id int64) (media common.MatchedMedia, err error) {
	matchedMediaList, err := m.getFilteredMatchedMediaFunc(func(media common.MatchedMedia) bool {
		return media.Type == mediaType && media.Id == id
	})
	if err != nil {
		return common.MatchedMedia{}, err
	}
	if len(matchedMediaList) == 0 {
		return common.MatchedMedia{}, common.ErrMediaNotFound
	} else if len(matchedMediaList) > 1 {
		return common.MatchedMedia{}, fmt.Errorf("multiple matched media found with type %s and id %d", mediaType, id)
	}
	return matchedMediaList[0], nil
}

func (m *Manager) getFilteredMatchedMediaFunc(filterFunc func(media common.MatchedMedia) bool) (media []common.MatchedMedia, err error) {
	if m.matchedMediaCache == nil {
		if err := m.refreshCache(); err != nil {
			return nil, err
		}
	}
	filteredMediaList := make([]common.MatchedMedia, 0)
	for _, mediaEntry := range m.matchedMediaCache {
		if filterFunc(mediaEntry) {
			filteredMediaList = append(filteredMediaList, mediaEntry)
		}
	}
	return filteredMediaList, nil
}

func compareBool(a, b bool) int {
	if a == b {
		return 0
	}
	if a {
		return 1
	}
	return -1
}

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
	seasonParts := make([]common.MatchedMediaPart, 0)
	for _, part := range filteredMatchedMedia.Parts {
		if part.Season != season {
			continue
		}
		seasonParts = append(seasonParts, part)
	}
	return m.deleteMediaParts(id, common.MediaTypeSeries, seasonParts)
}

func (m *Manager) deleteMediaParts(mediaId int64, mediaType common.MediaType, parts []common.MatchedMediaPart) error {
	type torrent struct {
		client string
		id     string
	}
	torrentsToDelete := make([]torrent, 0)
	fileIdsToDeleteMap := make(map[int64]struct{})
	for _, part := range parts {
		if _, ok := fileIdsToDeleteMap[part.Id]; !ok {
			fileIdsToDeleteMap[part.Id] = struct{}{}
		}
		if part.TorrentFinding == nil {
			continue
		}
		partTorrent := torrent{
			client: part.TorrentFinding.Client,
			id:     part.TorrentFinding.Id,
		}
		if !slices.ContainsFunc(torrentsToDelete, func(compareTorrent torrent) bool {
			if part.TorrentFinding == nil {
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
	mediaIndex := slices.IndexFunc(m.matchedMediaCache, func(media common.MatchedMedia) bool {
		return media.Type == mediaType && media.Id == mediaId
	})
	if mediaIndex == -1 {
		slog.Warn("No matched media found for this media.", "mediaType", mediaType, "mediaId", mediaId)
		return nil
	}
	mediaEntry := m.matchedMediaCache[mediaIndex]
	for partId := range fileIdsToDeleteMap {
		index := slices.IndexFunc(mediaEntry.Parts, func(mediaPart common.MatchedMediaPart) bool {
			return mediaPart.Id == partId
		})
		if index == -1 {
			continue
		}
		mediaEntry.Parts = append(mediaEntry.Parts[:index], mediaEntry.Parts[index+1:]...)
	}
	if len(mediaEntry.Parts) > 0 {
		m.matchedMediaCache[mediaIndex] = mediaEntry
	} else {
		m.matchedMediaCache = append(m.matchedMediaCache[:mediaIndex], m.matchedMediaCache[mediaIndex+1:]...)
	}
	return nil
}
