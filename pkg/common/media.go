package common

import "github.com/almanac1631/scrubarr/pkg/media"

type MovieInfo struct {
	media.Movie
	ExistsInTorrentClient bool
}

type SortKey string

const (
	SortKeyName          SortKey = "name"
	SortKeySize          SortKey = "size"
	SortKeyAdded         SortKey = "added"
	SortKeyTorrentStatus SortKey = "torrent_status"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type SortInfo struct {
	Key   SortKey
	Order SortOrder
}

type Manager interface {
	GetMovieInfos(page int, sortInfo SortInfo) (movies []MovieInfo, hasNext bool, err error)
}
