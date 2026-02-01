package webserver

import (
	"html/template"
	"time"

	"github.com/almanac1631/scrubarr/internal/utils"
	"github.com/almanac1631/scrubarr/pkg/domain"
)

var templateFunctions = template.FuncMap{
	"formatBytes": utils.FormatBytes,
	"formatDate": func(t time.Time) string {
		// Go uses a reference date (Mon Jan 2 15:04:05 MST 2006) for layout
		return t.Format("2006-01-02")
	},
	"checkCurrentSort": func(sortKey domain.SortKey, sortOrder domain.SortOrder, currentSortInfo domain.SortInfo) bool {
		return currentSortInfo.Key == sortKey && currentSortInfo.Order == sortOrder
	},
	"getSeasonGroupedParts": getSeasonGroupedParts,
	"durationToNanoseconds": func(duration time.Duration) int64 {
		return duration.Nanoseconds()
	},
}
