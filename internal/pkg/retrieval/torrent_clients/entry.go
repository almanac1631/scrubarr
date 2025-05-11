package torrent_clients

import (
	"fmt"
	"time"
)

type TorrentClientEntry struct {
	ID                any
	TorrentClientName string
	TorrentName       string
	DownloadFilePath  string
	DownloadedAt      time.Time
	Ratio             float32
	FileSizeBytes     int64
	TrackerHost       string
}

func (entry TorrentClientEntry) String() string {
	return fmt.Sprintf("TorrentName=%s, DownloadFilePath=%s", entry.TorrentName, entry.DownloadFilePath)
}

type TorrentEntryRetriever interface {
	RetrieveTorrentEntries() (map[string]TorrentClientEntry, error)
}
