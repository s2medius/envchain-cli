package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type opConnectClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// OnePasswordConnectBackend retrieves secrets from a 1Password Connect Server.
type OnePasswordConnectBackend struct {
	baseURL string
	token   string
	vaultID string
	client  opConnectClient
}

type opConnectItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Fields []struct {
		Label string `json:"label"`
		Value string `json:"value"`
	} `json:"fields"`
}

// NewOnePasswordConnectBackend creates a new 1Password Connect backend.
func NewOnePasswordConnectBackend(cfg map[string]string) (*OnePasswordConnectBackend, error) {
	token := cfg["token"]
	if token == "" {
		return nil, fmt.Errorf("onepassword_connect: missing required field 'token'")
	}
	baseURL := cfg["url"]
	if baseURL == "" {
		return nil, fmt.Errorf("onepassword_connect: missing required field 'url'")
	}
	vaultID := cfg["vault_id"]
	if vaultID == "" {
		return nil, fmt.Errorf("onepassword_connect: missing required field 'vault_id'")
	}
	return &OnePasswordConnectBackend{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		vaultID: vaultID,
		client:  &http.Client{},
	}, nil
}

// Get retrieves a secret field value. key format: "<item-title>/<field-label>"
func (b *OnePasswordConnectBackend) Get(key string) (string, error) {
	parts := strings.SplitN(key, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("onepassword_connect: key must be in format 'item-title/field-label', got %q", key)
	}
	itemTitle, fieldLabel := parts[0], parts[1]

	url := fmt.Sprintf("%s/v1/vaults/%s/items?filter=title eq \"%s\"", b.baseURL, b.vaultID, itemTitle)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("onepassword_connect: building request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+b.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("onepassword_connect: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("onepassword_connect: unexpected status %d for item %q", resp.StatusCode, itemTitle)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("onepassword_connect: reading response: %w", err)
	}

	var items []opConnectItem
	if err := json.Unmarshal(body, &items); err != nil {
		return "", fmt.Errorf("onepassword_connect: parsing response: %w", err)
	}
	if len(items) == 0 {
		return "", fmt.Errorf("onepassword_connect: item %q not found", itemTitle)
	}

	for _, f := range items[0].Fields {
		if strings.EqualFold(f.Label, fieldLabel) {
			return f.Value, nil
		}
	}
	return "", fmt.Errorf("onepassword_connect: field %q not found in item %q", fieldLabel, itemTitle)
}

func (b *OnePasswordConnectBackend) String() string {
	return fmt.Sprintf("onepassword_connect(url=%s, vault=%s)", b.baseURL, b.vaultID)
}
