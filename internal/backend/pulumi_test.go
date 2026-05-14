package backend

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestPulumiBackend(t *testing.T, handler http.HandlerFunc) (*PulumiBackend, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	b, err := NewPulumiBackend(map[string]string{
		"token":       "test-token",
		"org":         "my-org",
		"environment": "prod",
		"api_url":     ts.URL,
	}, ts.Client())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return b, ts
}

func TestNewPulumiBackend_MissingToken(t *testing.T) {
	_, err := NewPulumiBackend(map[string]string{
		"org": "my-org", "environment": "prod",
	}, nil)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewPulumiBackend_MissingOrg(t *testing.T) {
	_, err := NewPulumiBackend(map[string]string{
		"token": "tok", "environment": "prod",
	}, nil)
	if err == nil {
		t.Fatal("expected error for missing org")
	}
}

func TestNewPulumiBackend_MissingEnvironment(t *testing.T) {
	_, err := NewPulumiBackend(map[string]string{
		"token": "tok", "org": "my-org",
	}, nil)
	if err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestNewPulumiBackend_DefaultAPIURL(t *testing.T) {
	b, err := NewPulumiBackend(map[string]string{
		"token": "tok", "org": "my-org", "environment": "prod",
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.apiURL != defaultPulumiAPIURL {
		t.Errorf("expected default api_url %q, got %q", defaultPulumiAPIURL, b.apiURL)
	}
}

func TestPulumiBackend_Get_Found(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "token test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		body, _ := json.Marshal(map[string]interface{}{
			"values": map[string]interface{}{"MY_SECRET": "supersecret"},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
	b, ts := newTestPulumiBackend(t, handler)
	defer ts.Close()

	val, err := b.Get(context.Background(), "MY_SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected %q, got %q", "supersecret", val)
	}
}

func TestPulumiBackend_Get_NotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := json.Marshal(map[string]interface{}{
			"values": map[string]interface{}{},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
	b, ts := newTestPulumiBackend(t, handler)
	defer ts.Close()

	_, err := b.Get(context.Background(), "MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestPulumiBackend_Get_ServerError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	b, ts := newTestPulumiBackend(t, handler)
	defer ts.Close()

	_, err := b.Get(context.Background(), "KEY")
	if err == nil {
		t.Fatal("expected error on server error")
	}
}

func TestPulumiBackend_String(t *testing.T) {
	b, ts := newTestPulumiBackend(t, func(w http.ResponseWriter, r *http.Request) {})
	defer ts.Close()
	got := b.String()
	expected := "pulumi(org=my-org, environment=prod)"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
