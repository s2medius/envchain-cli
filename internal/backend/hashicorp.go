package backend

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
)

// HashiCorpCloudClient is the interface for HCP Vault Secrets API calls.
type HashiCorpCloudClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HashiCorpBackend fetches secrets from HCP Vault Secrets.
type HashiCorpBackend struct {
	client      HashiCorpCloudClient
	baseURL     string
	token       string
	orgID       string
	projectID   string
	appName     string
}

// NewHashiCorpBackend creates a new HashiCorpBackend from the provided options.
func NewHashiCorpBackend(opts map[string]string) (*HashiCorpBackend, error) {
	token := opts["token"]
	if token == "" {
		return nil, fmt.Errorf("hashicorp: token is required")
	}
	orgID := opts["org_id"]
	if orgID == "" {
		return nil, fmt.Errorf("hashicorp: org_id is required")
	}
	projectID := opts["project_id"]
	if projectID == "" {
		return nil, fmt.Errorf("hashicorp: project_id is required")
	}
	appName := opts["app_name"]
	if appName == "" {
		return nil, fmt.Errorf("hashicorp: app_name is required")
	}
	baseURL := opts["base_url"]
	if baseURL == "" {
		baseURL = "https://api.cloud.hashicorp.com"
	}
	return &HashiCorpBackend{
		client:    &http.Client{},
		baseURL:   baseURL,
		token:     token,
		orgID:     orgID,
		projectID: projectID,
		appName:   appName,
	}, nil
}

func (h *HashiCorpBackend) Get(key string) (string, error) {
	url := fmt.Sprintf(
		"%s/secrets/2023-11-28/organizations/%s/projects/%s/apps/%s/secrets/%s:open",
		h.baseURL, h.orgID, h.projectID, h.appName, key,
	)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("hashicorp: failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+h.token)
	resp, err := h.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("hashicorp: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("hashicorp: secret %q not found", key)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("hashicorp: unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("hashicorp: failed to read response: %w", err)
	}
	var result struct {
		Secret struct {
			StaticVersion struct {
				Value string `json:"value"`
			} `json:"static_version"`
		} `json:"secret"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("hashicorp: failed to parse response: %w", err)
	}
	return result.Secret.StaticVersion.Value, nil
}

func (h *HashiCorpBackend) String() string {
	return fmt.Sprintf("hashicorp(app=%s)", h.appName)
}
