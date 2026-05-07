package backend

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// ssmClient defines the subset of SSM API used by SSMBackend.
type ssmClient interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

// SSMBackend resolves secrets from AWS Systems Manager Parameter Store.
type SSMBackend struct {
	client ssmClient
	path   string
}

// NewSSMBackend creates an SSMBackend. opts must contain "path" (parameter prefix)
// and optionally "region".
func NewSSMBackend(opts map[string]string) (*SSMBackend, error) {
	path, ok := opts["path"]
	if !ok || path == "" {
		return nil, fmt.Errorf("ssm backend: missing required option 'path'")
	}

	cfgOpts := []func(*config.LoadOptions) error{}
	if region, ok := opts["region"]; ok && region != "" {
		cfgOpts = append(cfgOpts, config.WithRegion(region))
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), cfgOpts...)
	if err != nil {
		return nil, fmt.Errorf("ssm backend: failed to load AWS config: %w", err)
	}

	return &SSMBackend{
		client: ssm.NewFromConfig(awsCfg),
		path:   strings.TrimSuffix(path, "/"),
	}, nil
}

// Get retrieves the parameter at <path>/<key> with decryption enabled.
func (s *SSMBackend) Get(key string) (string, error) {
	name := s.path + "/" + key
	out, err := s.client.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("ssm backend: get %q: %w", name, err)
	}
	if out.Parameter == nil || out.Parameter.Value == nil {
		return "", fmt.Errorf("ssm backend: parameter %q has no value", name)
	}
	return *out.Parameter.Value, nil
}

// String returns a human-readable description of the backend.
func (s *SSMBackend) String() string {
	return fmt.Sprintf("ssm(%s)", s.path)
}
