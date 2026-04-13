package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/store"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new envoy store in the current directory",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().StringP("passphrase", "p", "", "Passphrase to encrypt the store (or set ENVOY_PASSPHRASE)")
	initCmd.Flags().StringP("output", "o", "", "Path to write the store file (default: .envoy/store.enc)")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	passphrase, _ := cmd.Flags().GetString("passphrase")
	outputPath, _ := cmd.Flags().GetString("output")

	if outputPath == "" {
		outputPath = store.DefaultStorePath
	}

	passphrase, err := crypto.ResolvePassphrase(passphrase)
	if err != nil {
		return fmt.Errorf("passphrase required: use --passphrase or set ENVOY_PASSPHRASE: %w", err)
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create store directory %q: %w", dir, err)
	}

	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("store already exists at %q; remove it first or choose a different path", outputPath)
	}

	emptyEnv := map[string]string{}
	if err := store.Save(outputPath, passphrase, emptyEnv); err != nil {
		return fmt.Errorf("failed to initialise store: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Initialised empty envoy store at %s\n", outputPath)
	return nil
}
