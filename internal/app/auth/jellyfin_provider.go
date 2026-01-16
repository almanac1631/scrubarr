package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var _ Provider = (*JellyfinProvider)(nil)

const (
	deviceId = "d28be177795b5542eda282b41c38fe60e1e6583a779481bbb420cd01e433ef45"
	version  = "unset"
)

type JellyfinProvider struct {
	baseUrl string
}

func NewJellyfinProvider(baseUrl string) *JellyfinProvider {
	return &JellyfinProvider{baseUrl: baseUrl}
}

func (provider JellyfinProvider) CheckCredentials(username string, password []byte) (bool, error) {
	url := provider.baseUrl + "/Users/AuthenticateByName"
	jsonData, _ := json.Marshal(map[string]string{
		"Username": username,
		"Pw":       string(password),
	})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")
	embyAuthorization := fmt.Sprintf("MediaBrowser Client=\"Scrubarr\", Device=\"Scrubarr\", DeviceId=%q, Version=%q", deviceId, version)
	req.Header.Set("X-Emby-Authorization", embyAuthorization)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send auth request to jellyfin: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == 401 {
		return false, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read auth response from jellyfin: %w", err)
	}

	var result JellyfinResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("failed to parse auth response from jellyfin: %w", err)
	}

	if result.AccessToken == "" {
		return false, nil
	}

	if !result.User.Policy.IsAdministrator {
		return false, nil
	}
	return true, nil
}

type JellyfinResponse struct {
	User struct {
		Name   string `json:"Name"`
		Policy struct {
			IsAdministrator bool `json:"IsAdministrator"`
		} `json:"Policy"`
	} `json:"User"`
	AccessToken string `json:"AccessToken"`
}

func (provider JellyfinProvider) Name() string {
	return "jellyfin"
}
