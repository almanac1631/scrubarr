package webserver

type DiskQuota struct {
	UsedSpacePercentage              float64
	UsedSpace, TotalSpace, FreeSpace int64
}

type QuotaService interface {
	GetDiskQuota() (DiskQuota, error)
}
