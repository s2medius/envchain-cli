package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestOPConnectBackend(t *testing.T, handler http.HandlerFunc) (*OnePasswordConnectBackend, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	b := &OnePasswordConnectBackend{
		baseURL: ts.URL,
		token:   "test-token",
		vaultID: "vault-abc",
		client:  ts.Client(),
	}
	return b, ts
}

func TestNewOnePasswordConnectBackend_MissingToken(t *testing.T) {
	_, err := NewOnePasswordConnectBackend(map[string]string{
		"url":      "http://localhost:8080",
		"vault_id": "v1",
	})
	if err == nil || !strings.Contains(err.Error(), "token") {
		t.Fatalf("expected token error, got %v", err)
	}
}

func TestNewOnePasswordConnectBackend_MissingURL(t *testing.T) {
	_, err := NewOnePasswordConnectBackend(map[string]string{
		"token":    "tok",
		"vault_id": "v1",
	})
	if err == nil || !strings.Contains(err.Error(), "url") {
		t.Fatalf("expected url error, got %v", err)
	}
}

func TestNewOnePasswordConnectBackend_MissingVaultID(t *testing.T) {
	_, err := NewOnePasswordConnectBackend(map[string]string{
		"token": "tok",
		"url":   "http://localhost:8080",
	})
	if err == nil || !strings.Contains(err.Error(), "vault_id") {
		t.Fatalf("expected vault_id error, got %v", err)
	}
}

func TestOnePasswordConnectBackend_Get_Found(t *testing.T) {
	items := []opConnectItem{{
		ID:    "item1",
		Title: "myapp",
		Fields: []struct {
			Label string `json:"label"`
			Value string `json:"value"`
		}{{Label: "API_KEY", Value: "supersecret"}},
	}}
	b, ts := newTestOPConnectBackend(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})
	defer ts.Close()

	val, err := b.Get("myapp/API_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected 'supersecret', got %q", val)
	}
}

func TestOnePasswordConnectBackend_Get_NotFound(t *testing.T) {
	b, ts := newTestOPConnectBackend(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]opConnectItem{})
	})
	defer ts.Close()

	_, err := b.Get("myapp/API_KEY")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestOnePasswordConnectBackend_Get_InvalidKeyFormat(t *testing.T) {
	b, ts := newTestOPConnectBackend(t, func(w http.ResponseWriter, r *http.Request) {})
	defer ts.Close()

	_, err := b.Get("invalidkey")
	if err == nil || !strings.Contains(err.Error(), "format") {
		t.Fatalf("expected format error, got %v", err)
	}
}

func TestOnePasswordConnectBackend_String(t *testing.T) {
	b := &OnePasswordConnectBackend{baseURL: "http://localhost:8080", vaultID: "vault-abc"}
	s := b.String()
	if !strings.Contains(s, "onepassword_connect") || !strings.Contains(s, "vault-abc") {
		t.Errorf("unexpected String() output: %s", s)
	}
}
