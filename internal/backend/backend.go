package backend

import (
	"fmt"
)

// Backend is the interface all secret backends must satisfy.
type Backend interface {
	Get(key string) (string, error)
	List() (map[string]string, error)
	String() string
}

// ErrSecretNotFound is returned when a requested secret key does not exist.
type ErrSecretNotFound struct {
	Path string
}

func (e ErrSecretNotFound) Error() string {
	return fmt.Sprintf("secret not found: %s", e.Path)
}

// BackendConfig holds the type and options for a single backend.
type BackendConfig struct {
	Type    string            `yaml:"type"`
	Options map[string]string `yaml:"options"`
}

// New constructs a Backend from a BackendConfig.
func New(cfg BackendConfig) (Backend, error) {
	switch cfg.Type {
	case "env":
		return NewEnvBackend(cfg.Options), nil
	case "file":
		return NewFileBackend(cfg.Options)
	case "vault":
		return NewVaultBackend(cfg.Options)
	default:
		return nil, fmt.Errorf("unsupported backend type: %q", cfg.Type)
	}
}
