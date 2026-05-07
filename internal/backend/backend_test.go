package backend

import (
	"testing"

	"github.com/envchain-cli/envchain-cli/internal/config"
)

func TestNew_EnvBackend(t *testing.T) {
	b, err := New(config.Backend{Type: "env", Options: map[string]string{}})
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
	b, err := New(config.Backend{Type: "file", Options: map[string]string{"path": ".env"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestNew_VaultBackend(t *testing.T) {
	_, err := New(config.Backend{Type: "vault", Options: map[string]string{}})
	if err == nil {
		t.Fatal("expected error for missing vault options")
	}
}

func TestNew_UnsupportedBackend(t *testing.T) {
	_, err := New(config.Backend{Type: "unknown", Options: map[string]string{}})
	if err == nil {
		t.Fatal("expected error for unsupported backend")
	}
}

func TestNew_SSMBackend(t *testing.T) {
	_, err := New(config.Backend{Type: "ssm", Options: map[string]string{}})
	// SSM may succeed or fail depending on AWS env; we just check no panic
	_ = err
}

func TestNew_SecretsManagerBackend(t *testing.T) {
	_, err := New(config.Backend{Type: "secretsmanager", Options: map[string]string{}})
	_ = err
}

func TestNew_GCPBackend_MissingProject(t *testing.T) {
	_, err := New(config.Backend{Type: "gcp", Options: map[string]string{}})
	if err == nil {
		t.Fatal("expected error for missing gcp project")
	}
}
