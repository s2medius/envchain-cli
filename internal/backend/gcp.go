package backend

import (
	"context"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

// GCPBackend retrieves secrets from Google Cloud Secret Manager.
type GCPBackend struct {
	client  gcpSecretClient
	project string
}

type gcpSecretClient interface {
	AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...option.ClientOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
	Close() error
}

type realGCPClient struct {
	inner *secretmanager.Client
}

func (r *realGCPClient) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...option.ClientOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	return r.inner.AccessSecretVersion(ctx, req)
}

func (r *realGCPClient) Close() error {
	return r.inner.Close()
}

// NewGCPBackend creates a GCPBackend from config options.
// Required options: "project"
func NewGCPBackend(opts map[string]string) (*GCPBackend, error) {
	project, ok := opts["project"]
	if !ok || project == "" {
		return nil, fmt.Errorf("gcp backend: missing required option 'project'")
	}
	ctx := context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("gcp backend: failed to create client: %w", err)
	}
	return &GCPBackend{client: &realGCPClient{inner: c}, project: project}, nil
}

// Get retrieves the latest version of a secret by name.
func (g *GCPBackend) Get(key string) (string, error) {
	return g.GetVersion(key, "latest")
}

// GetVersion retrieves a specific version of a secret by name.
// Use "latest" to retrieve the most recent enabled version.
func (g *GCPBackend) GetVersion(key, version string) (string, error) {
	if version == "" {
		version = "latest"
	}
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", g.project, key, version)
	resp, err := g.client.AccessSecretVersion(context.Background(), &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return "", fmt.Errorf("gcp backend: failed to access secret %q (version %s): %w", key, version, err)
	}
	return strings.TrimRight(string(resp.Payload.Data), "\n"), nil
}

func (g *GCPBackend) String() string {
	return fmt.Sprintf("gcp(project=%s)", g.project)
}
