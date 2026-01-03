package common

import "time"

type TorrentClientFinding struct {
	Added time.Time
}

type TorrentFile struct {
	Path string
}

type TorrentEntry struct {
	Client string
	Id     string
	Name   string
	Added  time.Time
	Files  []*TorrentFile
}

type TorrentClientManager interface {
	SearchForMedia(originalFilePath string) (finding *TorrentClientFinding, err error)
}

type TorrentClientRetriever interface {
	GetTorrentEntries() ([]*TorrentEntry, error)
	Name() string
}
