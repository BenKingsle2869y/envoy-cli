package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/cmd"
	"envoy-cli/internal/store"
)

var envFillCmd = &cobra.Command{
	Use:   "fill",
	Short: "Fill missing env vars from the active store into the current shell environment",
	Long:  envFillDoc,
	RunE:  runEnvFill,
}

func init() {
	envFillCmd.Flags().StringP("passphrase", "p", "", "Passphrase to decrypt the store")
	envFillCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing environment variables")
	envFillCmd.Flags().BoolP("export", "e", false, "Print export statements instead of applying")
	rootCmd.AddCommand(envFillCmd)
}

func runEnvFill(cmd *cobra.Command, args []string) error {
	passphrase, _ := cmd.Flags().GetString("passphrase")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	exportMode, _ := cmd.Flags().GetBool("export")

	if passphrase == "" {
		passphrase = os.Getenv("ENVOY_PASSPHRASE")
	}
	if passphrase == "" {
		return fmt.Errorf("passphrase is required (use --passphrase or ENVOY_PASSPHRASE)")
	}

	ctx := ActiveContext()
	path := StorePathForContext(ctx)

	s, err := store.Load(path, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	filled := 0
	for _, entry := range s.Entries {
		_, exists := os.LookupEnv(entry.Key)
		if exists && !overwrite {
			continue
		}
		if exportMode {
			fmt.Printf("export %s=%q\n", entry.Key, entry.Value)
		} else {
			if err := os.Setenv(entry.Key, entry.Value); err != nil {
				return fmt.Errorf("failed to set %s: %w", entry.Key, err)
			}
		}
		filled++
	}

	if !exportMode {
		fmt.Printf("Filled %d variable(s) into environment.\n", filled)
	}
	return nil
}
