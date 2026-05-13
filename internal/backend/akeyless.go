package backend

import (
	"context"
	"fmt"
	"net/http"

	akeyless "github.com/akeylesslabs/akeyless-go/v3"
)

// AkeylessClient defines the interface used for fetching secrets.
type AkeylessClient interface	{
	GetSecretValue(ctx context.Context, name, token string) (string, error)
}

type akeylessHTTPClient struct {
	api    *akeyless.APIClient
	gateway string
}

func (c *akeylessHTTPClient) GetSecretValue(ctx context.Context, name, token string) (string, error) {
	body := akeyless.GetSecretValue{
		Names: []string{name},
		Token: &token,
	}
	res, _, err := c.api.V2Api.GetSecretValue(ctx).Body(body).Execute()
	if err != nil {
		return "", fmt.Errorf("akeyless: get secret %q: %w", name, err)
	}
	val, ok := res[name]
	if !ok {
		return "", fmt.Errorf("akeyless: secret %q not found in response", name)
	}
	return val, nil
}

// AkeylessBackend retrieves secrets from Akeyless.
type AkeylessBackend struct {
	client  AkeylessClient
	token   string
	gateway string
}

// NewAkeylessBackend creates a new AkeylessBackend from the given options.
func NewAkeylessBackend(opts map[string]string) (*AkeylessBackend, error) {
	token, ok := opts["token"]
	if !ok || token == "" {
		return nil, fmt.Errorf("akeyless: missing required option 'token'")
	}
	gateway := opts["gateway"]
	if gateway == "" {
		gateway = "https://api.akeyless.io"
	}
	cfg := akeyless.NewConfiguration()
	cfg.Servers = akeyless.ServerConfigurations{
		{URL: gateway},
	}
	cfg.HTTPClient = &http.Client{}
	client := &akeylessHTTPClient{
		api:     akeyless.NewAPIClient(cfg),
		gateway: gateway,
	}
	return &AkeylessBackend{client: client, token: token, gateway: gateway}, nil
}

// Get retrieves the value of the named secret.
func (b *AkeylessBackend) Get(ctx context.Context, key string) (string, error) {
	return b.client.GetSecretValue(ctx, key, b.token)
}

// String returns a human-readable description of the backend.
func (b *AkeylessBackend) String() string {
	return fmt.Sprintf("akeyless(gateway=%s)", b.gateway)
}
