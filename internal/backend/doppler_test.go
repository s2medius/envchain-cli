package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestDopplerBackend(t *testing.T, handler http.HandlerFunc) (*DopplerBackend, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	b := &DopplerBackend{
		token:   "test-token",
		project: "myproject",
		config:  "production",
		client:  ts.Client(),
	}
	// Override the base URL by pointing requests to the test server.
	// We achieve this by wrapping the client to rewrite the host.
	b.client = &rewriteHostClient{base: ts.URL, inner: ts.Client()}
	return b, ts
}

// rewriteHostClient rewrites request URLs to point to a test server.
type rewriteHostClient struct {
	base  string
	inner dopplerClient
}

func (r *rewriteHostClient) Do(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = req.URL.Hostname() // keep host parsing
	// Replace the full URL host with test server host
	parsed, _ := http.NewRequest(req.Method, r.base+req.URL.RequestURI(), req.Body)
	parsed.Header = req.Header
	return r.inner.Do(parsed)
}

func TestNewDopplerBackend_MissingToken(t *testing.T) {
	_, err := NewDopplerBackend("", "proj", "cfg", nil)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewDopplerBackend_MissingProject(t *testing.T) {
	_, err := NewDopplerBackend("tok", "", "cfg", nil)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestNewDopplerBackend_MissingConfig(t *testing.T) {
	_, err := NewDopplerBackend("tok", "proj", "", nil)
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestDopplerBackend_Get_Found(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"secret": map[string]interface{}{
				"value": map[string]string{"raw": "supersecret"},
			},
		})
	}
	b, _ := newTestDopplerBackend(t, handler)

	val, err := b.Get("MY_SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected 'supersecret', got %q", val)
	}
}

func TestDopplerBackend_Get_NotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
	b, _ := newTestDopplerBackend(t, handler)

	_, err := b.Get("MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestDopplerBackend_Get_ClientError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}
	b, _ := newTestDopplerBackend(t, handler)

	_, err := b.Get("MY_SECRET")
	if err == nil {
		t.Fatal("expected error for server error")
	}
}

func TestDopplerBackend_String(t *testing.T) {
	b := &DopplerBackend{project: "myproject", config: "production"}
	got := b.String()
	expected := "doppler(project=myproject, config=production)"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
