package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestFlyIOBackend(t *testing.T, handler http.HandlerFunc) *FlyIOBackend {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	b, err := NewFlyIOBackend(map[string]string{
		"token":   "test-token",
		"app_id":  "my-app",
		"api_url": srv.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return b
}

func TestNewFlyIOBackend_MissingToken(t *testing.T) {
	_, err := NewFlyIOBackend(map[string]string{"app_id": "my-app"})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewFlyIOBackend_MissingAppID(t *testing.T) {
	_, err := NewFlyIOBackend(map[string]string{"token": "tok"})
	if err == nil {
		t.Fatal("expected error for missing app_id")
	}
}

func TestNewFlyIOBackend_DefaultAPIURL(t *testing.T) {
	b, err := NewFlyIOBackend(map[string]string{"token": "tok", "app_id": "app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.apiURL != flyioDefaultAPIURL {
		t.Errorf("expected default API URL %q, got %q", flyioDefaultAPIURL, b.apiURL)
	}
}

func TestFlyIOBackend_Get_Found(t *testing.T) {
	b := newTestFlyIOBackend(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"MY_SECRET": "hello"})
	})
	val, err := b.Get("MY_SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hello" {
		t.Errorf("expected %q, got %q", "hello", val)
	}
}

func TestFlyIOBackend_Get_NotFound(t *testing.T) {
	b := newTestFlyIOBackend(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"OTHER": "val"})
	})
	_, err := b.Get("MY_SECRET")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestFlyIOBackend_Get_ClientError(t *testing.T) {
	b := newTestFlyIOBackend(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	_, err := b.Get("KEY")
	if err == nil {
		t.Fatal("expected error on non-200 status")
	}
}

func TestFlyIOBackend_String(t *testing.T) {
	b := newTestFlyIOBackend(t, nil)
	if s := b.String(); s != "flyio(app=my-app)" {
		t.Errorf("unexpected String(): %q", s)
	}
}
