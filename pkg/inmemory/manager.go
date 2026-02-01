package inmemory

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/almanac1631/scrubarr/pkg/domain"
)

var _ domain.Manager = (*Manager)(nil)

type Manager struct {
	matchedMediasCache []domain.MatchedMedia

	mediaManager domain.MediaManager

	torrentManager domain.TorrentClientManager

	trackerManager domain.TrackerManager
}

func NewManager(mediaManager domain.MediaManager, torrentManager domain.TorrentClientManager, trackerManager domain.TrackerManager) *Manager {
	return &Manager{
		nil, mediaManager, torrentManager, trackerManager,
	}
}

const pageSize = 10

func (m *Manager) refreshCache() error {
	m.matchedMediasCache = make([]domain.MatchedMedia, 0)
	medias, err := m.mediaManager.GetMedia()
	if err != nil {
		return fmt.Errorf("could not retrieve media from media manager: %w", err)
	}

	for _, mediaEntry := range medias {
		size := int64(0)
		parts := make([]domain.MatchedMediaPart, 0, len(mediaEntry.MediaParts))
		added := mediaEntry.Added
		for _, part := range mediaEntry.MediaParts {
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

			torrentStatus := domain.TorrentStatusMissing
			var torrentClient, torrentId string
			if torrentFinding != nil {
				torrentStatus = domain.TorrentStatusPresent
				torrentClient = torrentFinding.Client
				torrentId = torrentFinding.Id
			}

			ratioStatus, ratio := m.getRatioStatus(torrentFinding, tracker)
			ageStatus, age := m.getAgeStatus(torrentFinding, tracker)

			mediaPart := domain.MatchedMediaPart{
				MediaPart: part,
				TorrentInformation: domain.TorrentInformation{
					Client:      torrentClient,
					Id:          torrentId,
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
		matchedMedias := domain.MatchedMedia{
			MediaMetadata: mediaEntry.MediaMetadata,
			Parts:         parts,
			Size:          size,
		}
		matchedMedias.Added = added
		m.matchedMediasCache = append(m.matchedMediasCache, matchedMedias)
	}
	return nil
}

func (m *Manager) getTracker(torrentFinding *domain.TorrentEntry, mediaEntry domain.MediaEntry) (domain.Tracker, error) {
	if torrentFinding != nil {
		tracker, err := m.trackerManager.GetTracker(torrentFinding.Trackers)
		if err != nil {
			if errors.Is(err, domain.ErrTrackerNotFound) {
				slog.Warn("Could not find tracker name for media entry.",
					"mediaType", mediaEntry.Type, "mediaId", mediaEntry.Id, "part", torrentFinding.Name,
					"trackers", torrentFinding.Trackers, "findingId", torrentFinding.Id, "findingClient", torrentFinding.Client)
				return domain.Tracker{}, nil
			}
			return domain.Tracker{}, err
		}
		return tracker, nil
	}
	return domain.Tracker{}, nil
}

func (m *Manager) getRatioStatus(torrentFinding *domain.TorrentEntry, tracker domain.Tracker) (ratioStatus domain.TorrentAttributeStatus, ratio float64) {
	ratioStatus = domain.TorrentAttributeStatusUnknown
	ratio = -1.0
	if tracker.IsValid() && torrentFinding != nil {
		ratio = torrentFinding.Ratio
		if torrentFinding.Ratio >= tracker.MinRatio {
			ratioStatus = domain.TorrentAttributeStatusFulfilled
		} else {
			ratioStatus = domain.TorrentAttributeStatusPending
		}
	}
	return ratioStatus, ratio
}

func (m *Manager) getAgeStatus(torrentFinding *domain.TorrentEntry, tracker domain.Tracker) (ageStatus domain.TorrentAttributeStatus, age time.Duration) {
	ageStatus = domain.TorrentAttributeStatusUnknown
	age = time.Duration(-1)
	if tracker.IsValid() && torrentFinding != nil {
		age = time.Since(torrentFinding.Added)
		if time.Since(torrentFinding.Added) > tracker.MinAge {
			ageStatus = domain.TorrentAttributeStatusFulfilled
		} else {
			ageStatus = domain.TorrentAttributeStatusPending
		}
	}
	return ageStatus, age
}
