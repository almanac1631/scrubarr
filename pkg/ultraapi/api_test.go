package ultraapi

import (
	"fmt"
	"os"
	"testing"
)

func TestApiClient(t *testing.T) {
	endpoint := os.Getenv("ULTRA_ENDPOINT")
	if endpoint == "" {
		t.Skipf("No ULTRA_ENDPOINT set, skipping")
	}
	authToken := []byte(os.Getenv("ULTRA_AUTH_TOKEN"))
	apiClient := New(endpoint, authToken)
	t.Run("can get total stats", func(t *testing.T) {
		totalStats, err := apiClient.GetTotalStats()
		if err != nil {
			t.Errorf("failed to get total stats: %v", err)
		}
		fmt.Printf("%T: %+v\n", totalStats, totalStats)
	})
	t.Run("can get traffic", func(t *testing.T) {
		traffic, err := apiClient.GetTraffic()
		if err != nil {
			t.Errorf("failed to get traffic: %v", err)
		}
		fmt.Printf("%T: %+v\n", traffic, traffic)
	})
	t.Run("can get disk quota", func(t *testing.T) {
		diskQuota, err := apiClient.GetDiskQuota()
		if err != nil {
			t.Errorf("failed to get disk quota: %v", err)
		}
		fmt.Printf("%T: %+v\n", diskQuota, diskQuota)
	})
}
