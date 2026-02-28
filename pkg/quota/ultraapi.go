package quota

import (
	"fmt"

	"github.com/almanac1631/scrubarr/internal/app/webserver"
	"github.com/almanac1631/scrubarr/pkg/ultraapi"
)

type UltraApiQuotaService struct {
	ultraApi *ultraapi.Instance
}

func NewUltraApiQuotaService(endpoint string, authToken []byte) *UltraApiQuotaService {
	instance := ultraapi.New(endpoint, authToken)
	return &UltraApiQuotaService{ultraApi: instance}
}

func (u UltraApiQuotaService) GetDiskQuota() (webserver.DiskQuota, error) {
	ultraDiskQuota, err := u.ultraApi.GetDiskQuota()
	if err != nil {
		return webserver.DiskQuota{}, fmt.Errorf("error getting disk quota from ultra api: %v", err)
	}
	storageInfo := ultraDiskQuota.StorageInfo
	totalSpace, err := parseStorageValue(int64(storageInfo.TotalStorageValue), storageInfo.TotalStorageUnit)
	if err != nil {
		return webserver.DiskQuota{}, err
	}
	usedSpace, err := parseStorageValue(int64(storageInfo.UsedStorageValue), storageInfo.UsedStorageUnit)
	if err != nil {
		return webserver.DiskQuota{}, err
	}
	return webserver.DiskQuota{
		UsedSpacePercentage: float64(usedSpace) / float64(totalSpace) * 100.0,
		UsedSpace:           usedSpace,
		TotalSpace:          totalSpace,
		FreeSpace:           totalSpace - usedSpace,
	}, nil
}

func parseStorageValue(value int64, unit string) (int64, error) {
	if unit == "G" {
		return value * 1024 * 1024 * 1024, nil
	} else if unit == "M" {
		return value * 1024 * 1024, nil
	} else if unit == "K" {
		return value * 1024, nil
	}
	return -1, fmt.Errorf("unknown storage value unit %q", unit)
}
