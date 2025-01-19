package arr_apps

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
}
