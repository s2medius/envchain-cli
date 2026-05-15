package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultNetlifyAPIURL = "https://api.netlify.com/api/v1"

type netlifyClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NetlifyBackend fetches environment variables from Netlify site environment.
type NetlifyBackend struct {
	token   string
	siteID  string
	apiURL  string
	client  netlifyClient
}

// NewNetlifyBackend creates a new NetlifyBackend.
func NewNetlifyBackend(opts map[string]string) (*NetlifyBackend, error) {
	token := opts["token"]
	if token == "" {
		return nil, fmt.Errorf("netlify: missing required option 'token'")
	}
	siteID := opts["site_id"]
	if siteID == "" {
		return nil, fmt.Errorf("netlify: missing required option 'site_id'")
	}
	apiURL := opts["api_url"]
	if apiURL == "" {
		apiURL = defaultNetlifyAPIURL
	}
	return &NetlifyBackend{
		token:  token,
		siteID: siteID,
		apiURL: apiURL,
		client: &http.Client{},
	}, nil
}

func (b *NetlifyBackend) Get(key string) (string, error) {
	url := fmt.Sprintf("%s/sites/%s/env/%s", b.apiURL, b.siteID, key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("netlify: failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+b.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("netlify: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("netlify: key %q not found", key)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("netlify: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("netlify: failed to read response: %w", err)
	}

	var result struct {
		Values []struct {
			Value   string `json:"value"`
			Context string `json:"context"`
		} `json:"values"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("netlify: failed to parse response: %w", err)
	}
	for _, v := range result.Values {
		if v.Context == "all" || v.Context == "production" || v.Context == "" {
			return v.Value, nil
		}
	}
	if len(result.Values) > 0 {
		return result.Values[0].Value, nil
	}
	return "", fmt.Errorf("netlify: no value found for key %q", key)
}

func (b *NetlifyBackend) String() string {
	return fmt.Sprintf("netlify(site=%s)", b.siteID)
}
