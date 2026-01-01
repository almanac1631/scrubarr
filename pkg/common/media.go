package common

import (
	"time"
)

type Media struct {
	Title            string
	Size             int64
	Added            time.Time
	OriginalFilePath string
	Url              string
}

type MediaInfo struct {
	Media
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
	GetMediaInfos(page int, sortInfo SortInfo) (medias []MediaInfo, hasNext bool, err error)
}
