package scrubarr

import (
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/arr_apps"
	_ "github.com/almanac1631/scrubarr/internal/pkg/retrieval/arr_apps"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/torrent_clients"
	_ "github.com/almanac1631/scrubarr/internal/pkg/retrieval/torrent_clients"
	"github.com/knadh/koanf/v2"
	"log/slog"
)

type instantiateFuncType func(koanf *koanf.Koanf, allowedFileEndings []string) (common.EntryRetriever, error)

type instantiateRetrieverInfo struct {
	category     string
	softwareName string
}

var instantiateFunctions = map[instantiateRetrieverInfo]instantiateFuncType{
	{category: "arr_app", softwareName: "sonarr"}:          instantiateSonarrRetriever,
	{category: "arr_app", softwareName: "radarr"}:          instantiateRadarrRetriever,
	{category: "torrent_client", softwareName: "deluge"}:   instantiateDelugeRetriever,
	{category: "torrent_client", softwareName: "rtorrent"}: instantiateRtorrentRetriever,
}

func initializeEntryRetrievers(koanf *koanf.Koanf) (map[common.RetrieverInfo]common.EntryRetriever, error) {
	entryRetrievers := map[common.RetrieverInfo]common.EntryRetriever{}
	for keyRetrieverInfo, instantiateFunc := range instantiateFunctions {
		err := checkAndRegisterRetriever(koanf, keyRetrieverInfo.category, keyRetrieverInfo.softwareName, entryRetrievers, instantiateFunc)
		if err != nil {
			return nil, fmt.Errorf("could not register retriever type %q: %w", keyRetrieverInfo, err)
		}
	}
	return entryRetrievers, nil
}

func checkAndRegisterRetriever(koanf *koanf.Koanf, category, softwareName string, retrieverRegistry map[common.RetrieverInfo]common.EntryRetriever, instantiateMethod instantiateFuncType) error {
	path := fmt.Sprintf("connections.%s", softwareName)
	retrieverConfigs, ok := koanf.Get(path).(map[string]interface{})
	if !ok {
		return fmt.Errorf("could not find and parse retriever config at %q", path)
	}
	allowedFileEndings := koanf.MustStrings("general.allowed_file_endings")
	for retrieverName, _ := range retrieverConfigs {
		logger := slog.With("category", category, "softwareName", softwareName, "path", path)
		retrieverPath := fmt.Sprintf("%s.%s", path, retrieverName)

		folderEntryEnabled := koanf.Bool(fmt.Sprintf("%s.enabled", retrieverPath))
		if !folderEntryEnabled {
			logger.Info("retriever disabled")
			continue
		}
		entryConfig := koanf.Cut(retrieverPath)
		retriever, err := instantiateMethod(entryConfig, allowedFileEndings)
		if err != nil {
			return fmt.Errorf("failed to register retriever %q: %w", retrieverName, err)
		}
		retrieverInfo := common.RetrieverInfo{
			Category:     category,
			SoftwareName: softwareName,
			Name:         retrieverName,
		}
		retrieverRegistry[retrieverInfo] = retriever
	}
	return nil
}

func instantiateSonarrRetriever(koanf *koanf.Koanf, allowedFileEndings []string) (common.EntryRetriever, error) {
	hostname := koanf.MustString("hostname")
	apiKey := koanf.MustString("api_key")
	retriever, err := arr_apps.NewSonarrMediaRetriever(allowedFileEndings, hostname, apiKey)
	if err != nil {
		return nil, err
	}
	return retriever, nil
}

func instantiateDelugeRetriever(koanf *koanf.Koanf, allowedFileEndings []string) (common.EntryRetriever, error) {
	hostname := koanf.MustString("hostname")
	port := uint(koanf.MustInt("port"))
	username := koanf.MustString("username")
	password := koanf.MustString("password")
	retriever, err := torrent_clients.NewDelugeFilePathMappingRetriever(allowedFileEndings, hostname, port, username, password)
	if err != nil {
		return nil, err
	}
	return retriever, nil
}

func instantiateRadarrRetriever(koanf *koanf.Koanf, allowedFileEndings []string) (common.EntryRetriever, error) {
	hostname := koanf.MustString("hostname")
	apiKey := koanf.MustString("api_key")
	retriever, err := arr_apps.NewRadarrMediaRetriever(allowedFileEndings, hostname, apiKey)
	if err != nil {
		return nil, err
	}
	return retriever, nil
}

func instantiateRtorrentRetriever(koanf *koanf.Koanf, allowedFileEndings []string) (common.EntryRetriever, error) {
	hostname := koanf.MustString("hostname")
	username := koanf.MustString("username")
	password := koanf.MustString("password")
	retriever, err := torrent_clients.NewRtorrentEntryRetriever(allowedFileEndings, hostname, username, password)
	if err != nil {
		return nil, err
	}
	return retriever, nil
}
