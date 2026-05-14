package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultPulumiAPIURL = "https://api.pulumi.com"

type pulumiClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// PulumiBackend retrieves secrets from Pulumi ESC (Environments, Secrets, and Configuration).
type PulumiBackend struct {
	token      string
	org        string
	environment string
	apiURL     string
	client     pulumiClient
}

// NewPulumiBackend creates a new PulumiBackend.
func NewPulumiBackend(cfg map[string]string, client pulumiClient) (*PulumiBackend, error) {
	token := cfg["token"]
	if token == "" {
		return nil, fmt.Errorf("pulumi: missing required config key: token")
	}
	org := cfg["org"]
	if org == "" {
		return nil, fmt.Errorf("pulumi: missing required config key: org")
	}
	env := cfg["environment"]
	if env == "" {
		return nil, fmt.Errorf("pulumi: missing required config key: environment")
	}
	apiURL := cfg["api_url"]
	if apiURL == "" {
		apiURL = defaultPulumiAPIURL
	}
	if client == nil {
		client = &http.Client{}
	}
	return &PulumiBackend{
		token:       token,
		org:         org,
		environment: env,
		apiURL:      apiURL,
		client:      client,
	}, nil
}

// Get retrieves a secret value by key from Pulumi ESC.
func (b *PulumiBackend) Get(ctx context.Context, key string) (string, error) {
	url := fmt.Sprintf("%s/api/esc/environments/%s/%s/open", b.apiURL, b.org, b.environment)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("pulumi: failed to build request: %w", err)
	}
	req.Header.Set("Authorization", "token "+b.token)
	req.Header.Set("Accept", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("pulumi: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("pulumi: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("pulumi: failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("pulumi: failed to parse response: %w", err)
	}

	values, ok := result["values"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("pulumi: unexpected response shape")
	}
	val, ok := values[key]
	if !ok {
		return "", fmt.Errorf("pulumi: key %q not found", key)
	}
	return fmt.Sprintf("%v", val), nil
}

// String returns a human-readable description of the backend.
func (b *PulumiBackend) String() string {
	return fmt.Sprintf("pulumi(org=%s, environment=%s)", b.org, b.environment)
}
