package webserver

import "net/url"

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
