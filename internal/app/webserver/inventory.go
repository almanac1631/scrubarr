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
	Name    string
	MinSize float64 // in GB
	OkOnly  bool
}

var (
	fuzzySpaceChars     = regexp.MustCompile(`[._\-‐‑‒–—―−]`)
	fuzzyRemovedChars   = regexp.MustCompile(`[()\[\]{}<>（）［］｛｝「」『』【】〔〕《》〈〉‹›«»:;!?'"]`)
	fuzzyRepeatedSpaces = regexp.MustCompile(`\s+`)
)

// Both sides are normalized, so e.g. "2.Test", "2_Test" and "2 Test" match each other.
func (f MediaFilterInfo) FuzzyNameMatch(title string) bool {
	query := normalizeForFuzzyMatch(f.Name)
	return strings.Contains(normalizeForFuzzyMatch(title), query)
}

func normalizeForFuzzyMatch(s string) string {
	s = strings.ToLower(s)
	s = fuzzySpaceChars.ReplaceAllString(s, " ")
	s = fuzzyRemovedChars.ReplaceAllString(s, "")
	s = fuzzyRepeatedSpaces.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func (f MediaFilterInfo) MinSizeBytes() int64 {
	if f.MinSize <= 0 {
		return 0
	}
	return int64(f.MinSize * 1000 * 1000 * 1000)
}

type InventoryService interface {
	GetMediaInventory(page int, sortInfo SortInfo, filter MediaFilterInfo) (mediaRows []MediaRow, hasNext bool, err error)

	GetExpandedMediaRow(id string) (mediaRow MediaRow, err error)

	DeleteMedia(id string) error

	RefreshCache() error

	GetOrphanedTorrents(page int, sortInfo SortInfo) (rows []OrphanedTorrentRow, hasNext bool, err error)

	DeleteOrphanedTorrent(id string) error
}
