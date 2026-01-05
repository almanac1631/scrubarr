package common

import (
	"errors"
	"time"
)

var ErrTorrentNotFound = errors.New("torrent not found")

type TorrentFile struct {
	Path string
	Size int64
}

type TorrentEntry struct {
	Client string
	Id     string
	Name   string
	Added  time.Time
	Files  []*TorrentFile
}

type TorrentClientManager interface {
	SearchForMedia(originalFilePath string, size int64) (finding *TorrentEntry, err error)
	DeleteFinding(client string, id string) error
}

type TorrentClientRetriever interface {
	GetTorrentEntries() ([]*TorrentEntry, error)
	DeleteTorrent(id string) error
	Name() string
}
