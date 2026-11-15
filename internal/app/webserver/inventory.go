package webserver

import (
	"errors"
	"regexp"
	"strings"
)

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

type MediaFilterInfo struct {
	Name    NameFilter
	MinSize SizeGB
	OkOnly  bool
}

type OrphanedTorrentFilterInfo struct {
	Name    NameFilter
	MinSize SizeGB
	Client  string
	OkOnly  bool
}

// SizeGB is a size in gigabytes.
type SizeGB float64

func (s SizeGB) Bytes() int64 {
	if s <= 0 {
		return 0
	}
	return int64(float64(s) * 1000 * 1000 * 1000)
}

type NameFilter string

var (
	fuzzySpaceChars     = regexp.MustCompile(`[._\-‐‑‒–—―−]`)
	fuzzyRemovedChars   = regexp.MustCompile(`[()\[\]{}<>（）［］｛｝「」『』【】〔〕《》〈〉‹›«»:;!?'"]`)
	fuzzyRepeatedSpaces = regexp.MustCompile(`\s+`)
)

func (f NameFilter) FuzzyNameMatch(title string) bool {
	return strings.Contains(normalizeForFuzzyMatch(title), normalizeForFuzzyMatch(string(f)))
}

func normalizeForFuzzyMatch(s string) string {
	s = strings.ToLower(s)
	s = fuzzySpaceChars.ReplaceAllString(s, " ")
	s = fuzzyRemovedChars.ReplaceAllString(s, "")
	s = fuzzyRepeatedSpaces.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

type InventoryService interface {
	GetMediaInventory(page int, sortInfo SortInfo, filter MediaFilterInfo) (mediaRows []MediaRow, hasNext bool, err error)

	GetExpandedMediaRow(id string) (mediaRow MediaRow, err error)

	DeleteMedia(id string) error

	RefreshCache() error

	GetOrphanedTorrents(page int, sortInfo SortInfo, filter OrphanedTorrentFilterInfo) (rows []OrphanedTorrentRow, hasNext bool, err error)

	GetTorrentClientNames() []string

	DeleteOrphanedTorrent(id string) error
}
