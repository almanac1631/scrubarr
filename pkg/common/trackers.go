package common

import (
	"errors"
	"time"
)

var ErrTrackerNotFound = errors.New("tracker not found")

type Tracker struct {
	Name     string
	MinRatio float64
	MinAge   time.Duration
}

type TrackerManager interface {
	GetTracker(trackers []string) (Tracker, error)
}
