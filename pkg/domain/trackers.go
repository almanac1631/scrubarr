package domain

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

func (t Tracker) IsValid() bool {
	return t.Name != ""
}

type TrackerManager interface {
	GetTracker(trackers []string) (Tracker, error)
}
