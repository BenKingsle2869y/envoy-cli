package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/store"
)

var getCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Get the value of an environment variable from the store",
	Args:  cobra.ExactArgs(1),
	RunE:  runGet,
}

func init() {
	getCmd.Flags().StringP("passphrase", "p", "", "Passphrase to decrypt the store")
	getCmd.Flags().StringP("file", "f", store.DefaultStorePath, "Path to the encrypted store file")
	rootCmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) error {
	passphrase, err := crypto.ResolvePassphrase(cmd)
	if err != nil {
		return fmt.Errorf("passphrase error: %w", err)
	}

	filePath, _ := cmd.Flags().GetString("file")

	s, err := store.Load(filePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	key := args[0]
	value, ok := s.Env[key]
	if !ok {
		return fmt.Errorf("key %q not found in store", key)
	}

	fmt.Println(value)
	return nil
}
