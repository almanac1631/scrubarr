package domain

import "time"

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

type MediaFile struct {
	Id               int64
	Season           int
	OriginalFilePath string
	Size             int64
}

type MediaEntry struct {
	MediaMetadata
	Files []MediaFile
}

type MediaSourceManager interface {
	CachedManager
	GetMedia() ([]*MediaEntry, error)
	DeleteMediaFiles(mediaType MediaType, fileIds []int64, stopParentMonitoring bool) error
}

type MediaSource interface {
	GetMedia() ([]MediaEntry, error)
	SupportedMediaType() MediaType
	DeleteMediaFiles(fileIds []int64, stopParentMonitoring bool) error
}
