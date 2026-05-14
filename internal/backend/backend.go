package backend

import (
	"fmt"
)

// Backend is the interface implemented by all secret backends.
type Backend interface {
	Get(key string) (string, error)
	String() string
}

// New constructs a Backend from a type name and configuration map.
func New(backendType string, cfg map[string]string) (Backend, error) {
	switch backendType {
	case "env":
		return NewEnvBackend(cfg)
	case "file":
		return NewFileBackend(cfg)
	case "vault":
		return NewVaultBackend(cfg)
	case "ssm":
		return NewSSMBackend(cfg)
	case "secretsmanager":
		return NewSecretsManagerBackend(cfg)
	case "gcp":
		return NewGCPBackend(cfg)
	case "azure":
		return NewAzureBackend(cfg)
	case "1password":
		return NewOnePasswordBackend(cfg)
	case "onepassword_connect":
		return NewOnePasswordConnectBackend(cfg)
	case "doppler":
		return NewDopplerBackend(cfg)
	case "infisical":
		return NewInfisicalBackend(cfg)
	case "github":
		return NewGitHubBackend(cfg)
	case "keychain":
		return NewKeychainBackend(cfg)
	case "bitwarden":
		return NewBitwardenBackend(cfg)
	case "lastpass":
		return NewLastPassBackend(cfg)
	case "hashicorp":
		return NewHashiCorpBackend(cfg)
	case "akeyless":
		return NewAkeylessBackend(cfg)
	case "conjur":
		return NewConjurBackend(cfg)
	default:
		return nil, fmt.Errorf("unsupported backend type: %q", backendType)
	}
}
