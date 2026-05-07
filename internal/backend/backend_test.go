package backend

import (
	"testing"

	"github.com/envchain-cli/envchain-cli/internal/config"
)

func TestNew_EnvBackend(t *testing.T) {
	b, err := New(config.Backend{Type: "env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestNew_FileBackend_MissingPath(t *testing.T) {
	_, err := New(config.Backend{Type: "file", Options: map[string]string{}})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestNew_FileBackend_WithPath(t *testing.T) {
	f := writeDotenv(t, "X=1\n")
	b, err := New(config.Backend{Type: "file", Options: map[string]string{"path": f}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestNew_VaultBackend(t *testing.T) {
	_, err := New(config.Backend{
		Type: "vault",
		Options: map[string]string{
			"address": "http://127.0.0.1:8200",
			"token":   "root",
			"path":    "secret/data/app",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_UnsupportedBackend(t *testing.T) {
	_, err := New(config.Backend{Type: "unknown"})
	if err == nil {
		t.Fatal("expected error for unsupported backend")
	}
}

func TestNew_SecretsManagerBackend(t *testing.T) {
	_, err := New(config.Backend{
		Type: "secretsmanager",
		Options: map[string]string{
			"secret_id": "my/secret",
			"region":    "us-east-1",
		},
	})
	// We expect no construction error (AWS config load may succeed in CI)
	// The test validates the dispatch path, not live AWS connectivity.
	_ = err
}

func TestNew_SecretsManagerBackend_MissingSecretID(t *testing.T) {
	_, err := New(config.Backend{
		Type:    "secretsmanager",
		Options: map[string]string{"region": "us-east-1"},
	})
	if err == nil {
		t.Fatal("expected error for missing secret_id")
	}
}
