package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestNomadBackend(t *testing.T, handler http.HandlerFunc) (*NomadBackend, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	b := &NomadBackend{
		address: ts.URL,
		token:   "test-token",
		path:    "project/app",
		client:  ts.Client(),
	}
	return b, ts
}

func TestNewNomadBackend_MissingAddress(t *testing.T) {
	_, err := NewNomadBackend(map[string]string{"token": "t", "path": "p"})
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewNomadBackend_MissingToken(t *testing.T) {
	_, err := NewNomadBackend(map[string]string{"address": "http://localhost", "path": "p"})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewNomadBackend_MissingPath(t *testing.T) {
	_, err := NewNomadBackend(map[string]string{"address": "http://localhost", "token": "t"})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestNomadBackend_Get_Found(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Nomad-Token") != "test-token" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		payload := map[string]interface{}{
			"Items": map[string]string{"DB_PASSWORD": "s3cr3t"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}
	b, ts := newTestNomadBackend(t, handler)
	defer ts.Close()

	val, err := b.Get("DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Errorf("expected 's3cr3t', got %q", val)
	}
}

func TestNomadBackend_Get_NotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"Items": map[string]string{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}
	b, ts := newTestNomadBackend(t, handler)
	defer ts.Close()

	_, err := b.Get("MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestNomadBackend_Get_404(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
	b, ts := newTestNomadBackend(t, handler)
	defer ts.Close()

	_, err := b.Get("ANY")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestNomadBackend_String(t *testing.T) {
	b := &NomadBackend{path: "project/app"}
	if s := b.String(); s != "nomad(project/app)" {
		t.Errorf("unexpected String(): %q", s)
	}
}
