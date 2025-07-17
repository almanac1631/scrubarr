package torrent_clients

import (
	"errors"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	delugeclient "github.com/gdm85/go-libdeluge"
	"log/slog"
	"path"
	"slices"
	"time"
)

var _ common.EntryRetriever = (*DelugeEntryRetriever)(nil)

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

func (d *DelugeEntryRetriever) RetrieveEntries() (common.RetrieverEntries, error) {
	torrentsStatusList, err := d.client.TorrentsStatus(delugeclient.StateSeeding, []string{})
	if err != nil {
		return nil, fmt.Errorf("could not get torrents status: %w", err)
	}

	torrentEntryMap := common.RetrieverEntries{}
	for hash, torrentStatus := range torrentsStatusList {
		iterTorrentStatusList := d.parseDelugeTorrentStatus(hash, torrentStatus)
		for _, torrentStatusEntry := range iterTorrentStatusList {
			torrentEntryMap[torrentStatusEntry.Name] = torrentStatusEntry
		}
	}
	return torrentEntryMap, nil
}

func (d *DelugeEntryRetriever) parseDelugeTorrentStatus(hash string, torrentStatus *delugeclient.TorrentStatus) []common.Entry {
	entryList := make([]common.Entry, 0)
	for _, file := range torrentStatus.Files {
		if !slices.Contains(d.allowedFileEndings, path.Ext(file.Path)) {
			continue
		}
		filePath := path.Join(torrentStatus.DownloadLocation, file.Path)
		downloadedAt := time.Unix(torrentStatus.CompletedTime, 0).In(time.UTC)
		name := path.Base(filePath)
		entry := common.Entry{
			Name: common.EntryName(name),
			AdditionalData: TorrentClientEntry{
				ID:                hash,
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

func (d *DelugeEntryRetriever) DeleteEntry(id any) error {
	torrentID, ok := id.(string)
	if !ok {
		return fmt.Errorf("could not convert id to string")
	}
	ok, err := d.client.RemoveTorrent(torrentID, true)
	if ok || err == nil {
		return nil
	}
	var wrappedErr delugeclient.RPCError
	if errors.As(err, &wrappedErr) && wrappedErr.ExceptionType == "InvalidTorrentError" {
		slog.Warn("torrent not found in deluge, assuming it was already removed", "id", torrentID)
		return nil
	}
	return fmt.Errorf("could not remove torrent: %w", err)
}
