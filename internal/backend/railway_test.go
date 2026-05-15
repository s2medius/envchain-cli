package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestRailwayBackend(t *testing.T, handler http.HandlerFunc) *RailwayBackend {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	b, err := NewRailwayBackend(map[string]string{
		"token":          "test-token",
		"project_id":     "proj-123",
		"environment_id": "env-456",
		"api_url":        srv.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return b
}

func TestNewRailwayBackend_MissingToken(t *testing.T) {
	_, err := NewRailwayBackend(map[string]string{
		"project_id": "proj-123", "environment_id": "env-456",
	})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewRailwayBackend_MissingProjectID(t *testing.T) {
	_, err := NewRailwayBackend(map[string]string{
		"token": "tok", "environment_id": "env-456",
	})
	if err == nil {
		t.Fatal("expected error for missing project_id")
	}
}

func TestNewRailwayBackend_MissingEnvironmentID(t *testing.T) {
	_, err := NewRailwayBackend(map[string]string{
		"token": "tok", "project_id": "proj-123",
	})
	if err == nil {
		t.Fatal("expected error for missing environment_id")
	}
}

func TestNewRailwayBackend_DefaultAPIURL(t *testing.T) {
	b, err := NewRailwayBackend(map[string]string{
		"token": "tok", "project_id": "p", "environment_id": "e",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.apiURL != defaultRailwayAPIURL {
		t.Errorf("expected default API URL %q, got %q", defaultRailwayAPIURL, b.apiURL)
	}
}

func TestRailwayBackend_Get_Found(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"variables": map[string]interface{}{
					"edges": []map[string]interface{}{
						{"node": map[string]string{"name": "DB_PASSWORD", "value": "s3cr3t"}},
						{"node": map[string]string{"name": "API_KEY", "value": "abc123"}},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}
	b := newTestRailwayBackend(t, handler)
	val, err := b.Get("DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Errorf("expected %q, got %q", "s3cr3t", val)
	}
}

func TestRailwayBackend_Get_NotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"variables": map[string]interface{}{"edges": []interface{}{}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}
	b := newTestRailwayBackend(t, handler)
	_, err := b.Get("MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRailwayBackend_Get_ServerError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	b := newTestRailwayBackend(t, handler)
	_, err := b.Get("ANY")
	if err == nil {
		t.Fatal("expected error on server error")
	}
}

func TestRailwayBackend_String(t *testing.T) {
	b := &RailwayBackend{projectID: "proj-123", environmentID: "env-456"}
	s := b.String()
	if s != "railway(project=proj-123, environment=env-456)" {
		t.Errorf("unexpected String() output: %q", s)
	}
}
