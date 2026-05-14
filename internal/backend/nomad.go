package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// nomadClient is an interface for Nomad API calls (enables testing).
type nomadClient interface {
	Get(url string) (*http.Response, error)
}

// NomadBackend retrieves secrets from HashiCorp Nomad's Variables API.
type NomadBackend struct {
	address string
	token   string
	path    string
	client  nomadClient
}

// NewNomadBackend creates a new NomadBackend from the given options map.
// Required keys: address, token, path.
func NewNomadBackend(opts map[string]string) (*NomadBackend, error) {
	address := opts["address"]
	if address == "" {
		return nil, fmt.Errorf("nomad: missing required option 'address'")
	}
	token := opts["token"]
	if token == "" {
		return nil, fmt.Errorf("nomad: missing required option 'token'")
	}
	path := opts["path"]
	if path == "" {
		return nil, fmt.Errorf("nomad: missing required option 'path'")
	}
	return &NomadBackend{
		address: address,
		token:   token,
		path:    path,
		client:  &http.Client{},
	}, nil
}

// Get retrieves the value for the given key from the Nomad Variable at the
// configured path. The variable's Items map is queried by key.
func (b *NomadBackend) Get(key string) (string, error) {
	url := fmt.Sprintf("%s/v1/var/%s", b.address, b.path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("nomad: failed to build request: %w", err)
	}
	req.Header.Set("X-Nomad-Token", b.token)

	hc, ok := b.client.(*http.Client)
	var resp *http.Response
	if ok {
		resp, err = hc.Do(req)
	} else {
		resp, err = b.client.Get(url)
	}
	if err != nil {
		return "", fmt.Errorf("nomad: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("nomad: variable path %q not found", b.path)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("nomad: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("nomad: failed to read response: %w", err)
	}

	var result struct {
		Items map[string]string `json:"Items"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("nomad: failed to parse response: %w", err)
	}

	val, ok2 := result.Items[key]
	if !ok2 {
		return "", fmt.Errorf("nomad: key %q not found in variable %q", key, b.path)
	}
	return val, nil
}

func (b *NomadBackend) String() string {
	return fmt.Sprintf("nomad(%s)", b.path)
}
