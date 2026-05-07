package cmd

import (
	"fmt"
	"os"

	"github.com/envchain-cli/envchain-cli/internal/config"
	"github.com/envchain-cli/envchain-cli/internal/resolver"
	"github.com/envchain-cli/envchain-cli/internal/runner"
	"github.com/spf13/cobra"
)

var (
	configFile string
	keys       []string
)

var runCmd = &cobra.Command{
	Use:   "run -- <command> [args...]",
	Short: "Inject secrets and run a command",
	Long: `Resolve secrets from configured backends and inject them as
environment variables into the given command process.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		res, err := resolver.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to build resolver: %w", err)
		}

		r := runner.New(res)
		if err := r.Run(args[0], args[1:], keys); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	runCmd.Flags().StringVarP(&configFile, "config", "c", "envchain.yaml", "path to envchain config file")
	runCmd.Flags().StringArrayVarP(&keys, "key", "k", nil, "secret keys to resolve (default: all configured keys)")
	RootCmd.AddCommand(runCmd)
}
