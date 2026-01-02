package common

type CachedRetriever interface {
	RefreshCache() error
}
