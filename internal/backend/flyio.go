package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const flyioDefaultAPIURL = "https://api.fly.io"

type flyioClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// FlyIOBackend fetches secrets from Fly.io Secrets API.
type FlyIOBackend struct {
	token  string
	appID  string
	apiURL string
	client flyioClient
}

// NewFlyIOBackend creates a new FlyIOBackend.
func NewFlyIOBackend(opts map[string]string) (*FlyIOBackend, error) {
	token, ok := opts["token"]
	if !ok || token == "" {
		return nil, fmt.Errorf("flyio: missing required option 'token'")
	}
	appID, ok := opts["app_id"]
	if !ok || appID == "" {
		return nil, fmt.Errorf("flyio: missing required option 'app_id'")
	}
	apiURL := flyioDefaultAPIURL
	if u, ok := opts["api_url"]; ok && u != "" {
		apiURL = strings.TrimRight(u, "/")
	}
	return &FlyIOBackend{
		token:  token,
		appID:  appID,
		apiURL: apiURL,
		client: &http.Client{},
	}, nil
}

func (b *FlyIOBackend) Get(key string) (string, error) {
	url := fmt.Sprintf("%s/v1/apps/%s/secrets", b.apiURL, b.appID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("flyio: failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+b.token)
	req.Header.Set("Accept", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("flyio: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("flyio: unexpected status %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("flyio: failed to decode response: %w", err)
	}

	val, ok := result[key]
	if !ok {
		return "", fmt.Errorf("flyio: key %q not found", key)
	}
	return val, nil
}

func (b *FlyIOBackend) String() string {
	return fmt.Sprintf("flyio(app=%s)", b.appID)
}
