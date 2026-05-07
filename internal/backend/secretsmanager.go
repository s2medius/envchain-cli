package backend

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type secretsManagerClient interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// SecretsManagerBackend resolves variables from AWS Secrets Manager.
type SecretsManagerBackend struct {
	client secretsManagerClient
	secretID string
	region   string
}

// NewSecretsManagerBackend creates a new SecretsManagerBackend.
func NewSecretsManagerBackend(secretID, region string) (*SecretsManagerBackend, error) {
	if secretID == "" {
		return nil, fmt.Errorf("secretsmanager: secret_id is required")
	}
	if region == "" {
		return nil, fmt.Errorf("secretsmanager: region is required")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("secretsmanager: failed to load AWS config: %w", err)
	}

	client := secretsmanager.NewFromConfig(cfg)
	return &SecretsManagerBackend{
		client:   client,
		secretID: secretID,
		region:   region,
	}, nil
}

// Get retrieves the value for the given key from the secret JSON.
func (b *SecretsManagerBackend) Get(key string) (string, error) {
	out, err := b.client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(b.secretID),
	})
	if err != nil {
		return "", fmt.Errorf("secretsmanager: failed to get secret %q: %w", b.secretID, err)
	}

	if out.SecretString == nil {
		return "", fmt.Errorf("secretsmanager: secret %q has no string value", b.secretID)
	}

	pairs := parseKeyValueSecret(*out.SecretString)
	val, ok := pairs[key]
	if !ok {
		return "", fmt.Errorf("secretsmanager: key %q not found in secret %q", key, b.secretID)
	}
	return val, nil
}

// String returns a human-readable description of the backend.
func (b *SecretsManagerBackend) String() string {
	return fmt.Sprintf("secretsmanager(secret=%s, region=%s)", b.secretID, b.region)
}

// parseKeyValueSecret parses a simple KEY=VALUE newline-delimited secret string.
func parseKeyValueSecret(s string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}
