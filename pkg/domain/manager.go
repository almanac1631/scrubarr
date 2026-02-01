package domain

import "errors"

var ErrMediaNotFound = errors.New("media not found")

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
	GetMatchedMedia(page int, sortInfo SortInfo) (media []MatchedMedia, hasNext bool, err error)

	GetMatchedMediaBySeriesId(seriesId int64) (media MatchedMedia, err error)

	DeleteMedia(mediaType MediaType, id int64) error

	DeleteSeason(id int64, season int) error
}
