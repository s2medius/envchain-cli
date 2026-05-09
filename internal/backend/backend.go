package backend

import (
	"fmt"

	"github.com/yourusername/envchain-cli/internal/config"
)

// Backend is the interface that all secret backends must implement.
type Backend interface {
	Get(key string) (string, error)
	String() string
}

// New creates a Backend from a config.Backend definition.
func New(cfg config.Backend) (Backend, error) {
	opts := cfg.Options
	if opts == nil {
		opts = map[string]string{}
	}

	switch cfg.Type {
	case "env":
		return NewEnvBackend(opts), nil

	case "file":
		path, ok := opts["path"]
		if !ok || path == "" {
			return nil, fmt.Errorf("file backend: missing required option 'path'")
		}
		return NewFileBackend(path)

	case "vault":
		return NewVaultBackend(opts)

	case "ssm":
		return NewSSMBackend(opts)

	case "secretsmanager":
		return NewSecretsManagerBackend(opts)

	case "gcp":
		return NewGCPBackend(opts)

	case "azure":
		return NewAzureBackend(opts)

	case "1password":
		return NewOnePasswordBackend(opts)

	case "doppler":
		return NewDopplerBackend(opts)

	case "infisical":
		return NewInfisicalBackend(opts)

	default:
		return nil, fmt.Errorf("unsupported backend type: %q", cfg.Type)
	}
}
