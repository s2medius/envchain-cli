package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type bitwardenClient interface {
	GetItem(itemID string) (map[string]string, error)
}

type bitwardenHTTPClient struct {
	baseURL    string
	accessToken string
	httpClient  *http.Client
}

func (c *bitwardenHTTPClient) GetItem(itemID string) (map[string]string, error) {
	url := fmt.Sprintf("%s/object/item/%s", strings.TrimRight(c.baseURL, "/"), itemID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bitwarden: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Fields []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"fields"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("bitwarden: failed to decode response: %w", err)
	}

	fields := make(map[string]string, len(result.Data.Fields))
	for _, f := range result.Data.Fields {
		fields[f.Name] = f.Value
	}
	return fields, nil
}

// BitwardenBackend retrieves secrets from a Bitwarden Secrets Manager via the
// local Bitwarden CLI REST server (bw serve).
type BitwardenBackend struct {
	client bitwardenClient
	itemID string
}

// NewBitwardenBackend creates a new BitwardenBackend.
// Required config keys: access_token, item_id.
// Optional: base_url (default: http://localhost:8087).
func NewBitwardenBackend(cfg map[string]string) (*BitwardenBackend, error) {
	token := cfg["access_token"]
	if token == "" {
		return nil, fmt.Errorf("bitwarden: access_token is required")
	}
	itemID := cfg["item_id"]
	if itemID == "" {
		return nil, fmt.Errorf("bitwarden: item_id is required")
	}
	baseURL := cfg["base_url"]
	if baseURL == "" {
		baseURL = "http://localhost:8087"
	}
	return &BitwardenBackend{
		client: &bitwardenHTTPClient{
			baseURL:     baseURL,
			accessToken: token,
			httpClient:  &http.Client{},
		},
		itemID: itemID,
	}, nil
}

func (b *BitwardenBackend) Get(key string) (string, error) {
	fields, err := b.client.GetItem(b.itemID)
	if err != nil {
		return "", err
	}
	if fields == nil {
		return "", fmt.Errorf("bitwarden: item %q not found", b.itemID)
	}
	val, ok := fields[key]
	if !ok {
		return "", fmt.Errorf("bitwarden: key %q not found in item %q", key, b.itemID)
	}
	return val, nil
}

func (b *BitwardenBackend) String() string {
	return fmt.Sprintf("bitwarden(item=%s)", b.itemID)
}
