package common

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
	Tracker        Tracker
	TorrentFinding *TorrentEntry
}

type MatchedMedia struct {
	MediaMetadata
	Parts []MatchedMediaPart
}
