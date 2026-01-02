package torrentclients

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/almanac1631/scrubarr/pkg/common"
	"github.com/autobrr/go-rtorrent"
)

type RtorrentRetriever struct {
	client               *rtorrent.Client
	torrentListCache     []rtorrent.Torrent
	torrentFileListCache map[string][]rtorrent.File
}

func NewRtorrentRetriever(hostname string, username string, password string) (*RtorrentRetriever, error) {
	client := rtorrent.NewClient(rtorrent.Config{
		Addr:      hostname,
		BasicUser: username,
		BasicPass: password,
	})
	_, err := client.Name(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not connect to remote rtorrent rpc api: %w", err)
	}
	return &RtorrentRetriever{client, nil, map[string][]rtorrent.File{}}, nil
}

func (r *RtorrentRetriever) RefreshCache() error {
	var err error
	r.torrentListCache, err = r.client.GetTorrents(context.Background(), rtorrent.ViewMain)
	if err != nil {
		return fmt.Errorf("could not get torrent list from rtorrent: %w", err)
	}
	for _, torrent := range r.torrentListCache {
		torrentFiles, err := r.client.GetFiles(context.Background(), torrent)
		if err != nil {
			return fmt.Errorf("could not get torrent files from rtorrent: %w", err)
		}
		r.torrentFileListCache[torrent.Hash] = torrentFiles
	}
	return nil
}

func (r *RtorrentRetriever) SearchForMedia(originalFilePath string) (finding *common.TorrentClientFinding, err error) {
	if r.torrentListCache == nil {
		if err := r.RefreshCache(); err != nil {
			return nil, err
		}
	}
	for _, torrent := range r.torrentListCache {
		torrentNameWithExt := torrent.Name + filepath.Ext(originalFilePath)
		if torrent.Name == originalFilePath || torrentNameWithExt == originalFilePath {
			return &common.TorrentClientFinding{
				Added: torrent.Finished,
			}, nil
		}
		torrentFiles, ok := r.torrentFileListCache[torrent.Hash]
		if !ok {
			return nil, fmt.Errorf("invalid state: rtorrent file list cache must contain torrent files for torrent: %q", torrent.Hash)
		}
		for _, file := range torrentFiles {
			if file.Path == originalFilePath {
				return &common.TorrentClientFinding{
					Added: torrent.Finished,
				}, nil
			}
		}
	}
	return nil, nil
}
