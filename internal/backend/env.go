package backend

import (
	"fmt"
	"os"
	"strings"
)

// EnvBackend resolves secrets from the current process environment.
type EnvBackend struct {
	prefix string
}

// NewEnvBackend creates an EnvBackend. If prefix is non-empty, only
// variables with that prefix are considered, and the prefix is stripped
// from the key before lookup.
func NewEnvBackend(prefix string) *EnvBackend {
	return &EnvBackend{prefix: prefix}
}

// Get returns the value of the environment variable identified by key.
// If a prefix is configured it is prepended to key before lookup.
func (e *EnvBackend) Get(key string) (string, error) {
	lookup := key
	if e.prefix != "" {
		lookup = e.prefix + key
	}

	val, ok := os.LookupEnv(lookup)
	if !ok {
		return "", ErrSecretNotFound{Key: key}
	}
	return val, nil
}

// List returns all keys available in the environment, optionally filtered
// and stripped by the configured prefix.
func (e *EnvBackend) List() ([]string, error) {
	var keys []string
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) < 1 {
			continue
		}
		name := parts[0]
		if e.prefix != "" {
			if !strings.HasPrefix(name, e.prefix) {
				continue
			}
			name = strings.TrimPrefix(name, e.prefix)
		}
		keys = append(keys, name)
	}
	return keys, nil
}

// String returns a human-readable description of the backend.
func (e *EnvBackend) String() string {
	if e.prefix != "" {
		return fmt.Sprintf("env(prefix=%s)", e.prefix)
	}
	return "env"
}
