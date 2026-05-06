package backend_test

import (
	"errors"
	"testing"

	"github.com/envchain-cli/envchain-cli/internal/backend"
)

func TestNew_EnvBackend(t *testing.T) {
	p, err := backend.New(backend.TypeEnv, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p == nil {
		t.Fatal("expected provider, got nil")
	}
	if p.Name() != "env" {
		t.Errorf("expected name 'env', got %q", p.Name())
	}
}

func TestNew_FileBackend_MissingPath(t *testing.T) {
	_, err := backend.New(backend.TypeFile, map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing path, got nil")
	}
}

func TestNew_FileBackend_WithPath(t *testing.T) {
	p, err := backend.New(backend.TypeFile, map[string]string{"path": "/tmp/secrets.env"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Name() != "file" {
		t.Errorf("expected name 'file', got %q", p.Name())
	}
}

func TestNew_UnsupportedBackend(t *testing.T) {
	_, err := backend.New("unknown", nil)
	if err == nil {
		t.Fatal("expected error for unsupported backend, got nil")
	}
	var unsupported *backend.ErrUnsupportedBackend
	if !errors.As(err, &unsupported) {
		t.Errorf("expected ErrUnsupportedBackend, got %T", err)
	}
}

func TestErrSecretNotFound_Message(t *testing.T) {
	err := &backend.ErrSecretNotFound{Backend: "env", Key: "MY_SECRET"}
	expected := `secret "MY_SECRET" not found in backend "env"`
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestErrUnsupportedBackend_Message(t *testing.T) {
	err := &backend.ErrUnsupportedBackend{Type: "consul"}
	expected := `unsupported backend type: "consul"`
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}
