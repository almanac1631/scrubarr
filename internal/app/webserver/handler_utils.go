package webserver

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/almanac1631/scrubarr/internal/utils"
)

// htmxOnly is a middleware that returns 404 for non-HTMX requests.
// Apply it at the router level for endpoints that are only reachable via HTMX.
func htmxOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if !utils.IsHTMXRequest(request) {
			http.Error(writer, "404 Not Found", http.StatusNotFound)
			return
		}
		next(writer, request)
	}
}

func getSortInfoFromUrlQuery(values url.Values) SortInfo {
	sortInfo := SortInfo{}
	sortKeyRaw := values.Get("sortKey")
	switch SortKey(sortKeyRaw) {
	case SortKeyName, SortKeySize, SortKeyAdded, SortKeyStatus:
		sortInfo.Key = SortKey(sortKeyRaw)
	default:
		sortInfo.Key = SortKeyName
	}
	sortOrderRaw := values.Get("sortOrder")
	switch SortOrder(sortOrderRaw) {
	case SortOrderAsc, SortOrderDesc:
		sortInfo.Order = SortOrder(sortOrderRaw)
	default:
		sortInfo.Order = SortOrderAsc
	}
	return sortInfo
}

func getMediaFilterInfoFromUrlQuery(values url.Values) MediaFilterInfo {
	filter := MediaFilterInfo{
		Name:   strings.TrimSpace(values.Get("name")),
		OkOnly: values.Get("okOnly") == "true",
	}
	minSizeRaw := strings.TrimSpace(values.Get("minSize"))
	if minSize, err := strconv.ParseFloat(minSizeRaw, 64); err == nil && minSize > 0 && !math.IsInf(minSize, 0) {
		filter.MinSize = minSize
	}
	return filter
}
