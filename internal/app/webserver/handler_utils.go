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
		sortInfo.Key = SortKeyAdded
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
	return MediaFilterInfo{
		Name:    getNameFilterFromUrlQuery(values),
		MinSize: getMinSizeFromUrlQuery(values),
		OkOnly:  values.Get("okOnly") == "true",
	}
}

func getOrphanedTorrentFilterInfoFromUrlQuery(values url.Values) OrphanedTorrentFilterInfo {
	return OrphanedTorrentFilterInfo{
		Name:             getNameFilterFromUrlQuery(values),
		MinSize:          getMinSizeFromUrlQuery(values),
		Client:           values.Get("client"),
		OkOnly:           values.Get("okOnly") == "true",
		ExcludeCrossSeed: values.Get("excludeCrossSeed") == "true",
	}
}

func getNameFilterFromUrlQuery(values url.Values) NameFilter {
	return NameFilter(strings.TrimSpace(values.Get("name")))
}

// getMinSizeFromUrlQuery parses minSize in GB. Invalid values return 0.
func getMinSizeFromUrlQuery(values url.Values) SizeGB {
	minSize, err := strconv.ParseFloat(strings.TrimSpace(values.Get("minSize")), 64)
	if err != nil || minSize <= 0 || math.IsInf(minSize, 0) {
		return 0
	}
	return SizeGB(minSize)
}
