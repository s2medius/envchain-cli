package backend

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// conjurClient is an interface for testability.
type conjurClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ConjurBackend retrieves secrets from CyberArk Conjur.
type ConjurBackend struct {
	applianceURL string
	account      string
	token        string
	client       conjurClient
}

// NewConjurBackend creates a new ConjurBackend from the provided options map.
// Required keys: address, account, token.
func NewConjurBackend(opts map[string]string) (*ConjurBackend, error) {
	address := opts["address"]
	if address == "" {
		return nil, fmt.Errorf("conjur backend: missing required option 'address'")
	}
	account := opts["account"]
	if account == "" {
		return nil, fmt.Errorf("conjur backend: missing required option 'account'")
	}
	token := opts["token"]
	if token == "" {
		return nil, fmt.Errorf("conjur backend: missing required option 'token'")
	}
	return &ConjurBackend{
		applianceURL: strings.TrimRight(address, "/"),
		account:      account,
		token:        token,
		client:       &http.Client{},
	}, nil
}

// Get retrieves a secret value by variable ID (e.g. "myapp/db/password").
func (b *ConjurBackend) Get(key string) (string, error) {
	url := fmt.Sprintf("%s/secrets/%s/variable/%s", b.applianceURL, b.account, key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("conjur backend: failed to build request: %w", err)
	}
	req.Header.Set("Authorization", "Token token=\"" + b.token + "\"")

	resp, err := b.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("conjur backend: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("conjur backend: secret not found: %s", key)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("conjur backend: unexpected status %d for key %s", resp.StatusCode, key)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("conjur backend: failed to read response: %w", err)
	}
	return string(body), nil
}

// String returns a human-readable description of the backend.
func (b *ConjurBackend) String() string {
	return fmt.Sprintf("conjur(%s, account=%s)", b.applianceURL, b.account)
}
