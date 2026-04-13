package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/store"
)

var unsetCmd = &cobra.Command{
	Use:   "unset KEY [KEY...]",
	Short: "Remove one or more environment variables from the store",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runUnset,
}

func init() {
	unsetCmd.Flags().StringP("passphrase", "p", "", "Passphrase to decrypt/encrypt the store")
	unsetCmd.Flags().StringP("file", "f", store.DefaultStorePath, "Path to the encrypted store file")
	rootCmd.AddCommand(unsetCmd)
}

func runUnset(cmd *cobra.Command, args []string) error {
	passphrase, err := crypto.ResolvePassphrase(cmd)
	if err != nil {
		return fmt.Errorf("passphrase error: %w", err)
	}

	filePath, _ := cmd.Flags().GetString("file")

	s, err := store.Load(filePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	for _, key := range args {
		if _, ok := s.Env[key]; !ok {
			fmt.Printf("Warning: key %q not found, skipping\n", key)
			continue
		}
		delete(s.Env, key)
		fmt.Printf("Unset %s\n", key)
	}

	if err := store.Save(filePath, passphrase, s); err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}
	return nil
}
