package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "envchain",
	Short: "Securely inject environment variables from secret backends into processes",
	Long: `envchain-cli loads secrets from one or more configured backends
(e.g. environment files, AWS SSM, HashiCorp Vault) and injects them
as environment variables before executing a subprocess.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config", "c",
		"envchain.json",
		"path to envchain config file",
	)
}
