package backend

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// mockSSMClient implements ssmClient for testing.
type mockSSMClient struct {
	params map[string]string
	err    error
}

func (m *mockSSMClient) GetParameter(_ context.Context, in *ssm.GetParameterInput, _ ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	val, ok := m.params[aws.ToString(in.Name)]
	if !ok {
		return nil, &types.ParameterNotFound{}
	}
	return &ssm.GetParameterOutput{
		Parameter: &types.Parameter{Value: aws.String(val)},
	}, nil
}

func newTestSSMBackend(path string, params map[string]string) *SSMBackend {
	return &SSMBackend{
		client: &mockSSMClient{params: params},
		path:   path,
	}
}

func TestSSMBackend_Get_Found(t *testing.T) {
	b := newTestSSMBackend("/myapp/prod", map[string]string{
		"/myapp/prod/DB_PASSWORD": "s3cr3t",
	})
	val, err := b.Get("DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Errorf("expected 's3cr3t', got %q", val)
	}
}

func TestSSMBackend_Get_NotFound(t *testing.T) {
	b := newTestSSMBackend("/myapp/prod", map[string]string{})
	_, err := b.Get("MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestSSMBackend_Get_ClientError(t *testing.T) {
	b := &SSMBackend{
		client: &mockSSMClient{err: errors.New("network failure")},
		path:   "/myapp",
	}
	_, err := b.Get("SOME_KEY")
	if err == nil {
		t.Fatal("expected error from client, got nil")
	}
}

func TestNewSSMBackend_MissingPath(t *testing.T) {
	_, err := NewSSMBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing path, got nil")
	}
}

func TestSSMBackend_String(t *testing.T) {
	b := newTestSSMBackend("/myapp/prod", nil)
	if got := b.String(); got != "ssm(/myapp/prod)" {
		t.Errorf("unexpected String(): %q", got)
	}
}
