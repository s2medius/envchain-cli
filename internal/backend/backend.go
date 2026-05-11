package backend

import (
	"context"
	"fmt"

	"github.com/user/envchain-cli/internal/config"
)

// Backend is the interface all secret backends must implement.
type Backend interface {
	Get(ctx context.Context, key string) (string, error)
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
	case "azure":
		return NewAzureBackend(cfg.Options)
	case "1password":
		return NewOnePasswordBackend(cfg.Options)
	case "doppler":
		return NewDopplerBackend(cfg.Options)
	case "infisical":
		return NewInfisicalBackend(cfg.Options)
	case "github":
		return NewGitHubBackend(cfg.Options)
	default:
		return nil, fmt.Errorf("unsupported backend type: %q", cfg.Type)
	}
}
