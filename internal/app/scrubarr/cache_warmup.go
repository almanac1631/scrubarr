package scrubarr

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/almanac1631/scrubarr/pkg/domain"
)

const cacheDir = "./cache"

func warmupCaches(saveCache, useCache bool, cachedRetrievers ...domain.CachedManager) error {
	if saveCache && useCache {
		return errors.New("cannot save and use cache simultaneously")
	}

	if useCache {
		return loadFromCache(cachedRetrievers)
	}

	errChan := make(chan error, len(cachedRetrievers))
	defer close(errChan)
	for _, cachedRetriever := range cachedRetrievers {
		go func() {
			errChan <- cachedRetriever.RefreshCache()
		}()
	}
	errorList := make([]error, 0)
	for i := 0; i < len(cachedRetrievers); i++ {
		if err := <-errChan; err != nil {
			errorList = append(errorList, err)
		}
	}
	if len(errorList) > 0 {
		return errors.Join(errorList...)
	}
	if saveCache {
		return saveToCache(cachedRetrievers)
	}
	return nil
}

func saveToCache(cachedRetrievers []domain.CachedManager) error {
	if err := os.Mkdir(cacheDir, 0777); err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create cache directory (%s): %w", cacheDir, err)
	}
	for _, cachedRetriever := range cachedRetrievers {
		writeCache := func() error {
			sha1Name := sha1.Sum([]byte(fmt.Sprintf("%T", cachedRetriever)))
			retrieverHash := hex.EncodeToString(sha1Name[:])
			file, err := os.Create(filepath.Join(cacheDir, retrieverHash))
			if err != nil {
				return err
			}
			defer func() {
				_ = file.Close()
			}()
			if err = cachedRetriever.SaveCache(file); err != nil {
				return err
			}
			return nil
		}
		if err := writeCache(); err != nil {
			return err
		}
	}
	return nil
}

func loadFromCache(cachedRetrievers []domain.CachedManager) error {
	for _, cachedRetriever := range cachedRetrievers {
		loadCache := func() error {
			sha1Name := sha1.Sum([]byte(fmt.Sprintf("%T", cachedRetriever)))
			retrieverHash := hex.EncodeToString(sha1Name[:])
			file, err := os.Open(filepath.Join(cacheDir, retrieverHash))
			if err != nil {
				return err
			}
			defer func() {
				_ = file.Close()
			}()
			return cachedRetriever.LoadCache(file)
		}
		if err := loadCache(); err != nil {
			return err
		}
	}
	return nil
}
