package backend

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

// mockGCPClient implements gcpSecretClient for testing.
type mockGCPClient struct {
	data map[string]string
	err  error
}

func (m *mockGCPClient) AccessSecretVersion(_ context.Context, req *secretmanagerpb.AccessSecretVersionRequest, _ ...option.ClientOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	val, ok := m.data[req.Name]
	if !ok {
		return nil, errors.New("secret not found")
	}
	return &secretmanagerpb.AccessSecretVersionResponse{
		Payload: &secretmanagerpb.SecretPayload{Data: []byte(val)},
	}, nil
}

func (m *mockGCPClient) Close() error { return nil }

func newTestGCPBackend(project string, data map[string]string) *GCPBackend {
	return &GCPBackend{
		client:  &mockGCPClient{data: data},
		project: project,
	}
}

func TestGCPBackend_Get_Found(t *testing.T) {
	key := "MY_SECRET"
	project := "my-project"
	fullName := "projects/" + project + "/secrets/" + key + "/versions/latest"
	b := newTestGCPBackend(project, map[string]string{fullName: "supersecret"})
	val, err := b.Get(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected 'supersecret', got %q", val)
	}
}

func TestGCPBackend_Get_NotFound(t *testing.T) {
	b := newTestGCPBackend("proj", map[string]string{})
	_, err := b.Get("MISSING")
	if err == nil {
		t.Fatal("expected error for missing secret")
	}
}

func TestGCPBackend_Get_ClientError(t *testing.T) {
	b := &GCPBackend{
		client:  &mockGCPClient{err: errors.New("rpc error")},
		project: "proj",
	}
	_, err := b.Get("KEY")
	if err == nil {
		t.Fatal("expected error from client")
	}
}

func TestGCPBackend_String(t *testing.T) {
	b := newTestGCPBackend("my-project", nil)
	if s := b.String(); s != "gcp(project=my-project)" {
		t.Errorf("unexpected String(): %q", s)
	}
}

func TestNewGCPBackend_MissingProject(t *testing.T) {
	_, err := NewGCPBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}
