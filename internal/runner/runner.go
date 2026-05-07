package runner

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/envchain-cli/envchain-cli/internal/resolver"
)

// Runner injects resolved environment variables into a child process.
type Runner struct {
	resolver *resolver.Resolver
}

// New creates a new Runner backed by the given Resolver.
func New(r *resolver.Resolver) *Runner {
	return &Runner{resolver: r}
}

// Run executes the given command with the resolved environment variables
// merged on top of the current process environment.
func (r *Runner) Run(name string, args []string, keys []string) error {
	if name == "" {
		return fmt.Errorf("runner: command name must not be empty")
	}

	envMap, err := r.ResolveAll(keys)
	if err != nil {
		return fmt.Errorf("runner: failed to resolve keys: %w", err)
	}

	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start with the current process environment.
	cmd.Env = os.Environ()

	// Append resolved secrets, overriding any existing values.
	for k, v := range envMap {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runner: command %q exited with error: %w", name, err)
	}
	return nil
}
