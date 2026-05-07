package backend

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

type azureSecretsClient interface {
	GetSecret(ctx context.Context, name, version string, opts *azsecrets.GetSecretOptions) (azsecrets.GetSecretResponse, error)
}

// AzureBackend retrieves secrets from Azure Key Vault.
type AzureBackend struct {
	client  azureSecretsClient
	vaultURL string
}

// NewAzureBackend creates an AzureBackend from config options.
// Required options: vault_url
func NewAzureBackend(opts map[string]string) (*AzureBackend, error) {
	vaultURL, ok := opts["vault_url"]
	if !ok || vaultURL == "" {
		return nil, fmt.Errorf("azure backend: missing required option 'vault_url'")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("azure backend: failed to create credential: %w", err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("azure backend: failed to create client: %w", err)
	}

	return &AzureBackend{client: client, vaultURL: vaultURL}, nil
}

// Get retrieves a secret value by key from Azure Key Vault.
func (b *AzureBackend) Get(ctx context.Context, key string) (string, error) {
	resp, err := b.client.GetSecret(ctx, key, "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "SecretNotFound") || strings.Contains(err.Error(), "404") {
			return "", fmt.Errorf("azure backend: secret %q not found", key)
		}
		return "", fmt.Errorf("azure backend: failed to get secret %q: %w", key, err)
	}

	if resp.Value == nil {
		return "", fmt.Errorf("azure backend: secret %q has nil value", key)
	}

	return *resp.Value, nil
}

// String returns a human-readable description of the backend.
func (b *AzureBackend) String() string {
	return fmt.Sprintf("AzureKeyVault(%s)", b.vaultURL)
}
