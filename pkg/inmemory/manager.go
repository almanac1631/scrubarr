package inmemory

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/almanac1631/scrubarr/pkg/common"
	"github.com/almanac1631/scrubarr/pkg/media"
	"github.com/almanac1631/scrubarr/pkg/torrentclients"
)

type Manager struct {
	mappedMoviesCache []common.MovieInfo

	radarrRetriever   *media.RadarrRetriever
	delugeRetriever   *torrentclients.DelugeRetriever
	rtorrentRetriever *torrentclients.RtorrentRetriever
}

func NewManager(radarrRetriever *media.RadarrRetriever, delugeRetriever *torrentclients.DelugeRetriever, rtorrentRetriever *torrentclients.RtorrentRetriever) *Manager {
	return &Manager{
		nil, radarrRetriever, delugeRetriever, rtorrentRetriever,
	}
}

const pageSize = 10

func (m *Manager) GetMovieInfos(page int, sortInfo common.SortInfo) ([]common.MovieInfo, bool, error) {
	if m.mappedMoviesCache == nil {
		radarrMovies, err := m.radarrRetriever.GetMovies()
		if err != nil {
			return nil, false, fmt.Errorf("failed to get movies from radarr: %w", err)
		}
		m.mappedMoviesCache = make([]common.MovieInfo, 0, len(radarrMovies))
		for _, movie := range radarrMovies {
			exists, err := m.delugeRetriever.SearchForMovie(movie.OriginalFilePath)
			if err != nil {
				return nil, false, fmt.Errorf("failed to search for movie in deluge: %w", err)
			}
			if !exists {
				exists, err = m.rtorrentRetriever.SearchForMovie(movie.OriginalFilePath)
				if err != nil {
					return nil, false, fmt.Errorf("failed to search for movie in rtorrent: %w", err)
				}
			}
			m.mappedMoviesCache = append(m.mappedMoviesCache, common.MovieInfo{
				Movie:                 movie,
				ExistsInTorrentClient: exists,
			})
		}
	}
	hasNext := false
	movies := make([]common.MovieInfo, len(m.mappedMoviesCache))
	copy(movies, m.mappedMoviesCache)
	slices.SortFunc(movies, func(a, b common.MovieInfo) int {
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
			result = CompareBool(a.ExistsInTorrentClient, b.ExistsInTorrentClient)
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

func CompareBool(a, b bool) int {
	if a == b {
		return 0
	}
	if a {
		return 1
	}
	return -1
}
