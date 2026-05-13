package backend

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockHashiCorpClient struct {
	response *http.Response
	err      error
}

func (m *mockHashiCorpClient) Do(_ *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func newTestHashiCorpBackend(t *testing.T, client HashiCorpCloudClient) *HashiCorpBackend {
	t.Helper()
	b := &HashiCorpBackend{
		client:    client,
		baseURL:   "https://api.cloud.hashicorp.com",
		token:     "test-token",
		orgID:     "org-123",
		projectID: "proj-456",
		appName:   "my-app",
	}
	return b
}

func TestNewHashiCorpBackend_MissingToken(t *testing.T) {
	_, err := NewHashiCorpBackend(map[string]string{})
	if err == nil || !strings.Contains(err.Error(), "token") {
		t.Fatalf("expected token error, got %v", err)
	}
}

func TestNewHashiCorpBackend_MissingOrgID(t *testing.T) {
	_, err := NewHashiCorpBackend(map[string]string{"token": "t"})
	if err == nil || !strings.Contains(err.Error(), "org_id") {
		t.Fatalf("expected org_id error, got %v", err)
	}
}

func TestNewHashiCorpBackend_MissingProjectID(t *testing.T) {
	_, err := NewHashiCorpBackend(map[string]string{"token": "t", "org_id": "o"})
	if err == nil || !strings.Contains(err.Error(), "project_id") {
		t.Fatalf("expected project_id error, got %v", err)
	}
}

func TestNewHashiCorpBackend_MissingAppName(t *testing.T) {
	_, err := NewHashiCorpBackend(map[string]string{"token": "t", "org_id": "o", "project_id": "p"})
	if err == nil || !strings.Contains(err.Error(), "app_name") {
		t.Fatalf("expected app_name error, got %v", err)
	}
}

func TestHashiCorpBackend_Get_Found(t *testing.T) {
	body := `{"secret":{"static_version":{"value":"supersecret"}}}`
	client := &mockHashiCorpClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		},
	}
	b := newTestHashiCorpBackend(t, client)
	val, err := b.Get("MY_SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected 'supersecret', got %q", val)
	}
}

func TestHashiCorpBackend_Get_NotFound(t *testing.T) {
	client := &mockHashiCorpClient{
		response: &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader(`{}`)),
		},
	}
	b := newTestHashiCorpBackend(t, client)
	_, err := b.Get("MISSING")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestHashiCorpBackend_Get_ServerError(t *testing.T) {
	client := &mockHashiCorpClient{
		response: &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader(`{}`)),
		},
	}
	b := newTestHashiCorpBackend(t, client)
	_, err := b.Get("KEY")
	if err == nil || !strings.Contains(err.Error(), "unexpected status") {
		t.Fatalf("expected status error, got %v", err)
	}
}

func TestHashiCorpBackend_String(t *testing.T) {
	b := newTestHashiCorpBackend(t, nil)
	if s := b.String(); s != "hashicorp(app=my-app)" {
		t.Errorf("unexpected string: %q", s)
	}
}
