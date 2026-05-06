package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/envchain/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "envchain.json")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	raw := `{"version":1,"backends":[{"type":"env","name":"local","params":{}}]}`
	p := writeTemp(t, raw)

	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Version != 1 {
		t.Errorf("expected version 1, got %d", cfg.Version)
	}
	if len(cfg.Backends) != 1 {
		t.Errorf("expected 1 backend, got %d", len(cfg.Backends))
	}
}

func TestLoad_MissingVersion(t *testing.T) {
	raw := `{"backends":[{"type":"env","name":"local"}]}`
	p := writeTemp(t, raw)

	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for missing version, got nil")
	}
}

func TestLoad_MissingBackendType(t *testing.T) {
	raw := `{"version":1,"backends":[{"name":"local"}]}`
	p := writeTemp(t, raw)

	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for missing backend type, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/envchain.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	p := writeTemp(t, `{not valid json}`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
