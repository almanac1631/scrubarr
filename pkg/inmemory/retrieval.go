package inmemory

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/almanac1631/scrubarr/pkg/common"
)

func (m *Manager) GetMatchedMedia(page int, sortInfo common.SortInfo) ([]common.MatchedEntry, bool, error) {
	if m.matchedEntriesCache == nil {
		if err := m.refreshCache(); err != nil {
			return nil, false, err
		}
	}
	hasNext := false
	matchedEntries := make([]common.MatchedEntry, len(m.matchedEntriesCache))
	copy(matchedEntries, m.matchedEntriesCache)

	torrentStatusScores := map[string]int{}
	for _, entry := range m.matchedEntriesCache {
		totalScore := 0
		for _, part := range entry.Parts {
			totalScore += part.TorrentInformation.GetScore()
		}
		torrentStatusScores[entry.Url] = totalScore / len(entry.Parts)
	}

	slices.SortFunc(matchedEntries, func(a, b common.MatchedEntry) int {
		var result int
		switch sortInfo.Key {
		case common.SortKeyName:
			result = strings.Compare(strings.ToLower(a.Title), strings.ToLower(b.Title))
			break
		case common.SortKeySize:
			result = cmp.Compare(a.Size, b.Size)
			break
		case common.SortKeyAdded:
			result = cmp.Compare(a.Added.Unix(), b.Added.Unix())
			break
		case common.SortKeyTorrentStatus:
			result = cmp.Compare(torrentStatusScores[a.Url], torrentStatusScores[b.Url])
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
	if pageSize*page < len(matchedEntries) {
		hasNext = true
		matchedEntries = matchedEntries[pageSize*(page-1) : pageSize*page]
	} else {
		matchedEntries = matchedEntries[pageSize*(page-1):]
	}
	return matchedEntries, hasNext, nil
}

func (m *Manager) GetMatchedMediaBySeriesId(seriesId int64) (media common.MatchedEntry, err error) {
	return m.getSingleMatchedMediaEntry(common.MediaTypeSeries, seriesId)
}

func (m *Manager) getSingleMatchedMediaEntry(mediaType common.MediaType, id int64) (media common.MatchedEntry, err error) {
	matchedMediaList, err := m.getFilteredMatchedMediaFunc(func(media common.MatchedEntry) bool {
		return media.Type == mediaType && media.Id == id
	})
	if err != nil {
		return common.MatchedEntry{}, err
	}
	if len(matchedMediaList) == 0 {
		return common.MatchedEntry{}, common.ErrMediaNotFound
	} else if len(matchedMediaList) > 1 {
		return common.MatchedEntry{}, fmt.Errorf("multiple matched media found with type %s and id %d", mediaType, id)
	}
	return matchedMediaList[0], nil
}

func (m *Manager) getFilteredMatchedMediaFunc(filterFunc func(media common.MatchedEntry) bool) (media []common.MatchedEntry, err error) {
	if m.matchedEntriesCache == nil {
		if err := m.refreshCache(); err != nil {
			return nil, err
		}
	}
	filteredMediaList := make([]common.MatchedEntry, 0)
	for _, mediaEntry := range m.matchedEntriesCache {
		if filterFunc(mediaEntry) {
			filteredMediaList = append(filteredMediaList, mediaEntry)
		}
	}
	return filteredMediaList, nil
}
