package common

import (
	"time"
)

type MatchedEntry struct {
	MediaMetadata
	Size  int64
	Parts []MatchedEntryPart
}

type MatchedEntryPart struct {
	MediaPart          MediaPart
	TorrentInformation TorrentInformation
}

type TorrentStatus string

const (
	TorrentStatusMissing    TorrentStatus = "missing"
	TorrentStatusPresent    TorrentStatus = "present"
	TorrentStatusIncomplete TorrentStatus = "incomplete"
)

type TorrentAttributeStatus string

const (
	TorrentAttributeStatusFulfilled TorrentAttributeStatus = "fulfilled"
	TorrentAttributeStatusPending   TorrentAttributeStatus = "pending"
	TorrentAttributeStatusUnknown   TorrentAttributeStatus = "unknown"
)

type TorrentInformation struct {
	Client      string
	Id          string
	Status      TorrentStatus
	Tracker     Tracker
	RatioStatus TorrentAttributeStatus
	Ratio       float64
	AgeStatus   TorrentAttributeStatus
	Age         time.Duration
}

func (torrentInformation TorrentInformation) GetScore() int {
	score := 100
	if torrentInformation.Status == TorrentStatusPresent {
		score -= 60
	}
	getStatusScore := func(status TorrentAttributeStatus) int {
		if status == TorrentAttributeStatusFulfilled {
			return 20
		} else if status == TorrentAttributeStatusPending {
			return 10
		}
		return 0
	}
	score -= getStatusScore(torrentInformation.RatioStatus)
	score -= getStatusScore(torrentInformation.AgeStatus)
	return score
}
