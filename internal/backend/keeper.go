package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// KeeperBackend retrieves secrets from Keeper Secrets Manager via its REST gateway.
type KeeperBackend struct {
	baseURL string
	token   string
	client  *http.Client
}

type keeperSecretResponse struct {
	Data struct {
		Fields []struct {
			Type  string   `json:"type"`
			Value []string `json:"value"`
		} `json:"fields"`
	} `json:"data"`
}

// NewKeeperBackend creates a KeeperBackend from the given options map.
// Required keys: token, record_uid
// Optional keys: base_url (default: https://keepersecurity.com/api/rest/sm/v1)
func NewKeeperBackend(opts map[string]string) (*KeeperBackend, error) {
	token := opts["token"]
	if token == "" {
		return nil, fmt.Errorf("keeper: missing required option 'token'")
	}
	baseURL := opts["base_url"]
	if baseURL == "" {
		baseURL = "https://keepersecurity.com/api/rest/sm/v1"
	}
	return &KeeperBackend{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{},
	}, nil
}

// Get retrieves the value of a secret field. The key format is "<record_uid>/<field_type>".
func (k *KeeperBackend) Get(key string) (string, error) {
	recordUID, fieldType, err := splitTwoParts(key, "/")
	if err != nil {
		return "", fmt.Errorf("keeper: invalid key format %q, expected '<record_uid>/<field_type>'", key)
	}

	url := fmt.Sprintf("%s/get_secret?uid=%s", k.baseURL, recordUID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("keeper: failed to build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+k.token)
	req.Header.Set("Accept", "application/json")

	resp, err := k.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("keeper: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("keeper: record %q not found", recordUID)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("keeper: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result keeperSecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("keeper: failed to decode response: %w", err)
	}

	for _, field := range result.Data.Fields {
		if field.Type == fieldType && len(field.Value) > 0 {
			return field.Value[0], nil
		}
	}
	return "", fmt.Errorf("keeper: field type %q not found in record %q", fieldType, recordUID)
}

func (k *KeeperBackend) String() string {
	return fmt.Sprintf("keeper(%s)", k.baseURL)
}
