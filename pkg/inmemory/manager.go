package inmemory

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/almanac1631/scrubarr/pkg/common"
	"github.com/almanac1631/scrubarr/pkg/media"
)

var _ common.Manager = (*Manager)(nil)

type Manager struct {
	matchedMediaCache []common.MatchedMedia

	radarrRetriever *media.RadarrRetriever
	sonarrRetriever *media.SonarrRetriever

	torrentManager common.TorrentClientManager
}

func NewManager(radarrRetriever *media.RadarrRetriever, sonarrRetriever *media.SonarrRetriever, torrentManager common.TorrentClientManager) *Manager {
	return &Manager{
		nil, radarrRetriever, sonarrRetriever, torrentManager,
	}
}

const pageSize = 10

func (m *Manager) refreshCache() error {
	radarrMovies, err := m.radarrRetriever.GetMovies()
	if err != nil {
		return fmt.Errorf("failed to get movies from radarr: %w", err)
	}
	m.matchedMediaCache = make([]common.MatchedMedia, 0, len(radarrMovies))
	for _, movie := range radarrMovies {
		originalFilePath := movie.Parts[0].OriginalFilePath
		finding, err := m.torrentManager.SearchForMedia(originalFilePath, movie.Parts[0].Size)
		if err != nil {
			return err
		}
		if finding != nil {
			movie.Added = finding.Added
		}
		m.matchedMediaCache = append(m.matchedMediaCache, common.MatchedMedia{
			MediaMetadata: movie.MediaMetadata,
			Parts: []common.MatchedMediaPart{{
				MediaPart:      movie.Parts[0],
				TorrentFinding: finding,
			}},
		})
	}
	sonarrSeries, err := m.sonarrRetriever.GetMedia()
	for _, mediaEntry := range sonarrSeries {
		parts := make([]common.MatchedMediaPart, 0, len(mediaEntry.Parts))
		added := mediaEntry.Added
		for _, part := range mediaEntry.Parts {
			finding, err := m.torrentManager.SearchForMedia(part.OriginalFilePath, part.Size)
			if err != nil {
				return err
			}
			if finding != nil && !finding.Added.IsZero() && finding.Added.Before(added) {
				added = finding.Added
			}
			mediaPart := common.MatchedMediaPart{
				MediaPart:      part,
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
			result = CompareBool(existsInTorrentClient(a), existsInTorrentClient(b))
			break
		default:
			slog.Error("received unknown sort key", "sortKey", sortInfo.Key)
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

func (m *Manager) GetMatchedMediaBySeriesId(seriesId int64) (media []common.MatchedMedia, err error) {
	return m.getFilteredMatchedMedia(common.MediaTypeSeries, seriesId)
}

func (m *Manager) getFilteredMatchedMedia(mediaType common.MediaType, id int64) (media []common.MatchedMedia, err error) {
	if m.matchedMediaCache == nil {
		if err := m.refreshCache(); err != nil {
			return nil, err
		}
	}
	filteredMediaList := make([]common.MatchedMedia, 0)
	for _, mediaEntry := range m.matchedMediaCache {
		if mediaEntry.Type != mediaType || mediaEntry.Id != id {
			continue
		}
		filteredMediaList = append(filteredMediaList, mediaEntry)
	}
	return filteredMediaList, nil
}

func CompareBool(a, b bool) int {
	if a == b {
		return 0
	}
	if a {
		return 1
	}
	return -1
}

func (m *Manager) DeleteMedia(mediaType common.MediaType, id int64) error {
	filteredMatchedMedia, err := m.getFilteredMatchedMedia(mediaType, id)
	if err != nil {
		return err
	}
	type torrent struct {
		client string
		id     string
	}
	torrentsToDelete := make([]torrent, 0)
	for _, matchedMedia := range filteredMatchedMedia {
		for _, part := range matchedMedia.Parts {
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
	}
	for _, torrentToDelete := range torrentsToDelete {
		//if err = m.torrentManager.DeleteFinding(torrentToDelete.client, torrentToDelete.id); err != nil {
		//	return fmt.Errorf("could not delete %s with id %d from torrent client %q (torrent id: %q): %w",
		//		mediaType, id, torrentToDelete.client, torrentToDelete.id, err)
		//}
		slog.Info("would delete torrent", "client", torrentToDelete.client, "torrentId", torrentToDelete.id)
	}
	return nil
}
