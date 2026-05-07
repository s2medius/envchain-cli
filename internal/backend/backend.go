package backend

import (
	"context"
	"fmt"
)

// Backend is the interface that all secret backends must implement.
type Backend interface {
	Get(ctx context.Context, key string) (string, error)
	String() string
}

// New creates a Backend from the given type and options map.
func New(backendType string, opts map[string]string) (Backend, error) {
	switch backendType {
	case "env":
		return NewEnvBackend(opts), nil
	case "file":
		return NewFileBackend(opts)
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
	default:
		return nil, fmt.Errorf("unsupported backend type: %q", backendType)
	}
}
