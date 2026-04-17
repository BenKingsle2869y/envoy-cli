package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var envDefaultCmd = &cobra.Command{
	Use:   "default <key>",
	Short: "Set a default value for a key if it is not already set",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvDefault,
}

func init() {
	envDefaultCmd.Flags().String("value", "", "Default value to set")
	envDefaultCmd.Flags().String("passphrase", "", "Passphrase for the store")
	_ = envDefaultCmd.MarkFlagRequired("value")
}

func runEnvDefault(cmd *cobra.Command, args []string) error {
	key := args[0]
	value, _ := cmd.Flags().GetString("value")
	passphrase, _ := cmd.Flags().GetString("passphrase")

	if passphrase == "" {
		passphrase = os.Getenv("ENVOY_PASSPHRASE")
	}
	if passphrase == "" {
		return fmt.Errorf("passphrase is required")
	}

	ctx := ActiveContext()
	path := StorePathForContext(ctx)

	st, err := loadStore(path, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	if _, exists := st.Get(key); exists {
		fmt.Fprintf(cmd.OutOrStdout(), "key %q already set, skipping\n", key)
		return nil
	}

	st.Set(key, value)
	if err := saveStore(path, passphrase, st); err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "default set: %s=%s\n", key, value)
	return nil
}
