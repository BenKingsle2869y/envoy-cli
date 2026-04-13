package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/store"
)

var setCmd = &cobra.Command{
	Use:   "set KEY=VALUE [KEY=VALUE...]",
	Short: "Set one or more environment variables in the store",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSet,
}

func init() {
	setCmd.Flags().StringP("passphrase", "p", "", "Passphrase to decrypt/encrypt the store")
	setCmd.Flags().StringP("file", "f", store.DefaultStorePath, "Path to the encrypted store file")
	rootCmd.AddCommand(setCmd)
}

func runSet(cmd *cobra.Command, args []string) error {
	passphrase, err := crypto.ResolvePassphrase(cmd)
	if err != nil {
		return fmt.Errorf("passphrase error: %w", err)
	}

	filePath, _ := cmd.Flags().GetString("file")

	s, err := store.Load(filePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid format %q: expected KEY=VALUE", arg)
		}
		key, value := strings.TrimSpace(parts[0]), parts[1]
		if key == "" {
			return fmt.Errorf("key must not be empty in %q", arg)
		}
		s.Env[key] = value
		fmt.Printf("Set %s\n", key)
	}

	if err := store.Save(filePath, passphrase, s); err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}
	return nil
}
