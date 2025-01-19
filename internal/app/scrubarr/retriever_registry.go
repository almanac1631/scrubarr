package scrubarr

import (
	"fmt"
	"github.com/almanac1631/scrubarr/internal/app/common"
	"github.com/almanac1631/scrubarr/internal/pkg/config"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/arr_apps"
	_ "github.com/almanac1631/scrubarr/internal/pkg/retrieval/arr_apps"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/folder_scanning"
	_ "github.com/almanac1631/scrubarr/internal/pkg/retrieval/folder_scanning"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/torrent_clients"
	_ "github.com/almanac1631/scrubarr/internal/pkg/retrieval/torrent_clients"
	"github.com/knadh/koanf/v2"
	"log/slog"
)

type instantiateFuncType func(koanf *koanf.Koanf, allowedFileEndings []string) (retrieval.EntryRetriever, error)

var retrieverRegistry = &cachedRetrieverRegistry{retrieverRegistry: common.MapBasedRetrieverRegistry{}}

func registerRetrievers(koanf *koanf.Koanf) error {
	for retrieverInfo, instantiateFunc := range map[common.RetrieverInfo]instantiateFuncType{
		{Category: "folder", SoftwareName: "folder"}:           instantiateFolderScanner,
		{Category: "arr_app", SoftwareName: "sonarr"}:          instantiateSonarrRetriever,
		{Category: "arr_app", SoftwareName: "radarr"}:          instantiateRadarrRetriever,
		{Category: "torrent_client", SoftwareName: "deluge"}:   instantiateDelugeRetriever,
		{Category: "torrent_client", SoftwareName: "rtorrent"}: instantiateRtorrentRetriever,
	} {
		slog.Info("loading retrievers", "retrieverCategory", retrieverInfo.Category, "retrieverSoftwareName", retrieverInfo.SoftwareName)
		err := checkAndRegisterRetriever(koanf, retrieverInfo, instantiateFunc)
		if err != nil {
			return fmt.Errorf("could not register retriever type %q: %w", retrieverInfo, err)
		}
	}
	return nil
}

func checkAndRegisterRetriever(koanf *koanf.Koanf, retrieverInfoGeneral common.RetrieverInfo, instantiateMethod instantiateFuncType) error {
	path := fmt.Sprintf("connections.%s", retrieverInfoGeneral.SoftwareName)
	retrieverConfigs, err := config.GetEntry[map[string]interface{}](koanf.Get, path)
	if err != nil {
		return err
	}
	allowedFileEndings := koanf.MustStrings("general.allowed_file_endings")
	for retrieverName, _ := range retrieverConfigs {
		retrieverPath := fmt.Sprintf("%s.%s", path, retrieverName)
		folderEntryEnabled, err := config.GetEntry[bool](koanf.Get, fmt.Sprintf("%s.enabled", retrieverPath))
		if err != nil {
			return err
		}
		if !folderEntryEnabled {
			slog.Info("retriever disabled", "retrieverType", retrieverInfoGeneral, "path", retrieverPath)
			continue
		}
		entryConfig := koanf.Cut(retrieverPath)
		retriever, err := instantiateMethod(entryConfig, allowedFileEndings)
		if err != nil {
			return fmt.Errorf("failed to register retriever %q: %w", retrieverName, err)
		}
		retrieverInfo := retrieverInfoGeneral
		retrieverInfo.Name = retrieverName
		retrieverRegistry.retrieverRegistry[retrieverInfo] = retriever
	}
	return nil
}

func instantiateFolderScanner(koanf *koanf.Koanf, allowedFileEndings []string) (retrieval.EntryRetriever, error) {
	folderPath := koanf.MustString("path")
	retriever, err := folder_scanning.NewFolderScanner(allowedFileEndings, folderPath)
	if err != nil {
		return nil, err
	}
	return retriever, nil
}

func instantiateSonarrRetriever(koanf *koanf.Koanf, allowedFileEndings []string) (retrieval.EntryRetriever, error) {
	hostname := koanf.MustString("hostname")
	apiKey := koanf.MustString("api_key")
	retriever, err := arr_apps.NewSonarrMediaRetriever(allowedFileEndings, hostname, apiKey)
	if err != nil {
		return nil, err
	}
	return retriever, nil
}

func instantiateDelugeRetriever(koanf *koanf.Koanf, allowedFileEndings []string) (retrieval.EntryRetriever, error) {
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

func instantiateRadarrRetriever(koanf *koanf.Koanf, allowedFileEndings []string) (retrieval.EntryRetriever, error) {
	hostname := koanf.MustString("hostname")
	apiKey := koanf.MustString("api_key")
	retriever, err := arr_apps.NewRadarrMediaRetriever(allowedFileEndings, hostname, apiKey)
	if err != nil {
		return nil, err
	}
	return retriever, nil
}

func instantiateRtorrentRetriever(koanf *koanf.Koanf, allowedFileEndings []string) (retrieval.EntryRetriever, error) {
	hostname := koanf.MustString("hostname")
	username := koanf.MustString("username")
	password := koanf.MustString("password")
	retriever, err := torrent_clients.NewRtorrentEntryRetriever(allowedFileEndings, hostname, username, password)
	if err != nil {
		return nil, err
	}
	return retriever, nil
}
