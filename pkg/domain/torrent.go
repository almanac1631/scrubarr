package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrTorrentNotFound = errors.New("torrent not found")

type TorrentFile struct {
	Path         string
	AbsolutePath string
	Size         int64
}

type TorrentEntry struct {
	Client   string
	Id       string
	Name     string
	Label    string
	Ratio    float64
	Added    time.Time
	Trackers []string
	Files    []*TorrentFile
}

func (t TorrentEntry) String() string {
	return fmt.Sprintf("{Client:%s Id:%s Name:%s Trackers: %+v}", t.Client, t.Id, t.Name, t.Trackers)
}

func (t TorrentEntry) IsCrossSeed() bool {
	label := strings.ToLower(t.Label)
	return strings.Contains(label, "cross-seed") || strings.Contains(label, "crossseed")
}

func (t TorrentEntry) IsRetain() bool {
	return strings.HasPrefix(strings.ToLower(t.Label), "retain-")
}

type TorrentSourceManager interface {
	CachedManager
	GetTorrents() ([]*TorrentEntry, error)
	DeleteTorrent(client string, id string) error
	ClientNames() []string
}

type TorrentSource interface {
	GetTorrentEntries() ([]*TorrentEntry, error)
	DeleteTorrent(id string) error
	Name() string
}
