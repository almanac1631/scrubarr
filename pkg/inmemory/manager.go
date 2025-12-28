package inmemory

import (
	"fmt"

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

func (m *Manager) GetMovieInfos(page int) ([]common.MovieInfo, bool, error) {
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
	movies := make([]common.MovieInfo, 0, pageSize)
	if pageSize*page < len(m.mappedMoviesCache) {
		hasNext = true
		movies = m.mappedMoviesCache[pageSize*(page-1) : pageSize*page]
	} else {
		movies = m.mappedMoviesCache[pageSize*(page-1):]
	}
	return movies, hasNext, nil
}
