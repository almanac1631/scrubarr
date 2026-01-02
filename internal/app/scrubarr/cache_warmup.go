package scrubarr

import (
	"errors"

	"github.com/almanac1631/scrubarr/pkg/common"
)

func warmupCaches(cachedRetrievers ...common.CachedRetriever) error {
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
	return nil
}
