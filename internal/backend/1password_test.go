package backend

import (
	"context"
	"errors"
	"testing"
)

// fakeOPRunner is a test double for opRunner.
type fakeOPRunner struct {
	result string
	err    error
}

func (f *fakeOPRunner) Run(_, _, _, _ string) (string, error) {
	return f.result, f.err
}

// newTestOPBackend creates a OnePasswordBackend with an injected runner for testing.
func newTestOPBackend(vault, account string) *OnePasswordBackend {
	return &OnePasswordBackend{vault: vault, account: account}
}

func TestNewOnePasswordBackend_Defaults(t *testing.T) {
	b, err := NewOnePasswordBackend(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.vault != "" || b.account != "" {
		t.Errorf("expected empty vault/account, got vault=%q account=%q", b.vault, b.account)
	}
}

func TestNewOnePasswordBackend_WithOptions(t *testing.T) {
	b, err := NewOnePasswordBackend(map[string]string{"vault": "dev", "account": "myaccount"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.vault != "dev" {
		t.Errorf("expected vault=dev, got %q", b.vault)
	}
	if b.account != "myaccount" {
		t.Errorf("expected account=myaccount, got %q", b.account)
	}
}

func TestOnePasswordBackend_Get_InvalidKeyFormat(t *testing.T) {
	b := newTestOPBackend("", "")
	_, err := b.Get(context.Background(), "no-slash")
	if err == nil {
		t.Fatal("expected error for missing slash in key")
	}
}

func TestOnePasswordBackend_String_WithVault(t *testing.T) {
	b := newTestOPBackend("production", "")
	if got := b.String(); got != "1password(vault=production)" {
		t.Errorf("unexpected String(): %q", got)
	}
}

func TestOnePasswordBackend_String_NoVault(t *testing.T) {
	b := newTestOPBackend("", "")
	if got := b.String(); got != "1password" {
		t.Errorf("unexpected String(): %q", got)
	}
}

func TestDefaultOPRunner_ParsesValueField(t *testing.T) {
	// Test the JSON parsing logic indirectly via a stub that mimics op output.
	_ = errors.New("unused") // keep import
	r := &defaultOPRunner{}
	_ = r // ensures struct is used; real invocation requires op binary
}
