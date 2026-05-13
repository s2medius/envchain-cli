package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type lastpassHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// LastPassBackend retrieves secrets from LastPass Enterprise via the LastPass API.
type LastPassBackend struct {
	apiURL   string
	username string
	apiKey   string
	client   lastpassHTTPClient
}

type lastpassSecret struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NewLastPassBackend creates a new LastPassBackend from the given options map.
// Required options: username, api_key.
// Optional options: api_url (defaults to https://lastpass.com/enterpriseapi.php).
func NewLastPassBackend(opts map[string]string) (*LastPassBackend, error) {
	username, ok := opts["username"]
	if !ok || username == "" {
		return nil, fmt.Errorf("lastpass backend: missing required option 'username'")
	}
	apiKey, ok := opts["api_key"]
	if !ok || apiKey == "" {
		return nil, fmt.Errorf("lastpass backend: missing required option 'api_key'")
	}
	apiURL := opts["api_url"]
	if apiURL == "" {
		apiURL = "https://lastpass.com/enterpriseapi.php"
	}
	return &LastPassBackend{
		apiURL:   apiURL,
		username: username,
		apiKey:   apiKey,
		client:   &http.Client{},
	}, nil
}

// Get retrieves the value for the given key from LastPass.
// The key format is "folder/name" where folder is the shared folder name.
func (b *LastPassBackend) Get(key string) (string, error) {
	parts := strings.SplitN(key, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("lastpass backend: key %q must be in format 'folder/name'", key)
	}
	folder, name := parts[0], parts[1]

	body := fmt.Sprintf(`{"cmd":"getsharedfoldersdata","apiuser":%q,"apikey":%q,"folder":%q}`,
		b.username, b.apiKey, folder)
	req, err := http.NewRequest(http.MethodPost, b.apiURL, strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("lastpass backend: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("lastpass backend: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("lastpass backend: unexpected status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("lastpass backend: failed to read response: %w", err)
	}

	var secrets []lastpassSecret
	if err := json.Unmarshal(data, &secrets); err != nil {
		return "", fmt.Errorf("lastpass backend: failed to parse response: %w", err)
	}

	for _, s := range secrets {
		if s.Name == name {
			return s.Value, nil
		}
	}
	return "", fmt.Errorf("lastpass backend: key %q not found in folder %q", name, folder)
}

func (b *LastPassBackend) String() string {
	return fmt.Sprintf("lastpass(username=%s)", b.username)
}
