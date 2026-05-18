package backend

import (
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
	b, err := New("file", map[string]string{"path": "/tmp/.env"})
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
		t.Fatal("expected error for missing vault options")
	}
}

func TestNew_UnsupportedBackend(t *testing.T) {
	_, err := New("nonexistent", map[string]string{})
	if err == nil {
		t.Fatal("expected error for unsupported backend type")
	}
}

func TestNew_FlyIOBackend_Missing(t *testing.T) {
	_, err := New("flyio", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing flyio options")
	}
}

func TestNew_FlyIOBackend_Valid(t *testing.T) {
	b, err := New("flyio", map[string]string{
		"token":  "tok",
		"app_id": "my-app",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}
