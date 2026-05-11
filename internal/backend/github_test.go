package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockGitHubClient struct {
	secrets map[string]string
	err     error
}

func (m *mockGitHubClient) GetSecret(_ context.Context, _, _, secretName string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	val, ok := m.secrets[secretName]
	if !ok {
		return "", fmt.Errorf("secret %q not found", secretName)
	}
	return val, nil
}

func newTestGitHubBackend(secrets map[string]string, err error) *GitHubBackend {
	return &GitHubBackend{
		client: &mockGitHubClient{secrets: secrets, err: err},
		owner:  "myorg",
		repo:   "myrepo",
	}
}

func TestNewGitHubBackend_MissingToken(t *testing.T) {
	_, err := NewGitHubBackend(map[string]string{"owner": "o", "repo": "r"})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewGitHubBackend_MissingOwner(t *testing.T) {
	_, err := NewGitHubBackend(map[string]string{"token": "t", "repo": "r"})
	if err == nil {
		t.Fatal("expected error for missing owner")
	}
}

func TestNewGitHubBackend_MissingRepo(t *testing.T) {
	_, err := NewGitHubBackend(map[string]string{"token": "t", "owner": "o"})
	if err == nil {
		t.Fatal("expected error for missing repo")
	}
}

func TestGitHubBackend_Get_Found(t *testing.T) {
	b := newTestGitHubBackend(map[string]string{"MY_SECRET": "supersecret"}, nil)
	val, err := b.Get(context.Background(), "MY_SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected 'supersecret', got %q", val)
	}
}

func TestGitHubBackend_Get_NotFound(t *testing.T) {
	b := newTestGitHubBackend(map[string]string{}, nil)
	_, err := b.Get(context.Background(), "MISSING")
	if err == nil {
		t.Fatal("expected error for missing secret")
	}
}

func TestGitHubBackend_Get_ClientError(t *testing.T) {
	b := newTestGitHubBackend(nil, fmt.Errorf("network error"))
	_, err := b.Get(context.Background(), "ANY")
	if err == nil {
		t.Fatal("expected error from client")
	}
}

func TestGitHubBackend_String(t *testing.T) {
	b := newTestGitHubBackend(nil, nil)
	if b.String() != "github(myorg/myrepo)" {
		t.Errorf("unexpected string: %s", b.String())
	}
}

func TestGitHubBackend_HTTPClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"value": "from-http"})
	}))
	defer ts.Close()

	b, err := NewGitHubBackend(map[string]string{
		"token":    "mytoken",
		"owner":    "org",
		"repo":     "repo",
		"base_url": ts.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, err := b.Get(context.Background(), "SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "from-http" {
		t.Errorf("expected 'from-http', got %q", val)
	}
}
