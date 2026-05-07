package backend

import (
	"fmt"

	"github.com/envchain-cli/envchain-cli/internal/config"
)

// Backend is the interface all secret backends must satisfy.
type Backend interface {
	Get(key string) (string, error)
	String() string
}

// New constructs a Backend from a config.Backend definition.
func New(cfg config.Backend) (Backend, error) {
	switch cfg.Type {
	case "env":
		return NewEnvBackend(cfg.Options)
	case "file":
		return NewFileBackend(cfg.Options)
	case "vault":
		return NewVaultBackend(cfg.Options)
	case "ssm":
		return NewSSMBackend(cfg.Options)
	case "secretsmanager":
		return NewSecretsManagerBackend(cfg.Options)
	case "gcp":
		return NewGCPBackend(cfg.Options)
	default:
		return nil, fmt.Errorf("unsupported backend type: %q", cfg.Type)
	}
}
