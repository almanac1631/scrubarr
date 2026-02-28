package scrubarr

import (
	"fmt"

	"github.com/almanac1631/scrubarr/internal/app/webserver"
	"github.com/almanac1631/scrubarr/pkg/quota"
	"github.com/knadh/koanf/v2"
)

func getQuotaService(k *koanf.Koanf) (webserver.QuotaService, error) {
	provider := k.MustString("quota.provider")
	if provider == "ultraapi" {
		endpoint := k.MustString("quota.ultraapi.endpoint")
		apiKey := k.MustBytes("quota.ultraapi.api_key")
		return quota.NewUltraApiQuotaService(endpoint, apiKey), nil
	}
	return nil, fmt.Errorf("unknown quota provider: %q", provider)
}
