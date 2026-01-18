package common

import "errors"

var ErrTrackerNotFound = errors.New("tracker not found")

type TrackerManager interface {
	GetTrackerName(trackers []string) (string, error)
}
