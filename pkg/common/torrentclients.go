package common

import "time"

type TorrentClientFinding struct {
	Added time.Time
}

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
	SearchForMedia(originalFilePath string, size int64) (finding *TorrentClientFinding, err error)
}

type TorrentClientRetriever interface {
	GetTorrentEntries() ([]*TorrentEntry, error)
	Name() string
}
