package domain

import "io"

type CachedManager interface {
	RefreshCache() error

	SaveCache(writer io.Writer) error

	LoadCache(reader io.ReadSeeker) error
}
