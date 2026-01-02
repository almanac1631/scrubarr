package torrentclients

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/almanac1631/scrubarr/pkg/common"
	delugeclient "github.com/gdm85/go-libdeluge"
)

type DelugeRetriever struct {
	client           *delugeclient.ClientV2
	torrentListCache map[string]*delugeclient.TorrentStatus
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
	return &DelugeRetriever{client, nil}, nil
}

func (retriever *DelugeRetriever) RefreshCache() error {
	var err error
	retriever.torrentListCache, err = retriever.client.TorrentsStatus(delugeclient.StateSeeding, []string{})
	if err != nil {
		return fmt.Errorf("could not get torrent list from deluge: %w", err)
	}
	return nil
}

func (retriever *DelugeRetriever) SearchForMedia(originalFilePath string) (finding *common.TorrentClientFinding, err error) {
	if retriever.torrentListCache == nil {
		if err := retriever.RefreshCache(); err != nil {
			return nil, err
		}
	}
	for _, torrent := range retriever.torrentListCache {
		torrentNameWithExt := torrent.Name + filepath.Ext(originalFilePath)
		if torrent.Name == originalFilePath || torrentNameWithExt == originalFilePath {
			return &common.TorrentClientFinding{
				Added: time.Unix(int64(torrent.TimeAdded), 0),
			}, nil
		}
		if len(torrent.Files) == 0 {
			continue
		}
		for _, file := range torrent.Files {
			fileNameCmp := filepath.Base(file.Path)
			if fileNameCmp == originalFilePath {
				return &common.TorrentClientFinding{
					Added: time.Unix(int64(torrent.TimeAdded), 0),
				}, nil
			}
		}
	}
	return nil, nil
}
