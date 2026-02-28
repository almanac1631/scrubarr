package inventory

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/almanac1631/scrubarr/pkg/domain"
)

func (s *Service) loadManagerCacheFromDisk(manager domain.CachedManager) error {
	cacheDir := os.Getenv("SCRUBARR_CACHE_DIR")
	if cacheDir == "" {
		cacheDir = "./cache"
	}
	logger := slog.With("manager", fmt.Sprintf("%T", manager))
	logger.Info("Using cache manager to refresh cache")
	sha1Name := sha1.Sum([]byte(fmt.Sprintf("%T", manager)))
	retrieverHash := hex.EncodeToString(sha1Name[:])
	file, err := os.Open(filepath.Join(cacheDir, retrieverHash))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	return manager.LoadCache(file)
}

func (s *Service) saveManagerCacheToDisk(manager domain.CachedManager) error {
	cacheDir := os.Getenv("SCRUBARR_CACHE_DIR")
	if cacheDir == "" {
		cacheDir = "./cache"
	}
	if err := os.Mkdir(cacheDir, 0777); err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create cache directory (%s): %w", cacheDir, err)
	}
	logger := slog.With("manager", fmt.Sprintf("%T", manager))
	logger.Info("Saving cache manager to disk")
	sha1Name := sha1.Sum([]byte(fmt.Sprintf("%T", manager)))
	retrieverHash := hex.EncodeToString(sha1Name[:])
	file, err := os.Create(filepath.Join(cacheDir, retrieverHash))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	return manager.SaveCache(file)
}

func (s *Service) RefreshCache() error {
	s.Lock()
	defer s.Unlock()
	errChan := make(chan error)
	defer close(errChan)
	refreshManagerCache := func(manager domain.CachedManager) {
		if s.useCache {
			errChan <- s.loadManagerCacheFromDisk(manager)
			return
		}
		err := manager.RefreshCache()
		if err != nil {
			errChan <- err
			return
		}
		if s.saveCache {
			errChan <- s.saveManagerCacheToDisk(manager)
		}
	}
	go refreshManagerCache(s.mediaSourceManager)
	go refreshManagerCache(s.torrentSourceManager)
	err := <-errChan
	if refreshErr := <-errChan; refreshErr != nil {
		if err != nil {
			err = errors.Join(err, refreshErr)
		} else {
			err = refreshErr
		}
	}
	if err != nil {
		return fmt.Errorf("refresh cache failed: %w", err)
	}

	media, err := s.mediaSourceManager.GetMedia()
	if err != nil {
		return fmt.Errorf("unable to get media: %w", err)
	}
	torrents, err := s.torrentSourceManager.GetTorrents()
	if err != nil {
		return fmt.Errorf("unable to get torrents: %w", err)
	}
	linkedMediaList, err := s.linker.LinkMedia(media, torrents)
	if err != nil {
		return fmt.Errorf("unable to link media with torrents: %w", err)
	}

	s.enrichedLinkedMediaCache = make([]enrichedLinkedMedia, len(linkedMediaList))
	for i, linkedMedia := range linkedMediaList {
		evaluationReport, err := s.retentionPolicy.Evaluate(linkedMedia)
		if err != nil {
			return fmt.Errorf("unable to evaluate retention policy: %w", err)
		}
		s.enrichedLinkedMediaCache[i] = enrichedLinkedMedia{
			linkedMedia:      linkedMedia,
			evaluationReport: evaluationReport,
			size:             getSize(linkedMedia),
			added:            getAdded(linkedMedia),
		}
	}
	return nil
}
