package backend

import "fmt"

// Backend is the interface implemented by all secret backends.
type Backend interface {
	Get(key string) (string, error)
	String() string
}

// New creates a Backend from a type name and options map.
func New(backendType string, opts map[string]string) (Backend, error) {
	switch backendType {
	case "env":
		return NewEnvBackend(opts)
	case "file":
		return NewFileBackend(opts)
	case "vault":
		return NewVaultBackend(opts)
	case "ssm":
		return NewSSMBackend(opts)
	case "secretsmanager":
		return NewSecretsManagerBackend(opts)
	case "gcp":
		return NewGCPBackend(opts)
	case "azure":
		return NewAzureBackend(opts)
	case "1password":
		return NewOnePasswordBackend(opts)
	case "doppler":
		return NewDopplerBackend(opts)
	case "infisical":
		return NewInfisicalBackend(opts)
	case "github":
		return NewGitHubBackend(opts)
	case "keychain":
		return NewKeychainBackend(opts)
	case "bitwarden":
		return NewBitwardenBackend(opts)
	case "lastpass":
		return NewLastPassBackend(opts)
	case "hashicorp":
		return NewHashiCorpBackend(opts)
	case "akeyless":
		return NewAkeylessBackend(opts)
	case "conjur":
		return NewConjurBackend(opts)
	default:
		return nil, fmt.Errorf("backend: unsupported type %q", backendType)
	}
}
