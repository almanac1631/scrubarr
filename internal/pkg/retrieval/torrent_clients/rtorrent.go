package torrent_clients

import (
	"context"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/autobrr/go-rtorrent"
	"maps"
	"path"
	"slices"
)

var _ common.EntryRetriever = (*RtorrentEntryRetriever)(nil)

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

func (r *RtorrentEntryRetriever) RetrieveEntries() (common.RetrieverEntries, error) {
	torrentList, err := r.client.GetTorrents(context.Background(), rtorrent.ViewMain)
	if err != nil {
		return nil, fmt.Errorf("could not get torrent list from rtorrent: %w", err)
	}
	mediaEntryList := common.RetrieverEntries{}
	for _, torrent := range torrentList {
		torrentFileList, err := r.client.GetFiles(context.Background(), torrent)
		if err != nil {
			return nil, fmt.Errorf("could not get torrent files for torrent %q: %w", torrent.Name, err)
		}
		maps.Copy(mediaEntryList, r.parseTorrentFileList(torrent, torrentFileList))
	}
	return mediaEntryList, nil
}

func (r *RtorrentEntryRetriever) parseTorrentFileList(torrent rtorrent.Torrent, torrentFileList []rtorrent.File) map[common.EntryName]common.Entry {
	mediaEntryList := map[common.EntryName]common.Entry{}
	for _, torrentFile := range torrentFileList {
		name := common.EntryName(path.Base(torrentFile.Path))
		fileExtension := path.Ext(torrentFile.Path)
		if !slices.Contains(r.allowedFileEndings, fileExtension) {
			continue
		}
		filePath := path.Join(torrent.Path, torrentFile.Path)
		entry := common.Entry{
			Name: name,
			AdditionalData: TorrentClientEntry{
				ID:                torrent.Hash,
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
	if len(mediaEntryList) == 1 {
		for _, entry := range mediaEntryList {
			delete(mediaEntryList, entry.Name)
			fileExtension := path.Ext(string(entry.Name))
			combinedTorrentName := common.EntryName(torrent.Name)
			torrentNameFileExt := getTorrentNameFileExt(torrent.Name, fileExtension)
			if torrentNameFileExt == fileExtension {
				combinedTorrentName = common.EntryName(torrent.Name)
			} else {
				combinedTorrentName = common.EntryName(fmt.Sprintf("%s%s", torrent.Name, fileExtension))
			}
			entry.Name = combinedTorrentName
			mediaEntryList[combinedTorrentName] = entry
		}
	}
	return mediaEntryList
}

func getTorrentNameFileExt(torrentName string, expectedFileExtension string) string {
	if path.Ext(torrentName) == "" {
		return ""
	} else if len(torrentName) < len(expectedFileExtension) {
		return ""
	}
	return torrentName[len(torrentName)-len(expectedFileExtension):]
}

func (r *RtorrentEntryRetriever) DeleteEntry(id any) error {
	hash, ok := id.(string)
	if !ok {
		return fmt.Errorf("could not convert id to hash string: %q", id)
	}
	torrent := rtorrent.Torrent{Hash: hash}
	err := r.client.Delete(context.Background(), torrent)
	if err != nil {
		return fmt.Errorf("could not delete torrent %q: %w", hash, err)
	}
	return nil
}
