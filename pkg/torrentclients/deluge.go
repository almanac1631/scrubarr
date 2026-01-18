package torrentclients

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/almanac1631/scrubarr/pkg/common"
	delugeclient "github.com/gdm85/go-libdeluge"
)

var _ common.TorrentClientRetriever = (*DelugeRetriever)(nil)

type DelugeRetriever struct {
	client *delugeclient.ClientV2
	dryRun bool
}

func NewDelugeRetriever(hostname string, port uint, username string, password string, dryRun bool) (*DelugeRetriever, error) {
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
	return &DelugeRetriever{client, dryRun}, nil
}

func (retriever *DelugeRetriever) GetTorrentEntries() ([]*common.TorrentEntry, error) {
	torrentList, err := retriever.client.TorrentsStatus(delugeclient.StateSeeding, []string{})
	if err != nil {
		return nil, fmt.Errorf("could not get torrent list from deluge rpc api: %w", err)
	}
	torrentEntries := make([]*common.TorrentEntry, 0, len(torrentList))
	for hash, torrent := range torrentList {
		torrentEntry := &common.TorrentEntry{
			Client:   retriever.Name(),
			Id:       hash,
			Name:     torrent.Name,
			Added:    time.Unix(torrent.CompletedTime, 0).In(time.UTC),
			Files:    []*common.TorrentFile{},
			Trackers: []string{torrent.TrackerHost},
		}
		for _, file := range torrent.Files {
			torrentEntry.Files = append(torrentEntry.Files, &common.TorrentFile{
				Path: file.Path,
				Size: file.Size,
			})
		}
		torrentEntries = append(torrentEntries, torrentEntry)
	}
	return torrentEntries, nil
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
			return common.ErrTorrentNotFound
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
