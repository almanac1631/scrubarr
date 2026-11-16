package torrentclients

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/almanac1631/scrubarr/pkg/domain"
	delugeclient "github.com/gdm85/go-libdeluge"
)

var _ domain.TorrentSource = (*DelugeRetriever)(nil)

type DelugeRetriever struct {
	client       *delugeclient.ClientV2
	dryRun       bool
	singleSystem bool
}

func NewDelugeRetriever(hostname string, port uint, username string, password string, dryRun bool, singleSystem bool) (*DelugeRetriever, error) {
	client := delugeclient.NewV2(delugeclient.Settings{
		Hostname: hostname,
		Port:     port,
		Login:    username,
		Password: password,
	})
	err := client.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect to remote deluge rpc api: %w", err)
	}
	return &DelugeRetriever{client, dryRun, singleSystem}, nil
}

func (retriever *DelugeRetriever) GetTorrentEntries() ([]*domain.TorrentEntry, error) {
	torrentList, err := retriever.client.TorrentsStatus(delugeclient.StateUnspecified, []string{})
	if err != nil {
		return nil, fmt.Errorf("could not get torrent list from deluge rpc api: %w", err)
	}
	labels, err := retriever.getLabels()
	if err != nil {
		return nil, err
	}
	torrentEntries := make([]*domain.TorrentEntry, 0, len(torrentList))
	for hash, torrent := range torrentList {
		// only seeding torrents or finished ones that are paused
		if torrent.State != string(delugeclient.StateSeeding) &&
			(torrent.State != string(delugeclient.StatePaused) || torrent.Progress < 100) {
			continue
		}
		torrentEntry := &domain.TorrentEntry{
			Client:   retriever.Name(),
			Id:       hash,
			Name:     torrent.Name,
			Label:    labels[hash],
			Added:    time.Unix(torrent.CompletedTime, 0).In(time.UTC),
			Files:    []*domain.TorrentFile{},
			Trackers: []string{torrent.TrackerHost},
			Ratio:    float64(torrent.Ratio),
		}
		for _, file := range torrent.Files {
			torrentFile := &domain.TorrentFile{
				Path: file.Path,
				Size: file.Size,
			}
			if retriever.singleSystem {
				torrentFile.AbsolutePath = joinRemotePath(torrent.SavePath, file.Path)
			}
			torrentEntry.Files = append(torrentEntry.Files, torrentFile)
		}
		torrentEntries = append(torrentEntries, torrentEntry)
	}
	return torrentEntries, nil
}

// getLabels returns a map of torrent hash to label. If the Label plugin is not enabled on the deluge server, an empty
// map is returned and no labels are set.
func (retriever *DelugeRetriever) getLabels() (map[string]string, error) {
	labelPlugin, err := retriever.client.LabelPlugin()
	if err != nil {
		return nil, fmt.Errorf("could not check enabled plugins on deluge rpc api: %w", err)
	}
	if labelPlugin == nil {
		return map[string]string{}, nil
	}
	labels, err := labelPlugin.GetTorrentsLabels(delugeclient.StateUnspecified, []string{})
	if err != nil {
		return nil, fmt.Errorf("could not get torrent labels from deluge rpc api: %w", err)
	}
	return labels, nil
}

func (retriever *DelugeRetriever) DeleteTorrent(id string) error {
	if retriever.dryRun {
		slog.Info("[DRY RUN] Skipping deluge torrent deletion.", "id", id)
		return nil
	}
	ok, err := retriever.client.RemoveTorrent(id, true)
	if err != nil {
		var wrappedErr delugeclient.RPCError
		if errors.As(err, &wrappedErr) && wrappedErr.ExceptionType == "InvalidTorrentError" {
			return domain.ErrTorrentNotFound
		}
		return fmt.Errorf("could not remove torrent from deluge rpc api: %w", err)
	} else if !ok {
		return fmt.Errorf("could not remove torrent from deluge rpc api but no error was thrown")
	}
	return nil
}

func (retriever *DelugeRetriever) Name() string {
	return "deluge"
}
