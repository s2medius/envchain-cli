package backend

import "fmt"

// Provider defines the interface for secret backends.
type Provider interface {
	// GetSecret retrieves a secret value by key.
	GetSecret(key string) (string, error)
	// Name returns the backend provider name.
	Name() string
}

// Type represents a supported backend type.
type Type string

const (
	TypeEnv    Type = "env"
	TypeFile   Type = "file"
	TypeVault  Type = "vault"
	TypeAWSSSM Type = "aws-ssm"
)

// ErrSecretNotFound is returned when a secret key does not exist in the backend.
type ErrSecretNotFound struct {
	Backend string
	Key     string
}

func (e *ErrSecretNotFound) Error() string {
	return fmt.Sprintf("secret %q not found in backend %q", e.Key, e.Backend)
}

// ErrUnsupportedBackend is returned when an unknown backend type is specified.
type ErrUnsupportedBackend struct {
	Type string
}

func (e *ErrUnsupportedBackend) Error() string {
	return fmt.Sprintf("unsupported backend type: %q", e.Type)
}

// New constructs a Provider for the given backend type and options.
func New(backendType Type, opts map[string]string) (Provider, error) {
	switch backendType {
	case TypeEnv:
		return NewEnvProvider(), nil
	case TypeFile:
		path, ok := opts["path"]
		if !ok || path == "" {
			return nil, fmt.Errorf("file backend requires 'path' option")
		}
		return NewFileProvider(path), nil
	default:
		return nil, &ErrUnsupportedBackend{Type: string(backendType)}
	}
}
