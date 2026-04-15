// vaultpeek is a CLI tool for browsing and diffing HashiCorp Vault secret paths
// across environments.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is set at build time via ldflags.
var version = "dev"

// rootCmd is the base command for the vaultpeek CLI.
var rootCmd = &cobra.Command{
	Use:   "vaultpeek",
	Short: "Browse and diff HashiCorp Vault secret paths across environments",
	Long: `vaultpeek is a CLI tool for inspecting and comparing Vault KV secrets.

It supports browsing secret paths, reading secret values, and diffing secrets
between two Vault environments or paths.

Configuration is read from environment variables (VAULT_ADDR, VAULT_TOKEN)
or provided via flags.`,
	SilenceUsage: true,
}

// versionCmd prints the current build version.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of vaultpeek",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vaultpeek %s\n", version)
	},
}

// Global flags shared across subcommands.
var (
	vaultAddr  string
	vaultToken string
	mountPath  string
)

func init() {
	// Persistent flags are available to all subcommands.
	rootCmd.PersistentFlags().StringVar(
		&vaultAddr,
		"vault-addr",
		"",
		"Vault server address (overrides VAULT_ADDR env var)",
	)
	rootCmd.PersistentFlags().StringVar(
		&vaultToken,
		"vault-token",
		"",
		"Vault token (overrides VAULT_TOKEN env var)",
	)
	rootCmd.PersistentFlags().StringVar(
		&mountPath,
		"mount",
		"secret",
		"KV v2 mount path",
	)

	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
