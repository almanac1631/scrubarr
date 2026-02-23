package torrentclients

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/almanac1631/scrubarr/pkg/domain"
	"github.com/autobrr/go-rtorrent"
)

var _ domain.TorrentSource = (*RtorrentRetriever)(nil)

type RtorrentRetriever struct {
	client *rtorrent.Client
	dryRun bool
}

func NewRtorrentRetriever(hostname string, username string, password string, dryRun bool) (*RtorrentRetriever, error) {
	client := rtorrent.NewClient(rtorrent.Config{
		Addr:      hostname,
		BasicUser: username,
		BasicPass: password,
	})
	_, err := client.Name(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not connect to remote rtorrent rpc api: %w", err)
	}
	return &RtorrentRetriever{client, dryRun}, nil
}

func (retriever *RtorrentRetriever) GetTorrentEntries() ([]*domain.TorrentEntry, error) {
	torrentList, err := retriever.client.GetTorrents(context.Background(), rtorrent.ViewMain)
	if err != nil {
		return nil, fmt.Errorf("could not get torrent list from rtorrent: %w", err)
	}
	torrentEntries := make([]*domain.TorrentEntry, 0, len(torrentList))
	for _, torrent := range torrentList {
		torrentEntry := &domain.TorrentEntry{
			Client:   retriever.Name(),
			Id:       torrent.Hash,
			Name:     torrent.Name,
			Added:    torrent.Finished,
			Files:    []*domain.TorrentFile{},
			Trackers: []string{},
			Ratio:    torrent.Ratio,
		}
		torrentFiles, err := retriever.client.GetFiles(context.Background(), torrent)
		if err != nil {
			return nil, fmt.Errorf("could not get torrent files from rtorrent: %w", err)
		}
		for _, torrentFile := range torrentFiles {
			torrentEntry.Files = append(torrentEntry.Files, &domain.TorrentFile{
				Path: torrentFile.Path,
				Size: int64(torrentFile.Size),
			})
		}
		torrentEntries = append(torrentEntries, torrentEntry)
		torrentTrackers, err := retriever.client.GetTrackers(context.Background(), torrent)
		if err != nil {
			return nil, fmt.Errorf("could not get trackerresolver from rtorrent: %w", err)
		}
		for _, tracker := range torrentTrackers {
			if !slices.Contains(torrentEntry.Trackers, tracker) {
				torrentEntry.Trackers = append(torrentEntry.Trackers, tracker)
			}
		}
	}
	return torrentEntries, nil
}

func (retriever *RtorrentRetriever) DeleteTorrent(id string) error {
	hash := id
	if retriever.dryRun {
		slog.Info("[DRY RUN] Skipping rtorrent torrent deletion.", "hash", hash)
		return nil
	}
	torrent := rtorrent.Torrent{Hash: hash}
	if err := retriever.client.SetForceDelete(context.Background(), torrent, true); err != nil {
		errString := err.Error()
		if strings.Contains(errString, "Could not find info-hash") || strings.Contains(errString, "info-hash not found") {
			return domain.ErrTorrentNotFound
		}
		return fmt.Errorf("could not force deletion for torrent %q: %w", hash, err)
	}
	if err := retriever.client.DeleteTied(context.Background(), torrent); err != nil {
		return fmt.Errorf("could not delete tied files for torrent %q: %w", hash, err)
	}
	if err := retriever.client.Delete(context.Background(), torrent); err != nil {
		return fmt.Errorf("could not delete files for torrent %q: %w", hash, err)
	}
	return nil
}

func (retriever *RtorrentRetriever) Name() string {
	return "rtorrent"
}
