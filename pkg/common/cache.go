package common

import "io"

type CachedRetriever interface {
	RefreshCache() error

	SaveCache(writer io.Writer) error

	LoadCache(reader io.ReadSeeker) error
}
