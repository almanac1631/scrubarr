package domain

import (
	"errors"
	"fmt"
	"time"
)

var ErrTorrentNotFound = errors.New("torrent not found")

type TorrentFile struct {
	Path string
	Size int64
}

type TorrentEntry struct {
	Client   string
	Id       string
	Name     string
	Ratio    float64
	Added    time.Time
	Trackers []string
	Files    []*TorrentFile
}

func (t TorrentEntry) String() string {
	return fmt.Sprintf("{Client:%s Id:%s Name:%s Trackers: %+v}", t.Client, t.Id, t.Name, t.Trackers)
}

type TorrentSourceManager interface {
	CachedManager
	GetTorrents() ([]*TorrentEntry, error)
	DeleteTorrent(client string, id string) error
}

type TorrentSource interface {
	GetTorrentEntries() ([]*TorrentEntry, error)
	DeleteTorrent(id string) error
	Name() string
}
