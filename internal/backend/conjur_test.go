package backend

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestConjurBackend(t *testing.T, handler http.HandlerFunc) (*ConjurBackend, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	b := &ConjurBackend{
		applianceURL: ts.URL,
		account:      "myorg",
		token:        "test-token",
		client:       ts.Client(),
	}
	return b, ts
}

func TestNewConjurBackend_MissingAddress(t *testing.T) {
	_, err := NewConjurBackend(map[string]string{"account": "myorg", "token": "tok"})
	if err == nil || !strings.Contains(err.Error(), "address") {
		t.Fatalf("expected address error, got %v", err)
	}
}

func TestNewConjurBackend_MissingAccount(t *testing.T) {
	_, err := NewConjurBackend(map[string]string{"address": "http://conjur", "token": "tok"})
	if err == nil || !strings.Contains(err.Error(), "account") {
		t.Fatalf("expected account error, got %v", err)
	}
}

func TestNewConjurBackend_MissingToken(t *testing.T) {
	_, err := NewConjurBackend(map[string]string{"address": "http://conjur", "account": "myorg"})
	if err == nil || !strings.Contains(err.Error(), "token") {
		t.Fatalf("expected token error, got %v", err)
	}
}

func TestConjurBackend_Get_Found(t *testing.T) {
	b, ts := newTestConjurBackend(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "myapp/db/password") {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "supersecret")
	})
	defer ts.Close()

	val, err := b.Get("myapp/db/password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected 'supersecret', got %q", val)
	}
}

func TestConjurBackend_Get_NotFound(t *testing.T) {
	b, ts := newTestConjurBackend(t, func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	defer ts.Close()

	_, err := b.Get("missing/key")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestConjurBackend_Get_ClientError(t *testing.T) {
	b, ts := newTestConjurBackend(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer ts.Close()

	_, err := b.Get("some/key")
	if err == nil || !strings.Contains(err.Error(), "unexpected status") {
		t.Fatalf("expected status error, got %v", err)
	}
}

func TestConjurBackend_String(t *testing.T) {
	b := &ConjurBackend{applianceURL: "https://conjur.example.com", account: "myorg"}
	s := b.String()
	if !strings.Contains(s, "conjur") || !strings.Contains(s, "myorg") {
		t.Errorf("unexpected String() output: %s", s)
	}
}
