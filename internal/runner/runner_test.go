package runner_test

import (
	"testing"

	"github.com/envchain-cli/envchain-cli/internal/config"
	"github.com/envchain-cli/envchain-cli/internal/resolver"
	"github.com/envchain-cli/envchain-cli/internal/runner"
)

func makeResolver(t *testing.T) *resolver.Resolver {
	t.Helper()
	cfg := &config.Config{
		Version: 1,
		Backends: []config.Backend{
			{Name: "env", Type: "env"},
		},
	}
	r, err := resolver.New(cfg)
	if err != nil {
		t.Fatalf("makeResolver: %v", err)
	}
	return r
}

func TestRun_EmptyCommand(t *testing.T) {
	r := runner.New(makeResolver(t))
	err := r.Run("", nil, nil)
	if err == nil {
		t.Fatal("expected error for empty command, got nil")
	}
}

func TestRun_SimpleCommand(t *testing.T) {
	r := runner.New(makeResolver(t))
	// "true" exits 0 on all POSIX systems.
	if err := r.Run("true", nil, nil); err != nil {
		t.Fatalf("unexpected error running 'true': %v", err)
	}
}

func TestRun_CommandWithArgs(t *testing.T) {
	r := runner.New(makeResolver(t))
	// echo with args should succeed.
	if err := r.Run("echo", []string{"hello", "world"}, nil); err != nil {
		t.Fatalf("unexpected error running 'echo': %v", err)
	}
}

func TestRun_CommandFailure(t *testing.T) {
	r := runner.New(makeResolver(t))
	// "false" always exits non-zero.
	err := r.Run("false", nil, nil)
	if err == nil {
		t.Fatal("expected error from 'false', got nil")
	}
}

func TestRun_NonExistentCommand(t *testing.T) {
	r := runner.New(makeResolver(t))
	err := r.Run("__nonexistent_cmd__", nil, nil)
	if err == nil {
		t.Fatal("expected error for non-existent command, got nil")
	}
}
