package scrubarr

import (
	"math/rand/v2"

	"github.com/almanac1631/scrubarr/internal/app/webserver"
)

type MockQuotaService struct {
}

func (m MockQuotaService) GetDiskQuota() (webserver.DiskQuota, error) {
	totalSpace := rand.Int64N(10_000_000_000_000)
	usedSpace := rand.Int64N(totalSpace)
	return webserver.DiskQuota{
		UsedSpacePercentage: float64(usedSpace) / float64(totalSpace) * 100.0,
		UsedSpace:           usedSpace,
		TotalSpace:          totalSpace,
		FreeSpace:           totalSpace - usedSpace,
	}, nil
}
