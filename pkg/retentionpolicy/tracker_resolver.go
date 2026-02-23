package retentionpolicy

import (
	"errors"

	"github.com/almanac1631/scrubarr/pkg/domain"
)

var ErrTrackerNotFound = errors.New("tracker not found")

type TrackerResolver interface {
	Resolve(trackers []string) (*domain.Tracker, error)
}
