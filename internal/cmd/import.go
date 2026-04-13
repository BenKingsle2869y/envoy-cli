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

	imported, skipped := applyParsedVars(s.Env, parsed, importOverwrite)

	if err := store.Save(store.DefaultStorePath(), s, passphrase); err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Imported %d variable(s), skipped %d existing.\n", imported, skipped)
	return nil
}

// applyParsedVars merges parsed key-value pairs into dst.
// If overwrite is false, existing keys are left unchanged and counted as skipped.
// Returns the number of imported and skipped entries.
func applyParsedVars(dst map[string]string, src map[string]string, overwrite bool) (imported, skipped int) {
	for k, v := range src {
		if _, exists := dst[k]; exists && !overwrite {
			skipped++
			continue
		}
		dst[k] = v
		imported++
	}
	return imported, skipped
}
