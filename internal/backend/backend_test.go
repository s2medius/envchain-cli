package backend

import (
	"strings"
	"testing"
)

func TestNew_EnvBackend(t *testing.T) {
	b, err := New("env", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestNew_FileBackend_MissingPath(t *testing.T) {
	_, err := New("file", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestNew_FileBackend_WithPath(t *testing.T) {
	b, err := New("file", map[string]string{"path": "/dev/null"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestNew_VaultBackend(t *testing.T) {
	_, err := New("vault", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing vault config")
	}
}

func TestNew_UnsupportedBackend(t *testing.T) {
	_, err := New("nonexistent", map[string]string{})
	if err == nil {
		t.Fatal("expected error for unsupported backend")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected 'unsupported' in error, got: %v", err)
	}
}

func TestNew_OnePasswordConnectBackend(t *testing.T) {
	_, err := New("onepassword_connect", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing onepassword_connect config")
	}
	if !strings.Contains(err.Error(), "token") {
		t.Errorf("expected 'token' in error, got: %v", err)
	}
}

func TestNew_OnePasswordConnectBackend_Valid(t *testing.T) {
	b, err := New("onepassword_connect", map[string]string{
		"token":    "tok",
		"url":      "http://localhost:8080",
		"vault_id": "vault-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
	if !strings.Contains(b.String(), "onepassword_connect") {
		t.Errorf("unexpected String(): %s", b.String())
	}
}
