package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/store"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Re-encrypt the local store with a new passphrase",
	RunE:  runRotate,
}

func init() {
	rotateCmd.Flags().String("store", store.DefaultStorePath, "Path to the encrypted store file")
	rotateCmd.Flags().String("new-passphrase", "", "New passphrase to re-encrypt the store (or set ENVOY_NEW_PASSPHRASE)")
	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, _ []string) error {
	storePath, _ := cmd.Flags().GetString("store")
	newPassFlag, _ := cmd.Flags().GetString("new-passphrase")

	// Resolve current passphrase
	currentPass, err := crypto.ResolvePassphrase("", storePath)
	if err != nil {
		return fmt.Errorf("could not resolve current passphrase: %w", err)
	}

	// Load store with current passphrase
	s, err := store.Load(storePath, currentPass)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	// Resolve new passphrase
	newPass, err := crypto.ResolvePassphrase(newPassFlag, "")
	if err != nil {
		// Fall back to ENVOY_NEW_PASSPHRASE env var handled inside ResolvePassphrase;
		// if still not found, surface the error.
		return fmt.Errorf("could not resolve new passphrase: %w", err)
	}

	if newPass == currentPass {
		return fmt.Errorf("new passphrase must differ from the current passphrase")
	}

	// Save store with new passphrase
	if err := store.Save(storePath, s, newPass); err != nil {
		return fmt.Errorf("failed to save rotated store: %w", err)
	}

	fmt.Println("Passphrase rotated successfully.")
	return nil
}
