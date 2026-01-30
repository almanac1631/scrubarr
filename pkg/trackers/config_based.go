package trackers

import (
	"fmt"
	"regexp"
	"time"

	"github.com/almanac1631/scrubarr/pkg/common"
	"github.com/knadh/koanf/v2"
)

var _ common.TrackerManager = (*ConfigBasedTrackerManager)(nil)

type trackerConfig struct {
	common.Tracker
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
		var minRatio float64
		minRatio, err := getSetConfigValue[float64](config, fmt.Sprintf("trackers.%s.min_ratio", trackerKey))
		if err != nil {
			return nil, err
		}
		minAge, err := getSetConfigValue[time.Duration](config, fmt.Sprintf("trackers.%s.min_age", trackerKey))
		if err != nil {
			return nil, err
		}
		patternRaw := config.MustString(fmt.Sprintf("trackers.%s.pattern", trackerKey))
		pattern, err := regexp.Compile(patternRaw)
		if err != nil {
			return nil, fmt.Errorf("could not compile pattern for tracker %q: %w", trackerKey, err)
		}
		manager.trackerConfigs = append(manager.trackerConfigs, trackerConfig{
			Tracker: common.Tracker{
				Name:     name,
				MinRatio: minRatio,
				MinAge:   minAge,
			},
			Pattern: pattern,
		})
	}
	return manager, nil
}

func getSetConfigValue[V float64 | time.Duration](config *koanf.Koanf, key string) (V, error) {
	if !config.Exists(key) {
		return 0, fmt.Errorf("no value for key %q found", key)
	}
	var zero V
	switch any(zero).(type) {
	case float64:
		return V(config.Float64(key)), nil
	case time.Duration:
		durationStr := config.MustString(key)
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			return zero, fmt.Errorf("could not parse duration %q: %w", durationStr, err)
		}
		return V(duration), nil
	default:
		return zero, fmt.Errorf("invalid type for key %q found", key)
	}
}

func (c ConfigBasedTrackerManager) GetTracker(trackers []string) (common.Tracker, error) {
	for _, config := range c.trackerConfigs {
		for _, tracker := range trackers {
			if config.Pattern.MatchString(tracker) {
				return config.Tracker, nil
			}
		}
	}
	return common.Tracker{}, common.ErrTrackerNotFound
}
