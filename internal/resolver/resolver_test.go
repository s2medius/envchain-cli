package resolver_test

import (
	"testing"

	"github.com/yourusername/envchain-cli/internal/config"
	"github.com/yourusername/envchain-cli/internal/resolver"
)

func makeConfig(types ...string) *config.Config {
	var backends []config.BackendConfig
	for _, t := range types {
		bc := config.BackendConfig{Type: t}
		if t == "env" {
			bc.Prefix = "TESTAPP_"
		}
		backends = append(backends, bc)
	}
	return &config.Config{Version: 1, Backends: backends}
}

func TestNew_ValidConfig(t *testing.T) {
	cfg := makeConfig("env")
	r, err := resolver.New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestNew_InvalidBackend(t *testing.T) {
	cfg := makeConfig("unsupported")
	_, err := resolver.New(cfg)
	if err == nil {
		t.Fatal("expected error for unsupported backend type")
	}
}

func TestResolve_KeyFoundInEnv(t *testing.T) {
	t.Setenv("TESTAPP_SECRET_KEY", "supersecret")

	cfg := &config.Config{
		Version: 1,
		Backends: []config.BackendConfig{
			{Type: "env", Prefix: "TESTAPP_"},
		},
	}
	r, err := resolver.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := r.Resolve([]string{"SECRET_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["SECRET_KEY"] != "supersecret" {
		t.Errorf("expected %q, got %q", "supersecret", result["SECRET_KEY"])
	}
}

func TestResolve_KeyNotFound(t *testing.T) {
	cfg := &config.Config{
		Version: 1,
		Backends: []config.BackendConfig{
			{Type: "env", Prefix: "TESTAPP_"},
		},
	}
	r, err := resolver.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = r.Resolve([]string{"NONEXISTENT_VAR"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}
