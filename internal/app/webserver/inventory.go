package webserver

import "errors"

var ErrMalformedMediaId = errors.New("malformed media id")
var ErrMediaNotFound = errors.New("media not found")

type SortKey string

const (
	SortKeyName   SortKey = "name"
	SortKeySize   SortKey = "size"
	SortKeyAdded  SortKey = "added"
	SortKeyStatus SortKey = "status"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type SortInfo struct {
	Key   SortKey
	Order SortOrder
}

type InventoryService interface {
	GetMediaInventory(page int, sortInfo SortInfo) (mediaRows []MediaRow, hasNext bool, err error)

	GetExpandedMediaRow(id string) (mediaRow MediaRow, err error)

	DeleteMedia(id string) error
}
