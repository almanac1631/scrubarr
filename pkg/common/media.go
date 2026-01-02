package common

import (
	"time"
)

type MediaType string

const (
	MediaTypeMovie  MediaType = "movie"
	MediaTypeSeries MediaType = "series"
)

type MediaMetadata struct {
	Type  MediaType
	Title string
	Url   string
	Added time.Time
}

type MediaPart struct {
	Season           int
	OriginalFilePath string
	Size             int64
}

type Media struct {
	MediaMetadata
	Parts []MediaPart
}

type MatchedMediaPart struct {
	MediaPart
	ExistsInTorrentClient bool
}

type MatchedMedia struct {
	MediaMetadata
	Parts []MatchedMediaPart
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
	GetMatchedMedia(page int, sortInfo SortInfo) (medias []MatchedMedia, hasNext bool, err error)
}
