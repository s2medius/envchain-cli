package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestNetlifyBackend(t *testing.T, handler http.HandlerFunc) *NetlifyBackend {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return &NetlifyBackend{
		token:  "test-token",
		siteID: "my-site",
		apiURL: server.URL,
		client: server.Client(),
	}
}

func TestNewNetlifyBackend_MissingToken(t *testing.T) {
	_, err := NewNetlifyBackend(map[string]string{"site_id": "abc"})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewNetlifyBackend_MissingSiteID(t *testing.T) {
	_, err := NewNetlifyBackend(map[string]string{"token": "tok"})
	if err == nil {
		t.Fatal("expected error for missing site_id")
	}
}

func TestNewNetlifyBackend_DefaultAPIURL(t *testing.T) {
	b, err := NewNetlifyBackend(map[string]string{"token": "tok", "site_id": "abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.apiURL != defaultNetlifyAPIURL {
		t.Errorf("expected default API URL, got %s", b.apiURL)
	}
}

func TestNetlifyBackend_Get_Found(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		payload := map[string]interface{}{
			"values": []map[string]string{
				{"value": "secret-value", "context": "all"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}
	b := newTestNetlifyBackend(t, handler)
	val, err := b.Get("MY_SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "secret-value" {
		t.Errorf("expected 'secret-value', got %q", val)
	}
}

func TestNetlifyBackend_Get_NotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
	b := newTestNetlifyBackend(t, handler)
	_, err := b.Get("MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestNetlifyBackend_Get_ClientError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	b := newTestNetlifyBackend(t, handler)
	_, err := b.Get("KEY")
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestNetlifyBackend_String(t *testing.T) {
	b := &NetlifyBackend{siteID: "my-site"}
	if b.String() != "netlify(site=my-site)" {
		t.Errorf("unexpected String() output: %s", b.String())
	}
}
