package ultraapi

type TotalStats struct {
	ServiceStatsInfo struct {
		FreeStorageBytes           int64   `json:"free_storage_bytes"`
		FreeStorageGb              float64 `json:"free_storage_gb"`
		LastTrafficReset           string  `json:"last_traffic_reset"`
		NextTrafficReset           string  `json:"next_traffic_reset"`
		TotalStorageUnit           string  `json:"total_storage_unit"`
		TotalStorageValue          int     `json:"total_storage_value"`
		TrafficAvailablePercentage float64 `json:"traffic_available_percentage"`
		TrafficUsedPercentage      float64 `json:"traffic_used_percentage"`
		UsedStorageUnit            string  `json:"used_storage_unit"`
		UsedStorageValue           int     `json:"used_storage_value"`
	} `json:"service_stats_info"`
}

type DiskQuota struct {
	StorageInfo struct {
		FreeStorageBytes  int64   `json:"free_storage_bytes"`
		FreeStorageGb     float64 `json:"free_storage_gb"`
		TotalStorageUnit  string  `json:"total_storage_unit"`
		TotalStorageValue int     `json:"total_storage_value"`
		UsedStorageUnit   string  `json:"used_storage_unit"`
		UsedStorageValue  int     `json:"used_storage_value"`
	} `json:"Storage Info"`
}

type Traffic struct {
	TrafficInfo struct {
		LastTrafficReset           string  `json:"last_traffic_reset"`
		NextTrafficReset           string  `json:"next_traffic_reset"`
		TrafficAvailablePercentage float64 `json:"traffic_available_percentage"`
		TrafficUsedPercentage      float64 `json:"traffic_used_percentage"`
	} `json:"Traffic info"`
}
