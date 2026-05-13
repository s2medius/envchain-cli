package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestLastPassBackend(t *testing.T, handler http.HandlerFunc) *LastPassBackend {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	b, err := NewLastPassBackend(map[string]string{
		"username": "user@example.com",
		"api_key":  "test-api-key",
		"api_url":  server.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return b
}

func TestNewLastPassBackend_MissingUsername(t *testing.T) {
	_, err := NewLastPassBackend(map[string]string{"api_key": "key"})
	if err == nil {
		t.Fatal("expected error for missing username")
	}
}

func TestNewLastPassBackend_MissingAPIKey(t *testing.T) {
	_, err := NewLastPassBackend(map[string]string{"username": "user@example.com"})
	if err == nil {
		t.Fatal("expected error for missing api_key")
	}
}

func TestNewLastPassBackend_DefaultAPIURL(t *testing.T) {
	b, err := NewLastPassBackend(map[string]string{
		"username": "user@example.com",
		"api_key":  "key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.apiURL != "https://lastpass.com/enterpriseapi.php" {
		t.Errorf("expected default api_url, got %q", b.apiURL)
	}
}

func TestLastPassBackend_Get_Found(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		secrets := []lastpassSecret{{Name: "DB_PASSWORD", Value: "secret123"}}
		json.NewEncoder(w).Encode(secrets)
	}
	b := newTestLastPassBackend(t, handler)
	val, err := b.Get("shared-folder/DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "secret123" {
		t.Errorf("expected %q, got %q", "secret123", val)
	}
}

func TestLastPassBackend_Get_NotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]lastpassSecret{})
	}
	b := newTestLastPassBackend(t, handler)
	_, err := b.Get("shared-folder/MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestLastPassBackend_Get_InvalidKeyFormat(t *testing.T) {
	b := newTestLastPassBackend(t, func(w http.ResponseWriter, r *http.Request) {})
	_, err := b.Get("no-slash-key")
	if err == nil {
		t.Fatal("expected error for invalid key format")
	}
}

func TestLastPassBackend_Get_ServerError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	b := newTestLastPassBackend(t, handler)
	_, err := b.Get("folder/KEY")
	if err == nil {
		t.Fatal("expected error for server error")
	}
}

func TestLastPassBackend_String(t *testing.T) {
	b := newTestLastPassBackend(t, func(w http.ResponseWriter, r *http.Request) {})
	s := b.String()
	if s != "lastpass(username=user@example.com)" {
		t.Errorf("unexpected String(): %q", s)
	}
}
