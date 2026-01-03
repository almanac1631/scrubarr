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
				MediaPart:             movie.Parts[0],
				ExistsInTorrentClient: finding != nil,
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
			parts = append(parts, common.MatchedMediaPart{
				MediaPart:             part,
				ExistsInTorrentClient: finding != nil,
			})
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
				doesNotExist := !part.ExistsInTorrentClient
				return doesNotExist
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
	if m.matchedMediaCache == nil {
		if err := m.refreshCache(); err != nil {
			return nil, err
		}
	}
	filteredMediaList := make([]common.MatchedMedia, 0)
	for _, mediaEntry := range m.matchedMediaCache {
		if mediaEntry.Type != common.MediaTypeSeries || mediaEntry.Id != seriesId {
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
