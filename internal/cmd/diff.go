package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/env"
	"envoy-cli/internal/store"
)

var diffCmd = &cobra.Command{
	Use:   "diff <environment> <file>",
	Short: "Show differences between a stored environment and a local .env file",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	environment := args[0]
	filePath := args[1]

	passphrase, err := crypto.ResolvePassphrase()
	if err != nil {
		return fmt.Errorf("could not resolve passphrase: %w", err)
	}

	s, err := store.Load(store.DefaultStorePath(), passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	stored, ok := s.Envs[environment]
	if !ok {
		return fmt.Errorf("environment %q not found in store", environment)
	}

	local, err := env.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse local file: %w", err)
	}

	changes := env.Diff(stored, local)

	if len(changes) == 0 {
		fmt.Fprintln(os.Stdout, "No differences found.")
		return nil
	}

	for _, c := range changes {
		switch c.Type {
		case env.Added:
			fmt.Fprintf(os.Stdout, "+ %s=%s\n", c.Key, c.NewValue)
		case env.Removed:
			fmt.Fprintf(os.Stdout, "- %s=%s\n", c.Key, c.OldValue)
		case env.Changed:
			fmt.Fprintf(os.Stdout, "~ %s: %s -> %s\n", c.Key, c.OldValue, c.NewValue)
		}
	}

	return nil
}
