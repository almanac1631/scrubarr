package torrent_clients

import (
	"context"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval"
	"github.com/autobrr/go-rtorrent"
	"maps"
	"path"
	"slices"
)

type RtorrentEntryRetriever struct {
	client             *rtorrent.Client
	allowedFileEndings []string
}

func NewRtorrentEntryRetriever(allowedFileEndings []string, hostname string, username string, password string) (*RtorrentEntryRetriever, error) {
	client := rtorrent.NewClient(rtorrent.Config{
		Addr:      hostname,
		BasicUser: username,
		BasicPass: password,
	})
	_, err := client.Name(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not connect to remote rtorrent rpc api: %w", err)
	}
	return &RtorrentEntryRetriever{client, allowedFileEndings}, nil
}

func (r *RtorrentEntryRetriever) RetrieveEntries() (map[retrieval.EntryName]retrieval.Entry, error) {
	torrentList, err := r.client.GetTorrents(context.Background(), rtorrent.ViewMain)
	if err != nil {
		return nil, fmt.Errorf("could not get torrent list from rtorrent: %w", err)
	}
	mediaEntryList := map[retrieval.EntryName]retrieval.Entry{}
	for _, torrent := range torrentList {
		torrentFileList, err := r.client.GetFiles(context.Background(), torrent)
		if err != nil {
			return nil, fmt.Errorf("could not get torrent files for torrent %q: %w", torrent.Name, err)
		}
		maps.Copy(mediaEntryList, r.parseTorrentFileList(torrent, torrentFileList))
	}
	return mediaEntryList, nil
}

func (r *RtorrentEntryRetriever) parseTorrentFileList(torrent rtorrent.Torrent, torrentFileList []rtorrent.File) map[retrieval.EntryName]retrieval.Entry {
	mediaEntryList := map[retrieval.EntryName]retrieval.Entry{}
	for _, torrentFile := range torrentFileList {
		name := retrieval.EntryName(path.Base(torrentFile.Path))
		fileExtension := path.Ext(torrentFile.Path)
		if !slices.Contains(r.allowedFileEndings, fileExtension) {
			continue
		}
		filePath := path.Join(torrent.Path, torrentFile.Path)
		entry := retrieval.Entry{
			Name: name,
			AdditionalData: TorrentClientEntry{
				TorrentClientName: "rtorrent",
				TorrentName:       torrent.Name,
				DownloadFilePath:  filePath,
				DownloadedAt:      torrent.Finished,
				Ratio:             float32(torrent.Ratio),
				FileSizeBytes:     int64(torrentFile.Size),
				TrackerHost:       "<unknown>",
			},
		}
		mediaEntryList[entry.Name] = entry
	}
	return mediaEntryList
}
