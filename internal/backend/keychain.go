//go:build darwin

package backend

import (
	"fmt"
	"os/exec"
	"strings"
)

// KeychainBackend retrieves secrets from macOS Keychain using the `security` CLI.
type KeychainBackend struct {
	service string
	execCommand func(name string, args ...string) ([]byte, error)
}

// NewKeychainBackend creates a new KeychainBackend.
// Requires `service` to be set in the backend config.
func NewKeychainBackend(cfg map[string]string) (*KeychainBackend, error) {
	service, ok := cfg["service"]
	if !ok || service == "" {
		return nil, fmt.Errorf("keychain backend: missing required config key 'service'")
	}
	return &KeychainBackend{
		service: service,
		execCommand: func(name string, args ...string) ([]byte, error) {
			return exec.Command(name, args...).Output()
		},
	}, nil
}

// Get retrieves a secret by account name from the macOS Keychain.
func (k *KeychainBackend) Get(key string) (string, error) {
	out, err := k.execCommand(
		"security",
		"find-generic-password",
		"-s", k.service,
		"-a", key,
		"-w",
	)
	if err != nil {
		return "", fmt.Errorf("keychain backend: key %q not found in service %q: %w", key, k.service, err)
	}
	return strings.TrimRight(string(out), "\n"), nil
}

// String returns a human-readable description of the backend.
func (k *KeychainBackend) String() string {
	return fmt.Sprintf("keychain(service=%s)", k.service)
}
