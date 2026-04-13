package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envoy",
	Short:oy — manage and sync .env files with",
	L `envoy is a lightweight CLI for managing and syncingross local and remote environments with encryption support.`,
}

 Execute runs the root command and error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
