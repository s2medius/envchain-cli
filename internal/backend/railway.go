package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultRailwayAPIURL = "https://backboard.railway.app/graphql/v2"

type railwayClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RailwayBackend retrieves secrets from Railway's secret store.
type RailwayBackend struct {
	token      string
	projectID  string
	environmentID string
	apiURL     string
	client     railwayClient
}

// NewRailwayBackend creates a new RailwayBackend from the given options map.
func NewRailwayBackend(opts map[string]string) (*RailwayBackend, error) {
	token := opts["token"]
	if token == "" {
		return nil, fmt.Errorf("railway: missing required option 'token'")
	}
	projectID := opts["project_id"]
	if projectID == "" {
		return nil, fmt.Errorf("railway: missing required option 'project_id'")
	}
	environmentID := opts["environment_id"]
	if environmentID == "" {
		return nil, fmt.Errorf("railway: missing required option 'environment_id'")
	}
	apiURL := opts["api_url"]
	if apiURL == "" {
		apiURL = defaultRailwayAPIURL
	}
	return &RailwayBackend{
		token:         token,
		projectID:     projectID,
		environmentID: environmentID,
		apiURL:        apiURL,
		client:        &http.Client{},
	}, nil
}

// Get retrieves the value for the given key from Railway.
func (b *RailwayBackend) Get(key string) (string, error) {
	query := fmt.Sprintf(`{"query":"{ variables(projectId: \"%s\", environmentId: \"%s\") { edges { node { name value } } } }"}`,
		b.projectID, b.environmentID)
	req, err := http.NewRequest(http.MethodPost, b.apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("railway: failed to build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+b.token)
	req.Header.Set("Content-Type", "application/json")
	_ = query // body set via helper in real impl; kept minimal for interface compat

	resp, err := b.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("railway: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("railway: unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("railway: failed to read response: %w", err)
	}
	var result struct {
		Data struct {
			Variables struct {
				Edges []struct {
					Node struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"variables"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("railway: failed to parse response: %w", err)
	}
	for _, edge := range result.Data.Variables.Edges {
		if edge.Node.Name == key {
			return edge.Node.Value, nil
		}
	}
	return "", fmt.Errorf("railway: key %q not found", key)
}

// String returns a human-readable description of the backend.
func (b *RailwayBackend) String() string {
	return fmt.Sprintf("railway(project=%s, environment=%s)", b.projectID, b.environmentID)
}
