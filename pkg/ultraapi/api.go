package ultraapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Instance struct {
	endpoint   string
	httpClient *http.Client
	authToken  []byte
}

func New(endpoint string, authToken []byte) *Instance {
	return &Instance{endpoint: endpoint, httpClient: http.DefaultClient, authToken: authToken}
}

func (instance Instance) GetTotalStats() (*TotalStats, error) {
	var totalStats TotalStats
	err := instance.request("total-stats", &totalStats)
	if err != nil {
		return nil, fmt.Errorf("failed to get total stats: %w", err)
	}
	return &totalStats, nil
}

func (instance Instance) GetTraffic() (*Traffic, error) {
	var traffic Traffic
	err := instance.request("get-traffic", &traffic)
	if err != nil {
		return nil, fmt.Errorf("failed to get traffic: %w", err)
	}
	return &traffic, nil
}

func (instance Instance) GetDiskQuota() (*DiskQuota, error) {
	var diskQuota DiskQuota
	err := instance.request("get-diskquota", &diskQuota)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk quota: %w", err)
	}
	return &diskQuota, nil
}

func (instance Instance) request(endpointPath string, receivingValue any) (err error) {
	requestUrl, err := url.JoinPath(instance.endpoint, endpointPath)
	if err != nil {
		return fmt.Errorf("failed to join url: %w", err)
	}
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+string(instance.authToken))
	resp, err := instance.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err != nil {
			//goland:noinspection GoUnhandledErrorResult
			resp.Body.Close()
			return
		}
		err = resp.Body.Close()
	}()
	var respBytes []byte
	respBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return ErrUnexpectedApiResp{resp.StatusCode, respBytes}
	}
	err = json.Unmarshal(respBytes, receivingValue)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal response: %w", err)
		return
	}
	return
}
