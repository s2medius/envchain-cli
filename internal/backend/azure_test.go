package backend

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

type mockAzureClient struct {
	secrets map[string]string
	err     error
}

func (m *mockAzureClient) GetSecret(_ context.Context, name, _ string, _ *azsecrets.GetSecretOptions) (azsecrets.GetSecretResponse, error) {
	if m.err != nil {
		return azsecrets.GetSecretResponse{}, m.err
	}
	val, ok := m.secrets[name]
	if !ok {
		return azsecrets.GetSecretResponse{}, errors.New("SecretNotFound: 404")
	}
	return azsecrets.GetSecretResponse{
		Secret: azsecrets.Secret{Value: &val},
	}, nil
}

func newTestAzureBackend(secrets map[string]string, err error) *AzureBackend {
	return &AzureBackend{
		client:   &mockAzureClient{secrets: secrets, err: err},
		vaultURL: "https://test-vault.vault.azure.net",
	}
}

func TestNewAzureBackend_MissingVaultURL(t *testing.T) {
	_, err := NewAzureBackend(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing vault_url")
	}
}

func TestAzureBackend_Get_Found(t *testing.T) {
	b := newTestAzureBackend(map[string]string{"MY_SECRET": "supersecret"}, nil)
	val, err := b.Get(context.Background(), "MY_SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected 'supersecret', got %q", val)
	}
}

func TestAzureBackend_Get_NotFound(t *testing.T) {
	b := newTestAzureBackend(map[string]string{}, nil)
	_, err := b.Get(context.Background(), "MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing secret")
	}
}

func TestAzureBackend_Get_ClientError(t *testing.T) {
	b := newTestAzureBackend(nil, errors.New("connection refused"))
	_, err := b.Get(context.Background(), "ANY_KEY")
	if err == nil {
		t.Fatal("expected error on client failure")
	}
}

func TestAzureBackend_String(t *testing.T) {
	b := newTestAzureBackend(nil, nil)
	s := b.String()
	if s != "AzureKeyVault(https://test-vault.vault.azure.net)" {
		t.Errorf("unexpected String() output: %q", s)
	}
}
