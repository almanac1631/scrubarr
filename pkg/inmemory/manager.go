package inmemory

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/almanac1631/scrubarr/pkg/common"
)

var _ common.Manager = (*Manager)(nil)

type Manager struct {
	matchedEntriesCache []common.MatchedEntry

	mediaManager common.MediaManager

	torrentManager common.TorrentClientManager

	trackerManager common.TrackerManager
}

func NewManager(mediaManager common.MediaManager, torrentManager common.TorrentClientManager, trackerManager common.TrackerManager) *Manager {
	return &Manager{
		nil, mediaManager, torrentManager, trackerManager,
	}
}

const pageSize = 10

func (m *Manager) refreshCache() error {
	m.matchedEntriesCache = make([]common.MatchedEntry, 0)
	medias, err := m.mediaManager.GetMedia()
	if err != nil {
		return fmt.Errorf("could not retrieve media from media manager: %w", err)
	}

	for _, mediaEntry := range medias {
		size := int64(0)
		parts := make([]common.MatchedEntryPart, 0, len(mediaEntry.Parts))
		added := mediaEntry.Added
		for _, part := range mediaEntry.Parts {
			size += part.Size
			torrentFinding, err := m.torrentManager.SearchForMedia(part.OriginalFilePath, part.Size)
			if err != nil {
				return err
			}
			// use latest available added date from torrent entry
			if torrentFinding != nil && !torrentFinding.Added.IsZero() && torrentFinding.Added.After(added) {
				added = torrentFinding.Added
			}

			tracker, err := m.getTracker(torrentFinding, mediaEntry)
			if err != nil {
				return err
			}

			torrentStatus := common.TorrentStatusMissing
			if torrentFinding != nil {
				torrentStatus = common.TorrentStatusPresent
			}

			ratioStatus, ratio := m.getRatioStatus(torrentFinding, tracker)
			ageStatus, age := m.getAgeStatus(torrentFinding, tracker)

			mediaPart := common.MatchedEntryPart{
				MediaPart: part,
				TorrentInformation: common.TorrentInformation{
					Status:      torrentStatus,
					Tracker:     tracker,
					RatioStatus: ratioStatus,
					AgeStatus:   ageStatus,
					Ratio:       ratio,
					Age:         age,
				},
			}
			parts = append(parts, mediaPart)
		}
		matchedEntries := common.MatchedEntry{
			MediaMetadata: mediaEntry.MediaMetadata,
			Parts:         parts,
			Size:          size,
		}
		matchedEntries.Added = added
		m.matchedEntriesCache = append(m.matchedEntriesCache, matchedEntries)
	}
	return nil
}

func (m *Manager) getTracker(torrentFinding *common.TorrentEntry, mediaEntry common.Media) (common.Tracker, error) {
	if torrentFinding != nil {
		tracker, err := m.trackerManager.GetTracker(torrentFinding.Trackers)
		if err != nil {
			if errors.Is(err, common.ErrTrackerNotFound) {
				slog.Warn("Could not find tracker name for media entry.",
					"mediaType", mediaEntry.Type, "mediaId", mediaEntry.Id, "part", torrentFinding.Name,
					"trackers", torrentFinding.Trackers, "findingId", torrentFinding.Id, "findingClient", torrentFinding.Client)
				return common.Tracker{}, nil
			}
			return common.Tracker{}, err
		}
		return tracker, nil
	}
	return common.Tracker{}, nil
}

func (m *Manager) getRatioStatus(torrentFinding *common.TorrentEntry, tracker common.Tracker) (ratioStatus common.TorrentAttributeStatus, ratio float64) {
	ratioStatus = common.TorrentAttributeStatusUnknown
	ratio = -1.0
	if tracker.IsValid() && torrentFinding != nil {
		ratio = torrentFinding.Ratio
		if torrentFinding.Ratio >= tracker.MinRatio {
			ratioStatus = common.TorrentAttributeStatusFulfilled
		} else {
			ratioStatus = common.TorrentAttributeStatusPending
		}
	}
	return ratioStatus, ratio
}

func (m *Manager) getAgeStatus(torrentFinding *common.TorrentEntry, tracker common.Tracker) (ageStatus common.TorrentAttributeStatus, age time.Duration) {
	ageStatus = common.TorrentAttributeStatusUnknown
	age = time.Duration(-1)
	if tracker.IsValid() && torrentFinding != nil {
		age = time.Since(torrentFinding.Added)
		if time.Since(torrentFinding.Added) > tracker.MinAge {
			ageStatus = common.TorrentAttributeStatusFulfilled
		} else {
			ageStatus = common.TorrentAttributeStatusPending
		}
	}
	return ageStatus, age
}
