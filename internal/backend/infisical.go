package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultInfisicalBaseURL = "https://app.infisical.com"

type infisicalClient interface {
	GetSecret(projectID, environment, secretName string) (string, error)
}

type infisicalHTTPClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func (c *infisicalHTTPClient) GetSecret(projectID, environment, secretName string) (string, error) {
	url := fmt.Sprintf("%s/api/v3/secrets/raw/%s?workspaceId=%s&environment=%s",
		c.baseURL, secretName, projectID, environment)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("secret %q not found", secretName)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("infisical: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Secret struct {
			SecretValue string `json:"secretValue"`
		} `json:"secret"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("infisical: failed to parse response: %w", err)
	}
	return result.Secret.SecretValue, nil
}

// InfisicalBackend fetches secrets from Infisical.
type InfisicalBackend struct {
	client      infisicalClient
	projectID   string
	environment string
}

// NewInfisicalBackend creates a new InfisicalBackend from the given options.
func NewInfisicalBackend(opts map[string]string) (*InfisicalBackend, error) {
	token := opts["token"]
	if token == "" {
		return nil, fmt.Errorf("infisical backend: missing required option 'token'")
	}
	projectID := opts["project_id"]
	if projectID == "" {
		return nil, fmt.Errorf("infisical backend: missing required option 'project_id'")
	}
	environment := opts["environment"]
	if environment == "" {
		environment = "dev"
	}
	baseURL := opts["base_url"]
	if baseURL == "" {
		baseURL = defaultInfisicalBaseURL
	}
	return &InfisicalBackend{
		client: &infisicalHTTPClient{
			baseURL:    baseURL,
			token:      token,
			httpClient: &http.Client{},
		},
		projectID:   projectID,
		environment: environment,
	}, nil
}

func (b *InfisicalBackend) Get(key string) (string, error) {
	return b.client.GetSecret(b.projectID, b.environment, key)
}

func (b *InfisicalBackend) String() string {
	return fmt.Sprintf("infisical(project=%s, env=%s)", b.projectID, b.environment)
}
