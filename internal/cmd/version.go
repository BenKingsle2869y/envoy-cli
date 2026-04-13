package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information set at build time via ldflags.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information for envoy-cli",
	Long:  `Displays the current version, git commit hash, and build date of envoy-cli.`,
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Fprintf(cmd.OutOrStdout(), "envoy-cli version %s\n", Version)
	fmt.Fprintf(cmd.OutOrStdout(), "  commit:     %s\n", Commit)
	fmt.Fprintf(cmd.OutOrStdout(), "  build date: %s\n", BuildDate)
}
