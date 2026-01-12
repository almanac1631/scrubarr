package common

import (
	"errors"
	"time"
)

var ErrMediaNotFound = errors.New("media not found")

type MediaType string

const (
	MediaTypeMovie  MediaType = "movie"
	MediaTypeSeries MediaType = "series"
)

type MediaMetadata struct {
	Id    int64
	Type  MediaType
	Title string
	Url   string
	Added time.Time
}

type MediaPart struct {
	Id               int64
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
	TorrentFinding *TorrentEntry
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
	GetMatchedMedia(page int, sortInfo SortInfo) (media []MatchedMedia, hasNext bool, err error)

	GetMatchedMediaBySeriesId(seriesId int64) (media MatchedMedia, err error)

	DeleteMedia(mediaType MediaType, id int64) error

	DeleteSeason(id int64, season int) error
}

type MediaRetriever interface {
	GetMedia() ([]Media, error)
	SupportedMediaType() MediaType
	DeleteMediaFiles(fileIds []int64, stopParentMonitoring bool) error
}
