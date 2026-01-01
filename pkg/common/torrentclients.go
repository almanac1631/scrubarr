package common

import "time"

type TorrentClientFinding struct {
	Added time.Time
}

type TorrentClientManager interface {
	SearchForMovie(originalFilePath string) (finding *TorrentClientFinding, err error)
}
