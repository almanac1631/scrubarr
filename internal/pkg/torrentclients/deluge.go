package torrentclients

import (
	"fmt"

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

func (retriever *DelugeRetriever) SearchForMovie(originalFilePath string) (bool, error) {
	if retriever.torrentListCache == nil {
		var err error
		retriever.torrentListCache, err = retriever.client.TorrentsStatus(delugeclient.StateSeeding, []string{})
		if err != nil {
			return false, fmt.Errorf("could not get torrent list from deluge: %w", err)
		}
	}
	for _, torrent := range retriever.torrentListCache {
		if len(torrent.Files) == 0 {
			continue
		}
		fileNameCmp := torrent.Files[0].Path
		if fileNameCmp == originalFilePath {
			return true, nil
		}
	}
	return false, nil
}
