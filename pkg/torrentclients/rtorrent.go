package torrentclients

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"slices"
	"strings"

	"github.com/almanac1631/scrubarr/pkg/domain"
	"github.com/autobrr/go-rtorrent"
	"github.com/autobrr/go-rtorrent/xmlrpc"
)

var _ domain.TorrentSource = (*RtorrentRetriever)(nil)

type RtorrentRetriever struct {
	client       *rtorrent.Client
	xmlrpcClient *xmlrpc.Client
	dryRun       bool
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
	xmlrpcClient := xmlrpc.NewClient(xmlrpc.Config{
		Addr:      hostname,
		BasicUser: username,
		BasicPass: password,
	})
	return &RtorrentRetriever{client: client, xmlrpcClient: xmlrpcClient, dryRun: dryRun}, nil
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
			Label:    torrent.Label,
			Added:    torrent.Finished,
			Files:    []*domain.TorrentFile{},
			Trackers: []string{},
			Ratio:    torrent.Ratio,
		}
		basePath, err := retriever.getBasePath(torrent.Hash)
		if err != nil {
			return nil, err
		}
		if basePath == "" {
			basePath, err = retriever.getFallbackBasePath(torrent)
			if err != nil {
				return nil, err
			}
		}
		torrentFiles, err := retriever.client.GetFiles(context.Background(), torrent)
		if err != nil {
			return nil, fmt.Errorf("could not get torrent files from rtorrent: %w", err)
		}
		for _, torrentFile := range torrentFiles {
			torrentEntry.Files = append(torrentEntry.Files, &domain.TorrentFile{
				Path:         torrentFile.Path,
				AbsolutePath: rtorrentFileAbsolutePath(basePath, torrentFile.Path),
				Size:         int64(torrentFile.Size),
			})
		}
		torrentEntries = append(torrentEntries, torrentEntry)
		torrentTrackers, err := retriever.client.GetTrackers(context.Background(), torrent)
		if err != nil {
			return nil, fmt.Errorf("could not get trackerresolver from rtorrent: %w", err)
		}
		for _, tracker := range torrentTrackers {
			host := extractHost(tracker)
			if !slices.Contains(torrentEntry.Trackers, host) {
				torrentEntry.Trackers = append(torrentEntry.Trackers, host)
			}
		}
	}
	return torrentEntries, nil
}

func (retriever *RtorrentRetriever) getBasePath(hash string) (string, error) {
	results, err := retriever.xmlrpcClient.Call(context.Background(), "d.base_path", hash)
	if err != nil {
		return "", fmt.Errorf("could not get base path of torrent %q from rtorrent: %w", hash, err)
	}
	values, ok := results.([]any)
	if !ok || len(values) == 0 {
		return "", fmt.Errorf("unexpected d.base_path response for torrent %q: %v", hash, results)
	}
	basePath, ok := values[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected d.base_path value for torrent %q: %v", hash, values[0])
	}
	return basePath, nil
}

// For multi-file torrents d.directory is already the root folder; for single-file ones it is the parent folder.
func (retriever *RtorrentRetriever) getFallbackBasePath(torrent rtorrent.Torrent) (string, error) {
	results, err := retriever.xmlrpcClient.Call(context.Background(), "d.is_multi_file", torrent.Hash)
	if err != nil {
		return "", fmt.Errorf("could not check if torrent %q is multi-file on rtorrent: %w", torrent.Hash, err)
	}
	values, ok := results.([]any)
	if !ok || len(values) == 0 {
		return "", fmt.Errorf("unexpected d.is_multi_file response for torrent %q: %v", torrent.Hash, results)
	}
	isMultiFile, ok := values[0].(int)
	if !ok {
		return "", fmt.Errorf("unexpected d.is_multi_file value for torrent %q: %v", torrent.Hash, values[0])
	}
	if isMultiFile != 0 {
		return torrent.Path, nil
	}
	return joinRemotePath(torrent.Path, torrent.Name), nil
}

// f.path may or may not start with the root folder name; returns "" for a relative base path.
func rtorrentFileAbsolutePath(basePath string, filePath string) string {
	separatorIndex := strings.LastIndexAny(basePath, `/\`)
	if separatorIndex < 0 {
		return ""
	}
	parentDir := basePath[:separatorIndex+1]
	baseName := basePath[separatorIndex+1:]
	switch {
	case filePath == baseName:
		return basePath
	case strings.HasPrefix(filePath, baseName+"/") || strings.HasPrefix(filePath, baseName+`\`):
		return joinRemotePath(parentDir, filePath)
	default:
		return joinRemotePath(basePath, filePath)
	}
}

func extractHost(tracker string) string {
	parsed, err := url.Parse(tracker)
	if err != nil || parsed.Hostname() == "" {
		return tracker
	}
	return parsed.Hostname()
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
