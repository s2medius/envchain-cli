package backend

import "fmt"

// Backend is the interface all secret backends must implement.
type Backend interface {
	Get(key string) (string, error)
	String() string
}

// New constructs a Backend from the given type name and options map.
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
	case "1password-connect":
		return NewOnePasswordConnectBackend(opts)
	case "pulumi":
		return NewPulumiBackend(opts)
	case "nomad":
		return NewNomadBackend(opts)
	case "keeper":
		return NewKeeperBackend(opts)
	case "vercel":
		return NewVercelBackend(opts)
	case "netlify":
		return NewNetlifyBackend(opts)
	case "railway":
		return NewRailwayBackend(opts)
	case "flyio":
		return NewFlyIOBackend(opts)
	default:
		return nil, fmt.Errorf("backend: unsupported type %q", backendType)
	}
}
