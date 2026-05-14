package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestKeeperBackend(t *testing.T, handler http.HandlerFunc) (*KeeperBackend, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	b := &KeeperBackend{
		baseURL: srv.URL,
		token:   "test-token",
		client:  srv.Client(),
	}
	return b, srv
}

func TestNewKeeperBackend_MissingToken(t *testing.T) {
	_, err := NewKeeperBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewKeeperBackend_DefaultBaseURL(t *testing.T) {
	b, err := NewKeeperBackend(map[string]string{"token": "tok"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.baseURL != "https://keepersecurity.com/api/rest/sm/v1" {
		t.Errorf("unexpected base URL: %s", b.baseURL)
	}
}

func TestNewKeeperBackend_CustomBaseURL(t *testing.T) {
	b, err := NewKeeperBackend(map[string]string{"token": "tok", "base_url": "https://custom.example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.baseURL != "https://custom.example.com" {
		t.Errorf("unexpected base URL: %s", b.baseURL)
	}
}

func TestKeeperBackend_Get_Found(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		payload := keeperSecretResponse{}
		payload.Data.Fields = []struct {
			Type  string   `json:"type"`
			Value []string `json:"value"`
		}{{Type: "password", Value: []string{"s3cr3t"}}}
		_ = json.NewEncoder(w).Encode(payload)
	}
	b, srv := newTestKeeperBackend(t, handler)
	defer srv.Close()

	val, err := b.Get("record123/password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Errorf("expected 's3cr3t', got %q", val)
	}
}

func TestKeeperBackend_Get_NotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
	b, srv := newTestKeeperBackend(t, handler)
	defer srv.Close()

	_, err := b.Get("missing/password")
	if err == nil {
		t.Fatal("expected error for missing record")
	}
}

func TestKeeperBackend_Get_InvalidKeyFormat(t *testing.T) {
	b := &KeeperBackend{baseURL: "http://x", token: "t", client: &http.Client{}}
	_, err := b.Get("noslash")
	if err == nil {
		t.Fatal("expected error for invalid key format")
	}
}

func TestKeeperBackend_String(t *testing.T) {
	b := &KeeperBackend{baseURL: "https://keepersecurity.com/api/rest/sm/v1"}
	if b.String() != "keeper(https://keepersecurity.com/api/rest/sm/v1)" {
		t.Errorf("unexpected String(): %s", b.String())
	}
}
