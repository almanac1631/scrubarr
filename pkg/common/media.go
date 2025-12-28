package common

import "github.com/almanac1631/scrubarr/pkg/media"

type MovieInfo struct {
	media.Movie
	ExistsInTorrentClient bool
}

type Manager interface {
	GetMovieInfos(page int) (movies []MovieInfo, hasNext bool, err error)
}
