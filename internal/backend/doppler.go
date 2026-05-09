package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const dopplerAPIBase = "https://api.doppler.com/v3"

type dopplerClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DopplerBackend fetches secrets from Doppler.
type DopplerBackend struct {
	token   string
	project string
	config  string
	client  dopplerClient
}

// NewDopplerBackend creates a new DopplerBackend.
// Requires token, project, and config (environment) to be set.
func NewDopplerBackend(token, project, config string, client dopplerClient) (*DopplerBackend, error) {
	if token == "" {
		return nil, fmt.Errorf("doppler: token is required")
	}
	if project == "" {
		return nil, fmt.Errorf("doppler: project is required")
	}
	if config == "" {
		return nil, fmt.Errorf("doppler: config (environment) is required")
	}
	if client == nil {
		client = &http.Client{}
	}
	return &DopplerBackend{
		token:   token,
		project: project,
		config:  config,
		client:  client,
	}, nil
}

// Get retrieves a single secret value by name from Doppler.
func (d *DopplerBackend) Get(key string) (string, error) {
	url := fmt.Sprintf("%s/configs/config/secret?project=%s&config=%s&name=%s",
		dopplerAPIBase, d.project, d.config, key)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("doppler: failed to build request: %w", err)
	}
	req.SetBasicAuth(d.token, "")
	req.Header.Set("Accept", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("doppler: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("doppler: secret %q not found", key)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("doppler: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Secret struct {
			Value struct {
				Raw string `json:"raw"`
			} `json:"value"`
		} `json:"secret"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("doppler: failed to decode response: %w", err)
	}
	return result.Secret.Value.Raw, nil
}

// String returns a human-readable description of the backend.
func (d *DopplerBackend) String() string {
	return fmt.Sprintf("doppler(project=%s, config=%s)", d.project, d.config)
}
