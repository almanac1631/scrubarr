package trackers

import (
	"fmt"
	"regexp"

	"github.com/almanac1631/scrubarr/pkg/common"
	"github.com/knadh/koanf/v2"
)

var _ common.TrackerManager = (*ConfigBasedTrackerManager)(nil)

type trackerConfig struct {
	Name    string
	Pattern *regexp.Regexp
}

type ConfigBasedTrackerManager struct {
	trackerConfigs []trackerConfig
}

func NewConfigBasedTrackerManager(config *koanf.Koanf) (*ConfigBasedTrackerManager, error) {
	manager := &ConfigBasedTrackerManager{
		trackerConfigs: make([]trackerConfig, 0),
	}
	for _, trackerKey := range config.MapKeys("trackers") {
		name := config.MustString(fmt.Sprintf("trackers.%s.name", trackerKey))
		patternRaw := config.MustString(fmt.Sprintf("trackers.%s.pattern", trackerKey))
		pattern, err := regexp.Compile(patternRaw)
		if err != nil {
			return nil, fmt.Errorf("could not compile pattern for tracker %q: %w", trackerKey, err)
		}
		manager.trackerConfigs = append(manager.trackerConfigs, trackerConfig{
			Name:    name,
			Pattern: pattern,
		})
	}
	return manager, nil
}

func (c ConfigBasedTrackerManager) GetTrackerName(trackers []string) (string, error) {
	for _, config := range c.trackerConfigs {
		for _, tracker := range trackers {
			if config.Pattern.MatchString(tracker) {
				return config.Name, nil
			}
		}
	}
	return "", common.ErrTrackerNotFound
}
