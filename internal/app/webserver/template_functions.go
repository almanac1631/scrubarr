package webserver

import (
	"fmt"
	"html/template"
	"time"

	"github.com/almanac1631/scrubarr/internal/utils"
)

var templateFunctions = template.FuncMap{
	"formatBytes": utils.FormatBytes,
	"formatDate": func(t time.Time) string {
		// Go uses a reference date (Mon Jan 2 15:04:05 MST 2006) for layout
		return t.Format("2006-01-02")
	},
	"checkCurrentSort": func(sortKey SortKey, sortOrder SortOrder, currentSortInfo SortInfo) bool {
		return currentSortInfo.Key == sortKey && currentSortInfo.Order == sortOrder
	},
	"durationToNanoseconds": func(duration time.Duration) int64 {
		return duration.Nanoseconds()
	},
	"floatToStr": func(float float64) string {
		return fmt.Sprintf("%.2f", float)
	},
}
