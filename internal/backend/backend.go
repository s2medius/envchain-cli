package backend

import (
	"fmt"

	"github.com/envchain-cli/envchain-cli/internal/config"
)

// Backend is the interface all secret backends must implement.
type Backend interface {
	Get(key string) (string, error)
	String() string
}

// New creates a Backend from a config.Backend definition.
func New(cfg config.Backend) (Backend, error) {
	switch cfg.Type {
	case "env":
		return NewEnvBackend(cfg.Prefix), nil

	case "file":
		path, _ := cfg.Options["path"]
		if path == "" {
			return nil, fmt.Errorf("file backend requires 'path' option")
		}
		return NewFileBackend(path)

	case "vault":
		address, _ := cfg.Options["address"]
		token, _ := cfg.Options["token"]
		path, _ := cfg.Options["path"]
		return NewVaultBackend(address, token, path)

	case "ssm":
		prefix, _ := cfg.Options["prefix"]
		region, _ := cfg.Options["region"]
		return NewSSMBackend(prefix, region)

	case "secretsmanager":
		secretID, _ := cfg.Options["secret_id"]
		region, _ := cfg.Options["region"]
		return NewSecretsManagerBackend(secretID, region)

	default:
		return nil, fmt.Errorf("unsupported backend type: %q", cfg.Type)
	}
}
