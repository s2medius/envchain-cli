//go:build darwin

package backend

import (
	"errors"
	"fmt"
	"testing"
)

func newTestKeychainBackend(service string, secrets map[string]string) *KeychainBackend {
	return &KeychainBackend{
		service: service,
		execCommand: func(name string, args ...string) ([]byte, error) {
			// args: find-generic-password -s <service> -a <account> -w
			if len(args) < 5 {
				return nil, errors.New("unexpected args")
			}
			account := args[4]
			if val, ok := secrets[account]; ok {
				return []byte(val + "\n"), nil
			}
			return nil, fmt.Errorf("security: item not found")
		},
	}
}

func TestNewKeychainBackend_MissingService(t *testing.T) {
	_, err := NewKeychainBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing service, got nil")
	}
}

func TestNewKeychainBackend_Valid(t *testing.T) {
	b, err := NewKeychainBackend(map[string]string{"service": "my-app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.service != "my-app" {
		t.Errorf("expected service 'my-app', got %q", b.service)
	}
}

func TestKeychainBackend_Get_Found(t *testing.T) {
	b := newTestKeychainBackend("my-app", map[string]string{
		"DB_PASSWORD": "s3cr3t",
	})
	val, err := b.Get("DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Errorf("expected 's3cr3t', got %q", val)
	}
}

func TestKeychainBackend_Get_NotFound(t *testing.T) {
	b := newTestKeychainBackend("my-app", map[string]string{})
	_, err := b.Get("MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestKeychainBackend_Get_MultipleKeys(t *testing.T) {
	secrets := map[string]string{
		"API_KEY":     "abc123",
		"DB_PASSWORD": "s3cr3t",
		"TOKEN":       "tok-xyz",
	}
	b := newTestKeychainBackend("my-app", secrets)
	for key, expected := range secrets {
		val, err := b.Get(key)
		if err != nil {
			t.Errorf("unexpected error for key %q: %v", key, err)
			continue
		}
		if val != expected {
			t.Errorf("key %q: expected %q, got %q", key, expected, val)
		}
	}
}

func TestKeychainBackend_String(t *testing.T) {
	b := newTestKeychainBackend("my-app", nil)
	expected := "keychain(service=my-app)"
	if s := b.String(); s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
