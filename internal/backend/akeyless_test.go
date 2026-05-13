package backend

import (
	"context"
	"errors"
	"testing"
)

// mockAkeylessClient is a test double for AkeylessClient.
type mockAkeylessClient struct {
	secrets map[string]string
	err     error
}

func (m *mockAkeylessClient) GetSecretValue(_ context.Context, name, _ string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	v, ok := m.secrets[name]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func newTestAkeylessBackend(secrets map[string]string, err error) *AkeylessBackend {
	return &AkeylessBackend{
		client:  &mockAkeylessClient{secrets: secrets, err: err},
		token:   "test-token",
		gateway: "https://api.akeyless.io",
	}
}

func TestNewAkeylessBackend_MissingToken(t *testing.T) {
	_, err := NewAkeylessBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewAkeylessBackend_DefaultGateway(t *testing.T) {
	b, err := NewAkeylessBackend(map[string]string{"token": "tok"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.gateway != "https://api.akeyless.io" {
		t.Errorf("expected default gateway, got %q", b.gateway)
	}
}

func TestNewAkeylessBackend_CustomGateway(t *testing.T) {
	b, err := NewAkeylessBackend(map[string]string{"token": "tok", "gateway": "https://gw.example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.gateway != "https://gw.example.com" {
		t.Errorf("expected custom gateway, got %q", b.gateway)
	}
}

func TestAkeylessBackend_Get_Found(t *testing.T) {
	b := newTestAkeylessBackend(map[string]string{"/prod/db/password": "s3cr3t"}, nil)
	val, err := b.Get(context.Background(), "/prod/db/password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Errorf("expected %q, got %q", "s3cr3t", val)
	}
}

func TestAkeylessBackend_Get_NotFound(t *testing.T) {
	b := newTestAkeylessBackend(map[string]string{}, nil)
	_, err := b.Get(context.Background(), "/missing/key")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestAkeylessBackend_Get_ClientError(t *testing.T) {
	b := newTestAkeylessBackend(nil, errors.New("network error"))
	_, err := b.Get(context.Background(), "/any/key")
	if err == nil {
		t.Fatal("expected error from client failure")
	}
}

func TestAkeylessBackend_String(t *testing.T) {
	b := newTestAkeylessBackend(nil, nil)
	s := b.String()
	if s != "akeyless(gateway=https://api.akeyless.io)" {
		t.Errorf("unexpected String(): %q", s)
	}
}
