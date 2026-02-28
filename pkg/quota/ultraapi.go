package quota

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/almanac1631/scrubarr/internal/app/webserver"
	"github.com/almanac1631/scrubarr/pkg/ultraapi"
)

type UltraApiQuotaService struct {
	*sync.Mutex
	ultraApi        *ultraapi.Instance
	lastDiskQuota   webserver.DiskQuota
	lastRetrievedTs time.Time
}

func NewUltraApiQuotaService(endpoint string, authToken []byte) *UltraApiQuotaService {
	instance := ultraapi.New(endpoint, authToken)
	return &UltraApiQuotaService{
		Mutex:    &sync.Mutex{},
		ultraApi: instance,
	}
}

func (service *UltraApiQuotaService) GetDiskQuota() (webserver.DiskQuota, error) {
	service.Lock()
	defer service.Unlock()
	if !service.lastRetrievedTs.IsZero() && time.Since(service.lastRetrievedTs) < time.Second*30 {
		return service.lastDiskQuota, nil
	}
	slog.Debug("Retrieving fresh disk quota from Ultra API...")
	ultraDiskQuota, err := service.ultraApi.GetDiskQuota()
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
	service.lastDiskQuota = webserver.DiskQuota{
		UsedSpacePercentage: float64(usedSpace) / float64(totalSpace) * 100.0,
		UsedSpace:           usedSpace,
		TotalSpace:          totalSpace,
		FreeSpace:           totalSpace - usedSpace,
	}
	service.lastRetrievedTs = time.Now()
	return service.lastDiskQuota, nil
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
