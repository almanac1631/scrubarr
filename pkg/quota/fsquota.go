//go:build !linux

package quota

import (
	"github.com/almanac1631/scrubarr/internal/app/webserver"
)

type FsQuotaService struct{}

func NewFsQuotaService() *FsQuotaService {
	return &FsQuotaService{}
}

func (service *FsQuotaService) GetDiskQuota() (webserver.DiskQuota, error) {
	panic("disk quota not supported on this platform")
}
