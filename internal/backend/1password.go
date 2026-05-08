package backend

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// OnePasswordBackend retrieves secrets from 1Password via the `op` CLI.
type OnePasswordBackend struct {
	vault   string
	account string
}

type opRunner interface {
	Run(item, field, vault, account string) (string, error)
}

type defaultOPRunner struct{}

func (r *defaultOPRunner) Run(item, field, vault, account string) (string, error) {
	args := []string{"item", "get", item, "--field", field, "--format", "json"}
	if vault != "" {
		args = append(args, "--vault", vault)
	}
	if account != "" {
		args = append(args, "--account", account)
	}
	out, err := exec.Command("op", args...).Output()
	if err != nil {
		return "", fmt.Errorf("op cli error: %w", err)
	}
	// op returns a JSON object; extract the "value" field simply
	raw := strings.TrimSpace(string(out))
	const valueKey = `"value":`
	idx := strings.Index(raw, valueKey)
	if idx == -1 {
		return "", fmt.Errorf("unexpected op output format")
	}
	rest := strings.TrimSpace(raw[idx+len(valueKey):])
	rest = strings.Trim(rest, `"\n,}`)
	return rest, nil
}

// NewOnePasswordBackend creates a new 1Password backend.
// opts keys: "vault" (optional), "account" (optional).
func NewOnePasswordBackend(opts map[string]string) (*OnePasswordBackend, error) {
	return &OnePasswordBackend{
		vault:   opts["vault"],
		account: opts["account"],
	}, nil
}

// Get retrieves the value for key formatted as "item/field".
func (b *OnePasswordBackend) Get(_ context.Context, key string) (string, error) {
	parts := strings.SplitN(key, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("1password: key must be in 'item/field' format, got %q", key)
	}
	r := &defaultOPRunner{}
	return r.Run(parts[0], parts[1], b.vault, b.account)
}

func (b *OnePasswordBackend) String() string {
	if b.vault != "" {
		return fmt.Sprintf("1password(vault=%s)", b.vault)
	}
	return "1password"
}
