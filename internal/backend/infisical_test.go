package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockInfisicalClient struct {
	secrets map[string]string
}

func (m *mockInfisicalClient) GetSecret(_, _, key string) (string, error) {
	val, ok := m.secrets[key]
	if !ok {
		return "", fmt.Errorf("secret %q not found", key)
	}
	return val, nil
}

func newTestInfisicalBackend(secrets map[string]string) *InfisicalBackend {
	return &InfisicalBackend{
		client:      &mockInfisicalClient{secrets: secrets},
		projectID:   "proj-123",
		environment: "dev",
	}
}

func TestNewInfisicalBackend_MissingToken(t *testing.T) {
	_, err := NewInfisicalBackend(map[string]string{"project_id": "p"})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewInfisicalBackend_MissingProjectID(t *testing.T) {
	_, err := NewInfisicalBackend(map[string]string{"token": "t"})
	if err == nil {
		t.Fatal("expected error for missing project_id")
	}
}

func TestNewInfisicalBackend_Defaults(t *testing.T) {
	b, err := NewInfisicalBackend(map[string]string{"token": "tok", "project_id": "proj"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.environment != "dev" {
		t.Errorf("expected default environment 'dev', got %q", b.environment)
	}
}

func TestInfisicalBackend_Get_Found(t *testing.T) {
	b := newTestInfisicalBackend(map[string]string{"MY_SECRET": "supersecret"})
	val, err := b.Get("MY_SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected 'supersecret', got %q", val)
	}
}

func TestInfisicalBackend_Get_NotFound(t *testing.T) {
	b := newTestInfisicalBackend(map[string]string{})
	_, err := b.Get("MISSING")
	if err == nil {
		t.Fatal("expected error for missing secret")
	}
}

func TestInfisicalBackend_String(t *testing.T) {
	b := newTestInfisicalBackend(nil)
	s := b.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}

func TestInfisicalHTTPClient_GetSecret_Found(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer mytoken" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		payload, _ := json.Marshal(map[string]interface{}{
			"secret": map[string]string{"secretValue": "hello"},
		})
		w.Write(payload)
	}))
	defer ts.Close()

	c := &infisicalHTTPClient{
		baseURL:    ts.URL,
		token:      "mytoken",
		httpClient: &http.Client{},
	}
	val, err := c.GetSecret("proj", "dev", "KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hello" {
		t.Errorf("expected 'hello', got %q", val)
	}
}
