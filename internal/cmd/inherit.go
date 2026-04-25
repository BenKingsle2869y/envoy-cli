package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envoy-cli/internal/store"
)

var inheritCmd = &cobra.Command{
	Use:   "inherit",
	Short: "Inherit missing keys from a parent context into the active context",
	Long: `Copies keys that exist in the parent context but are missing in the
active context. Use --overwrite to replace existing values.`,
}

func init() {
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "run <parent-context>",
		Short: "Inherit keys from a parent context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInherit(cmd, args, overwrite)
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing keys in the active context")
	inheritCmd.AddCommand(cmd)
	rootCmd.AddCommand(inheritCmd)
}

func runInherit(cmd *cobra.Command, args []string, overwrite bool) error {
	parentCtx := args[0]
	activeCtx := ActiveContext()

	if parentCtx == activeCtx {
		return fmt.Errorf("parent context and active context are the same: %q", activeCtx)
	}

	passphrase := resolvePassphrase(cmd)
	if passphrase == "" {
		return fmt.Errorf("passphrase is required (use --passphrase or ENVOY_PASSPHRASE)")
	}

	parentPath := StorePathForContext(parentCtx)
	parentStore, err := store.Load(parentPath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load parent context %q: %w", parentCtx, err)
	}

	activePath := StorePathForContext(activeCtx)
	activeStore, err := store.Load(activePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load active context %q: %w", activeCtx, err)
	}

	inherited := 0
	for _, entry := range parentStore.Entries {
		_, exists := activeStore.Get(entry.Key)
		if exists && !overwrite {
			continue
		}
		activeStore.Set(entry.Key, entry.Value)
		inherited++
	}

	if inherited == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No keys inherited.")
		return nil
	}

	if err := store.Save(activePath, passphrase, activeStore); err != nil {
		return fmt.Errorf("failed to save active context: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Inherited %d key(s) from %q into %q.\n", inherited, parentCtx, activeCtx)
	return nil
}
