package arr_apps

import "time"

type MediaType string

const (
	MediaTypeMovie  MediaType = "movie"
	MediaTypeSeries MediaType = "series"
)

type ArrAppEntry struct {
	Type          MediaType
	ParentName    string
	Monitored     bool
	MediaFilePath string
	DateAdded     time.Time
}
