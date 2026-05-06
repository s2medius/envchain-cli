package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func vaultTestServer(t *testing.T, status int, data map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(status)
		if status == http.StatusOK {
			body, _ := json.Marshal(map[string]interface{}{
				"data": map[string]interface{}{"data": data},
			})
			w.Write(body)
		}
	}))
}

func TestNewVaultBackend_MissingAddress(t *testing.T) {
	_, err := NewVaultBackend(map[string]string{"token": "t", "path": "p"})
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewVaultBackend_MissingToken(t *testing.T) {
	_, err := NewVaultBackend(map[string]string{"address": "http://localhost", "path": "p"})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewVaultBackend_MissingPath(t *testing.T) {
	_, err := NewVaultBackend(map[string]string{"address": "http://localhost", "token": "t"})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestVaultBackend_Get_Found(t *testing.T) {
	srv := vaultTestServer(t, http.StatusOK, map[string]string{"DB_PASS": "secret"})
	defer srv.Close()
	v, _ := NewVaultBackend(map[string]string{"address": srv.URL, "token": "tok", "path": "myapp"})
	val, err := v.Get("DB_PASS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "secret" {
		t.Errorf("expected 'secret', got %q", val)
	}
}

func TestVaultBackend_Get_NotFound(t *testing.T) {
	srv := vaultTestServer(t, http.StatusOK, map[string]string{})
	defer srv.Close()
	v, _ := NewVaultBackend(map[string]string{"address": srv.URL, "token": "tok", "path": "myapp"})
	_, err := v.Get("MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestVaultBackend_List(t *testing.T) {
	data := map[string]string{"A": "1", "B": "2"}
	srv := vaultTestServer(t, http.StatusOK, data)
	defer srv.Close()
	v, _ := NewVaultBackend(map[string]string{"address": srv.URL, "token": "tok", "path": "myapp"})
	result, err := v.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(result))
	}
}

func TestVaultBackend_String(t *testing.T) {
	v, _ := NewVaultBackend(map[string]string{"address": "http://vault:8200", "token": "tok", "path": "myapp"})
	s := v.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
