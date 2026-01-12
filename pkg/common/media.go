package common

type MediaManager interface {
	CachedRetriever
	GetMedia() ([]Media, error)
	DeleteMediaFiles(mediaType MediaType, fileIds []int64, stopParentMonitoring bool) error
}

type MediaRetriever interface {
	GetMedia() ([]Media, error)
	SupportedMediaType() MediaType
	DeleteMediaFiles(fileIds []int64, stopParentMonitoring bool) error
}
