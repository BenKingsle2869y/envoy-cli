package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var protectCmd = &cobra.Command{
	Use:   "protect",
	Short: "Manage protected (read-only) keys",
	Long:  "Mark keys as protected to prevent accidental modification or deletion.",
}

var protectAddCmd = &cobra.Command{
	Use:   "add <key>",
	Short: "Mark a key as protected",
	Args:  cobra.ExactArgs(1),
	RunE:  runProtectAdd,
}

var protectRemoveCmd = &cobra.Command{
	Use:   "remove <key>",
	Short: "Remove protection from a key",
	Args:  cobra.ExactArgs(1),
	RunE:  runProtectRemove,
}

var protectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all protected keys",
	Args:  cobra.NoArgs,
	RunE:  runProtectList,
}

func init() {
	protectCmd.AddCommand(protectAddCmd, protectRemoveCmd, protectListCmd)
	rootCmd.AddCommand(protectCmd)
}

func runProtectAdd(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	path := StorePathForContext(ctx)
	protected, err := LoadProtected(path)
	if err != nil {
		return fmt.Errorf("load protected: %w", err)
	}
	if err := MarkProtected(protected, args[0]); err != nil {
		return err
	}
	if err := SaveProtected(path, protected); err != nil {
		return fmt.Errorf("save protected: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "key %q is now protected\n", args[0])
	return nil
}

func runProtectRemove(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	path := StorePathForContext(ctx)
	protected, err := LoadProtected(path)
	if err != nil {
		return fmt.Errorf("load protected: %w", err)
	}
	if err := UnmarkProtected(protected, args[0]); err != nil {
		return err
	}
	if err := SaveProtected(path, protected); err != nil {
		return fmt.Errorf("save protected: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "key %q is no longer protected\n", args[0])
	return nil
}

func runProtectList(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	path := StorePathForContext(ctx)
	protected, err := LoadProtected(path)
	if err != nil {
		return fmt.Errorf("load protected: %w", err)
	}
	keys := ProtectedKeys(protected)
	if len(keys) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no protected keys")
		return nil
	}
	for _, k := range keys {
		fmt.Fprintln(cmd.OutOrStdout(), k)
	}
	return nil
}
