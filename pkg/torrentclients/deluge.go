package torrentclients

import (
	"fmt"
	"time"

	"github.com/almanac1631/scrubarr/pkg/common"
	delugeclient "github.com/gdm85/go-libdeluge"
)

type DelugeRetriever struct {
	client *delugeclient.ClientV2
}

func NewDelugeRetriever(hostname string, port uint, username string, password string) (*DelugeRetriever, error) {
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
	return &DelugeRetriever{client}, nil
}

func (retriever *DelugeRetriever) GetTorrentEntries() ([]*common.TorrentEntry, error) {
	torrentList, err := retriever.client.TorrentsStatus(delugeclient.StateSeeding, []string{})
	if err != nil {
		return nil, fmt.Errorf("could not get torrent list from deluge rpc api: %w", err)
	}
	torrentEntries := make([]*common.TorrentEntry, 0, len(torrentList))
	for hash, torrent := range torrentList {
		torrentEntry := &common.TorrentEntry{
			Client: retriever.Name(),
			Id:     hash,
			Name:   torrent.Name,
			Added:  time.Unix(int64(torrent.TimeAdded), 0),
			Files:  []*common.TorrentFile{},
		}
		for _, file := range torrent.Files {
			torrentEntry.Files = append(torrentEntry.Files, &common.TorrentFile{
				Path: file.Path,
			})
		}
		torrentEntries = append(torrentEntries, torrentEntry)
	}
	return torrentEntries, nil
}

func (retriever *DelugeRetriever) Name() string {
	return "deluge"
}
