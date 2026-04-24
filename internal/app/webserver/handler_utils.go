package webserver

import (
	"net/http"
	"net/url"

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
