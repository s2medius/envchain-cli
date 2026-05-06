package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// VaultBackend reads secrets from HashiCorp Vault KV v2.
type VaultBackend struct {
	address   string
	token     string
	mountPath string
	secretPath string
	client    *http.Client
}

type vaultResponse struct {
	Data struct {
		Data map[string]string `json:"data"`
	} `json:"data"`
}

// NewVaultBackend creates a VaultBackend from the provided options.
func NewVaultBackend(opts map[string]string) (*VaultBackend, error) {
	address, ok := opts["address"]
	if !ok || address == "" {
		return nil, fmt.Errorf("vault backend: missing required option 'address'")
	}
	token, ok := opts["token"]
	if !ok || token == "" {
		return nil, fmt.Errorf("vault backend: missing required option 'token'")
	}
	mount := opts["mount"]
	if mount == "" {
		mount = "secret"
	}
	path := opts["path"]
	if path == "" {
		return nil, fmt.Errorf("vault backend: missing required option 'path'")
	}
	return &VaultBackend{
		address:    strings.TrimRight(address, "/"),
		token:      token,
		mountPath:  mount,
		secretPath: path,
		client:     &http.Client{},
	}, nil
}

func (v *VaultBackend) fetchSecrets() (map[string]string, error) {
	url := fmt.Sprintf("%s/v1/%s/data/%s", v.address, v.mountPath, v.secretPath)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("vault backend: creating request: %w", err)
	}
	req.Header.Set("X-Vault-Token", v.token)
	resp, err := v.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault backend: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrSecretNotFound{Path: v.secretPath}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault backend: unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("vault backend: reading response: %w", err)
	}
	var vr vaultResponse
	if err := json.Unmarshal(body, &vr); err != nil {
		return nil, fmt.Errorf("vault backend: parsing response: %w", err)
	}
	return vr.Data.Data, nil
}

func (v *VaultBackend) Get(key string) (string, error) {
	secrets, err := v.fetchSecrets()
	if err != nil {
		return "", err
	}
	val, ok := secrets[key]
	if !ok {
		return "", ErrSecretNotFound{Path: key}
	}
	return val, nil
}

func (v *VaultBackend) List() (map[string]string, error) {
	return v.fetchSecrets()
}

func (v *VaultBackend) String() string {
	return fmt.Sprintf("vault(%s/%s/%s)", v.address, v.mountPath, v.secretPath)
}
