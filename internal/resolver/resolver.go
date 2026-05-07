package resolver

import (
	"fmt"

	"github.com/yourusername/envchain-cli/internal/backend"
	"github.com/yourusername/envchain-cli/internal/config"
)

// Resolver resolves environment variables from multiple backends
// based on the provided configuration.
type Resolver struct {
	backends []backend.Backend
}

// New creates a Resolver from the given config, initializing each backend.
func New(cfg *config.Config) (*Resolver, error) {
	var backends []backend.Backend

	for _, bc := range cfg.Backends {
		b, err := backend.New(bc)
		if err != nil {
			return nil, fmt.Errorf("resolver: failed to init backend %q: %w", bc.Type, err)
		}
		backends = append(backends, b)
	}

	return &Resolver{backends: backends}, nil
}

// Resolve looks up each requested key across all backends in order,
// returning the first value found. Returns an error if any key is missing.
func (r *Resolver) Resolve(keys []string) (map[string]string, error) {
	result := make(map[string]string, len(keys))

	for _, key := range keys {
		val, found, err := r.lookup(key)
		if err != nil {
			return nil, fmt.Errorf("resolver: error looking up %q: %w", key, err)
		}
		if !found {
			return nil, fmt.Errorf("resolver: key %q not found in any backend", key)
		}
		result[key] = val
	}

	return result, nil
}

// lookup searches all backends for the given key, returning the first match.
func (r *Resolver) lookup(key string) (string, bool, error) {
	for _, b := range r.backends {
		val, err := b.Get(key)
		if err != nil {
			return "", false, fmt.Errorf("backend %s: %w", b, err)
		}
		if val != "" {
			return val, true, nil
		}
	}
	return "", false, nil
}
