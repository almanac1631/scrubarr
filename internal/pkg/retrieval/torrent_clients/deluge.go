package torrent_clients

import (
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval"
	delugeclient "github.com/gdm85/go-libdeluge"
	"path"
	"slices"
	"time"
)

type DelugeEntryRetriever struct {
	client             *delugeclient.ClientV2
	allowedFileEndings []string
}

func NewDelugeFilePathMappingRetriever(allowedFileEndings []string, hostname string, port uint, username string, password string) (*DelugeEntryRetriever, error) {
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
	return &DelugeEntryRetriever{client, allowedFileEndings}, nil
}

func (d *DelugeEntryRetriever) RetrieveEntries() (map[retrieval.EntryName]retrieval.Entry, error) {
	torrentsStatusList, err := d.client.TorrentsStatus(delugeclient.StateSeeding, []string{})
	if err != nil {
		return nil, fmt.Errorf("could not get torrents status: %w", err)
	}

	torrentEntryMap := map[retrieval.EntryName]retrieval.Entry{}
	for _, torrentStatus := range torrentsStatusList {
		iterTorrentStatusList := d.parseDelugeTorrentStatus(torrentStatus)
		for _, torrentStatusEntry := range iterTorrentStatusList {
			torrentEntryMap[torrentStatusEntry.Name] = torrentStatusEntry
		}
	}
	return torrentEntryMap, nil
}

func (d *DelugeEntryRetriever) parseDelugeTorrentStatus(torrentStatus *delugeclient.TorrentStatus) []retrieval.Entry {
	entryList := make([]retrieval.Entry, 0)
	for _, file := range torrentStatus.Files {
		if !slices.Contains(d.allowedFileEndings, path.Ext(file.Path)) {
			continue
		}
		filePath := path.Join(torrentStatus.DownloadLocation, file.Path)
		downloadedAt := time.Unix(torrentStatus.CompletedTime, 0).In(time.UTC)
		name := path.Base(filePath)
		entry := retrieval.Entry{
			Name: retrieval.EntryName(name),
			AdditionalData: TorrentClientEntry{
				TorrentClientName: "deluge",
				TorrentName:       torrentStatus.Name,
				DownloadFilePath:  filePath,
				DownloadedAt:      downloadedAt,
				Ratio:             torrentStatus.Ratio,
				FileSizeBytes:     file.Size,
				TrackerHost:       torrentStatus.TrackerHost,
			},
		}
		entryList = append(entryList, entry)
	}
	return entryList
}
