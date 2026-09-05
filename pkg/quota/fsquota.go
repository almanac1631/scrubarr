package quota

import (
	"fmt"
	"os/user"

	"github.com/almanac1631/scrubarr/internal/app/webserver"
	"github.com/parkervcp/fsquota"
)

type FsQuotaService struct{}

func NewFsQuotaService() *FsQuotaService {
	return &FsQuotaService{}
}

func (service *FsQuotaService) GetDiskQuota() (webserver.DiskQuota, error) {
	currentUser, err := user.Current()
	if err != nil {
		return webserver.DiskQuota{}, fmt.Errorf("getting current user: %w", err)
	}

	info, err := fsquota.GetUserInfo(currentUser.HomeDir, currentUser)
	if err != nil {
		return webserver.DiskQuota{}, fmt.Errorf("getting filesystem quota for %q: %w", currentUser.HomeDir, err)
	}

	usedSpace := int64(info.BytesUsed)
	totalSpace := int64(info.Bytes.GetHard())

	if totalSpace <= 0 {
		return webserver.DiskQuota{}, fmt.Errorf("invalid hard quota for %q: %d", currentUser.HomeDir, totalSpace)
	}

	freeSpace := totalSpace - usedSpace

	return webserver.DiskQuota{
		UsedSpacePercentage: float64(usedSpace) / float64(totalSpace) * 100.0,
		UsedSpace:           usedSpace,
		TotalSpace:          totalSpace,
		FreeSpace:           freeSpace,
	}, nil
}
