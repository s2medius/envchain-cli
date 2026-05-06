package backend

import (
	"os"
	"testing"
)

func writeDotenv(t *testing.T, content string) string {
	t.Helper()
	tmp, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := tmp.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	tmp.Close()
	return tmp.Name()
}

func TestFileBackend_Get_Found(t *testing.T) {
	path := writeDotenv(t, "DB_PASS=hunter2\nAPI_KEY=\"abc123\"\n")
	b := NewFileBackend(path)

	val, err := b.Get("DB_PASS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hunter2" {
		t.Errorf("expected 'hunter2', got %q", val)
	}

	val, err = b.Get("API_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "abc123" {
		t.Errorf("expected 'abc123' (quotes stripped), got %q", val)
	}
}

func TestFileBackend_Get_NotFound(t *testing.T) {
	path := writeDotenv(t, "FOO=bar\n")
	b := NewFileBackend(path)

	_, err := b.Get("MISSING")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(ErrSecretNotFound); !ok {
		t.Errorf("expected ErrSecretNotFound, got %T", err)
	}
}

func TestFileBackend_IgnoresComments(t *testing.T) {
	path := writeDotenv(t, "# this is a comment\nKEY=value\n")
	b := NewFileBackend(path)

	keys, err := b.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 1 || keys[0] != "KEY" {
		t.Errorf("expected [KEY], got %v", keys)
	}
}

func TestFileBackend_MissingFile(t *testing.T) {
	b := NewFileBackend("/nonexistent/path/.env")
	_, err := b.Get("ANY")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFileBackend_String(t *testing.T) {
	b := NewFileBackend("/etc/app.env")
	if b.String() != "file(/etc/app.env)" {
		t.Errorf("unexpected String(): %s", b.String())
	}
}

func TestFileBackend_CachesOnSecondCall(t *testing.T) {
	path := writeDotenv(t, "ONCE=loaded\n")
	b := NewFileBackend(path)

	// First call loads from disk
	_, _ = b.Get("ONCE")
	// Remove the file — second call must still succeed via cache
	os.Remove(path)

	val, err := b.Get("ONCE")
	if err != nil {
		t.Fatalf("expected cached value, got error: %v", err)
	}
	if val != "loaded" {
		t.Errorf("expected 'loaded', got %q", val)
	}
}
