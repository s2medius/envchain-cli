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
	b, err := New("file", map[string]string{"path": "/tmp/test.env"})
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
	_, err := New("unknown", map[string]string{})
	if err == nil {
		t.Fatal("expected error for unsupported backend")
	}
}

func TestNew_OnePasswordBackend(t *testing.T) {
	b, err := New("1password", map[string]string{"vault": "dev"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
	if b.String() != "1password(vault=dev)" {
		t.Errorf("unexpected String(): %q", b.String())
	}
}

func TestNew_OnePasswordBackend_NoVault(t *testing.T) {
	b, err := New("1password", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.String() != "1password" {
		t.Errorf("unexpected String(): %q", b.String())
	}
}
