package backend

import (
	"os"
	"testing"
)

func TestEnvBackend_Get_Found(t *testing.T) {
	t.Setenv("TEST_SECRET_KEY", "hello")

	b := NewEnvBackend("")
	val, err := b.Get("TEST_SECRET_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hello" {
		t.Errorf("expected 'hello', got %q", val)
	}
}

func TestEnvBackend_Get_NotFound(t *testing.T) {
	os.Unsetenv("ENVCHAIN_MISSING")

	b := NewEnvBackend("")
	_, err := b.Get("ENVCHAIN_MISSING")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(ErrSecretNotFound); !ok {
		t.Errorf("expected ErrSecretNotFound, got %T", err)
	}
}

func TestEnvBackend_Get_WithPrefix(t *testing.T) {
	t.Setenv("APP_DB_PASS", "secret")

	b := NewEnvBackend("APP_")
	val, err := b.Get("DB_PASS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "secret" {
		t.Errorf("expected 'secret', got %q", val)
	}
}

func TestEnvBackend_List_WithPrefix(t *testing.T) {
	t.Setenv("APP_FOO", "1")
	t.Setenv("APP_BAR", "2")
	t.Setenv("OTHER_VAR", "3")

	b := NewEnvBackend("APP_")
	keys, err := b.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := map[string]bool{}
	for _, k := range keys {
		found[k] = true
	}
	if !found["FOO"] || !found["BAR"] {
		t.Errorf("expected FOO and BAR in list, got %v", keys)
	}
	if found["OTHER_VAR"] {
		t.Errorf("OTHER_VAR should not appear in prefix-filtered list")
	}
}

func TestEnvBackend_String(t *testing.T) {
	if NewEnvBackend("").String() != "env" {
		t.Error("expected 'env'")
	}
	if NewEnvBackend("APP_").String() != "env(prefix=APP_)" {
		t.Error("expected 'env(prefix=APP_)'")
	}
}
