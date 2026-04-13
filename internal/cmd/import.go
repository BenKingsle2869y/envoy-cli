package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/env"
	"envoy-cli/internal/store"
)

var importOverwrite bool

func init() {
	importCmd := &cobra.Command{
		Use:   "import [file]",
		Short: "Import variables from a .env file into the store",
		Args:  cobra.ExactArgs(1),
		RunE:  runImport,
	}

	importCmd.Flags().BoolVar(&importOverwrite, "overwrite", false, "Overwrite existing keys")
	RootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %q", filePath)
	}

	passphrase, err := crypto.ResolvePassphrase()
	if err != nil {
		return fmt.Errorf("passphrase error: %w", err)
	}

	s, err := store.Load(store.DefaultStorePath(), passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	parsed, err := env.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse file %q: %w", filePath, err)
	}

	if len(parsed) == 0 {
		fmt.Fprintln(os.Stdout, "No variables found in file.")
		return nil
	}

	imported := 0
	skipped := 0
	for k, v := range parsed {
		if _, exists := s.Env[k]; exists && !importOverwrite {
			skipped++
			continue
		}
		s.Env[k] = v
		imported++
	}

	if err := store.Save(store.DefaultStorePath(), s, passphrase); err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Imported %d variable(s), skipped %d existing.\n", imported, skipped)
	return nil
}
