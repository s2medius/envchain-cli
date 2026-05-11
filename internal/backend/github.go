package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type githubSecretsClient interface {
	GetSecret(ctx context.Context, owner, repo, secretName string) (string, error)
}

// GitHubBackend reads secrets from GitHub Actions secrets via the GitHub API.
type GitHubBackend struct {
	client githubSecretsClient
	owner  string
	repo   string
}

type defaultGitHubClient struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

func (c *defaultGitHubClient) GetSecret(ctx context.Context, owner, repo, secretName string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/actions/secrets/%s", c.baseURL, owner, repo, secretName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("secret %q not found", secretName)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	val, ok := result["value"].(string)
	if !ok {
		return "", fmt.Errorf("secret value not available (secrets are write-only via API)")
	}
	return val, nil
}

// NewGitHubBackend creates a new GitHubBackend from config options.
func NewGitHubBackend(options map[string]string) (*GitHubBackend, error) {
	token := options["token"]
	if token == "" {
		return nil, fmt.Errorf("github backend requires 'token'")
	}
	owner := options["owner"]
	if owner == "" {
		return nil, fmt.Errorf("github backend requires 'owner'")
	}
	repo := options["repo"]
	if repo == "" {
		return nil, fmt.Errorf("github backend requires 'repo'")
	}
	baseURL := options["base_url"]
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	return &GitHubBackend{
		client: &defaultGitHubClient{
			httpClient: &http.Client{},
			token:      token,
			baseURL:    baseURL,
		},
		owner: owner,
		repo:  repo,
	}, nil
}

func (b *GitHubBackend) Get(ctx context.Context, key string) (string, error) {
	return b.client.GetSecret(ctx, b.owner, b.repo, key)
}

func (b *GitHubBackend) String() string {
	return fmt.Sprintf("github(%s/%s)", b.owner, b.repo)
}

func init() {
	_ = strings.ToLower // ensure strings import used
}
