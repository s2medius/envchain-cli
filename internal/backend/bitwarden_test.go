package backend

import (
	"errors"
	"testing"
)

type mockBitwardenClient struct {
	fields map[string]string
	err    error
}

func (m *mockBitwardenClient) GetItem(_ string) (map[string]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.fields, nil
}

func newTestBitwardenBackend(fields map[string]string, err error) *BitwardenBackend {
	return &BitwardenBackend{
		client: &mockBitwardenClient{fields: fields, err: err},
		itemID: "test-item-id",
	}
}

func TestNewBitwardenBackend_MissingToken(t *testing.T) {
	_, err := NewBitwardenBackend(map[string]string{"item_id": "abc"})
	if err == nil {
		t.Fatal("expected error for missing access_token")
	}
}

func TestNewBitwardenBackend_MissingItemID(t *testing.T) {
	_, err := NewBitwardenBackend(map[string]string{"access_token": "tok"})
	if err == nil {
		t.Fatal("expected error for missing item_id")
	}
}

func TestNewBitwardenBackend_DefaultBaseURL(t *testing.T) {
	b, err := NewBitwardenBackend(map[string]string{
		"access_token": "tok",
		"item_id":      "abc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c, ok := b.client.(*bitwardenHTTPClient)
	if !ok {
		t.Fatal("expected *bitwardenHTTPClient")
	}
	if c.baseURL != "http://localhost:8087" {
		t.Errorf("expected default base_url, got %q", c.baseURL)
	}
}

func TestBitwardenBackend_Get_Found(t *testing.T) {
	b := newTestBitwardenBackend(map[string]string{"DB_PASSWORD": "secret123"}, nil)
	val, err := b.Get("DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "secret123" {
		t.Errorf("expected %q, got %q", "secret123", val)
	}
}

func TestBitwardenBackend_Get_NotFound(t *testing.T) {
	b := newTestBitwardenBackend(map[string]string{"OTHER": "val"}, nil)
	_, err := b.Get("MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestBitwardenBackend_Get_ItemNotFound(t *testing.T) {
	b := newTestBitwardenBackend(nil, nil)
	_, err := b.Get("ANY")
	if err == nil {
		t.Fatal("expected error when item not found")
	}
}

func TestBitwardenBackend_Get_ClientError(t *testing.T) {
	b := newTestBitwardenBackend(nil, errors.New("connection refused"))
	_, err := b.Get("KEY")
	if err == nil {
		t.Fatal("expected error on client failure")
	}
}

func TestBitwardenBackend_String(t *testing.T) {
	b := newTestBitwardenBackend(nil, nil)
	if s := b.String(); s != "bitwarden(item=test-item-id)" {
		t.Errorf("unexpected String(): %q", s)
	}
}
