package backend

import (
	"testing"
)

func TestNew_EnvBackend(t *testing.T) {
	b, err := New(BackendConfig{Type: "env", Options: map[string]string{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestNew_FileBackend_MissingPath(t *testing.T) {
	_, err := New(BackendConfig{Type: "file", Options: map[string]string{}})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestNew_FileBackend_WithPath(t *testing.T) {
	b, err := New(BackendConfig{Type: "file", Options: map[string]string{"path": "/tmp/test.env"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestNew_VaultBackend(t *testing.T) {
	b, err := New(BackendConfig{
		Type: "vault",
		Options: map[string]string{
			"address": "http://vault:8200",
			"token":   "mytoken",
			"path":    "myapp/config",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestNew_UnsupportedBackend(t *testing.T) {
	_, err := New(BackendConfig{Type: "unknown"})
	if err == nil {
		t.Fatal("expected error for unsupported backend")
	}
}

func TestErrSecretNotFound_Message(t *testing.T) {
	err := ErrSecretNotFound{Path: "MY_KEY"}
	expected := "secret not found: MY_KEY"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}
