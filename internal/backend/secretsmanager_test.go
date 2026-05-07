package backend

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type mockSMClient struct {
	secretString string
	err          error
}

func (m *mockSMClient) GetSecretValue(_ context.Context, _ *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(m.secretString),
	}, nil
}

func newTestSMBackend(secretString string, err error) *SecretsManagerBackend {
	return &SecretsManagerBackend{
		client:   &mockSMClient{secretString: secretString, err: err},
		secretID: "my/secret",
		region:   "us-east-1",
	}
}

func TestSecretsManagerBackend_Get_Found(t *testing.T) {
	b := newTestSMBackend("DB_PASS=hunter2\nAPI_KEY=abc123", nil)
	val, err := b.Get("DB_PASS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hunter2" {
		t.Errorf("expected hunter2, got %s", val)
	}
}

func TestSecretsManagerBackend_Get_NotFound(t *testing.T) {
	b := newTestSMBackend("OTHER=value", nil)
	_, err := b.Get("MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestSecretsManagerBackend_Get_ClientError(t *testing.T) {
	b := newTestSMBackend("", errors.New("access denied"))
	_, err := b.Get("ANY")
	if err == nil {
		t.Fatal("expected error from client")
	}
}

func TestSecretsManagerBackend_String(t *testing.T) {
	b := newTestSMBackend("", nil)
	s := b.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestNewSecretsManagerBackend_MissingSecretID(t *testing.T) {
	_, err := NewSecretsManagerBackend("", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing secret_id")
	}
}

func TestNewSecretsManagerBackend_MissingRegion(t *testing.T) {
	_, err := NewSecretsManagerBackend("my/secret", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestParseKeyValueSecret(t *testing.T) {
	input := "# comment\nFOO=bar\nBAZ=qux\n"
	pairs := parseKeyValueSecret(input)
	if pairs["FOO"] != "bar" {
		t.Errorf("expected bar, got %s", pairs["FOO"])
	}
	if _, ok := pairs["# comment"]; ok {
		t.Error("comments should be ignored")
	}
}
