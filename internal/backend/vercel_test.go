package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestVercelBackend(t *testing.T, handler http.Handler) *VercelBackend {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	b, err := NewVercelBackend(map[string]string{
		"token":      "test-token",
		"project_id": "my-project",
		"api_url":    srv.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return b
}

func TestNewVercelBackend_MissingToken(t *testing.T) {
	_, err := NewVercelBackend(map[string]string{"project_id": "p"})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewVercelBackend_MissingProjectID(t *testing.T) {
	_, err := NewVercelBackend(map[string]string{"token": "tok"})
	if err == nil {
		t.Fatal("expected error for missing project_id")
	}
}

func TestNewVercelBackend_DefaultAPIURL(t *testing.T) {
	b, err := NewVercelBackend(map[string]string{"token": "tok", "project_id": "p"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.apiURL != defaultVercelAPIURL {
		t.Errorf("expected default api_url %q, got %q", defaultVercelAPIURL, b.apiURL)
	}
}

func TestVercelBackend_Get_Found(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"value": "secret-value"})
	})
	b := newTestVercelBackend(t, handler)
	val, err := b.Get("MY_SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "secret-value" {
		t.Errorf("expected %q, got %q", "secret-value", val)
	}
}

func TestVercelBackend_Get_NotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	b := newTestVercelBackend(t, handler)
	_, err := b.Get("MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestVercelBackend_Get_ClientError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	})
	b := newTestVercelBackend(t, handler)
	_, err := b.Get("KEY")
	if err == nil {
		t.Fatal("expected error on server error")
	}
}

func TestVercelBackend_String(t *testing.T) {
	b := newTestVercelBackend(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	if b.String() != "vercel(project=my-project)" {
		t.Errorf("unexpected String(): %s", b.String())
	}
}
