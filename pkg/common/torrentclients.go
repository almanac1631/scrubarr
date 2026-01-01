package common

import "time"

type TorrentClientFinding struct {
	AddedOn time.Time
}

type TorrentClientManager interface {
	SearchForMovie(originalFilePath string) (finding *TorrentClientFinding, err error)
}
