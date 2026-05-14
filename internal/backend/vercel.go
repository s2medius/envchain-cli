package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultVercelAPIURL = "https://api.vercel.com"

type vercelClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// VercelBackend retrieves environment variables from Vercel via the Vercel REST API.
type VercelBackend struct {
	token     string
	projectID string
	teamID    string
	apiURL    string
	client    vercelClient
}

// NewVercelBackend creates a new VercelBackend from the provided options map.
// Required keys: token, project_id.
// Optional keys: team_id, api_url.
func NewVercelBackend(opts map[string]string) (*VercelBackend, error) {
	token := opts["token"]
	if token == "" {
		return nil, fmt.Errorf("vercel: token is required")
	}
	projectID := opts["project_id"]
	if projectID == "" {
		return nil, fmt.Errorf("vercel: project_id is required")
	}
	apiURL := opts["api_url"]
	if apiURL == "" {
		apiURL = defaultVercelAPIURL
	}
	return &VercelBackend{
		token:     token,
		projectID: projectID,
		teamID:    opts["team_id"],
		apiURL:    apiURL,
		client:    &http.Client{},
	}, nil
}

// Get retrieves the value of the named environment variable from Vercel.
func (b *VercelBackend) Get(key string) (string, error) {
	url := fmt.Sprintf("%s/v9/projects/%s/env/%s", b.apiURL, b.projectID, key)
	if b.teamID != "" {
		url += "?teamId=" + b.teamID
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("vercel: failed to build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+b.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("vercel: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("vercel: key %q not found", key)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("vercel: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("vercel: failed to decode response: %w", err)
	}
	return result.Value, nil
}

func (b *VercelBackend) String() string {
	return fmt.Sprintf("vercel(project=%s)", b.projectID)
}
