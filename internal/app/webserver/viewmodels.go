package webserver

import (
	"fmt"
	"time"

	"github.com/almanac1631/scrubarr/pkg/domain"
)

type basePageData struct {
	SortInfo  SortInfo
	NextPage  int
	Version   string
	DiskQuota DiskQuota
	PageTitle string
}

type TorrentInformation struct {
	LinkStatus TorrentLinkStatus

	Tracker domain.Tracker

	Ratio float64
	Age   time.Duration
}

type MediaRow struct {
	Id    string
	Type  domain.MediaType
	Title string
	Url   string
	Size  int64
	Added time.Time

	TorrentInformation TorrentInformation
	Decision           domain.Decision

	AllowDeletion bool

	ChildMediaRows []MediaRow
}

func (m MediaRow) String() string {
	return fmt.Sprintf("id=%s", m.Id)
}

type TorrentLinkStatus string

const (
	TorrentLinkPresent    TorrentLinkStatus = "present"
	TorrentLinkMissing    TorrentLinkStatus = "missing"
	TorrentLinkIncomplete TorrentLinkStatus = "incomplete"
)

type OrphanedTorrentRow struct {
	Id            string
	Name          string
	Client        string
	Ratio         float64
	Added         time.Time
	Age           time.Duration
	Size          int64
	Decision      domain.Decision
	Tracker       domain.Tracker
	AllowDeletion bool
}
