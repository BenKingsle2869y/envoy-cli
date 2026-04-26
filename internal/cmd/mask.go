package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/cmd"
	"envoy-cli/internal/store"
)

var maskCmd = &cobra.Command{
	Use:   "mask",
	Short: "Mask or unmask values in the active context",
}

var maskAddCmd = &cobra.Command{
	Use:   "add <key>",
	Short: "Mark a key as masked (value shown as ***)",
	Args:  cobra.ExactArgs(1),
	RunE:  runMaskAdd,
}

var maskRemoveCmd = &cobra.Command{
	Use:   "remove <key>",
	Short: "Unmark a key as masked",
	Args:  cobra.ExactArgs(1),
	RunE:  runMaskRemove,
}

var maskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all masked keys",
	RunE:  runMaskList,
}

func init() {
	maskCmd.AddCommand(maskAddCmd, maskRemoveCmd, maskListCmd)
	maskCmd.PersistentFlags().String("passphrase", "", "Passphrase for the store")
	rootCmd.AddCommand(maskCmd)
}

func runMaskAdd(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	path := store.DefaultStorePath(ctx)
	passphrase, err := resolvePassphrase(cmd)
	if err != nil {
		return err
	}
	s, err := store.Load(path, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}
	key := args[0]
	if _, ok := s.Get(key); !ok {
		return fmt.Errorf("key %q not found", key)
	}
	masks, err := LoadMasks(maskFilePath(ctx))
	if err != nil {
		return err
	}
	if err := MarkMasked(masks, key); err != nil {
		return err
	}
	if err := SaveMasks(maskFilePath(ctx), masks); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "key %q marked as masked\n", key)
	return nil
}

func runMaskRemove(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	masks, err := LoadMasks(maskFilePath(ctx))
	if err != nil {
		return err
	}
	key := args[0]
	if err := UnmarkMasked(masks, key); err != nil {
		return err
	}
	if err := SaveMasks(maskFilePath(ctx), masks); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "key %q unmasked\n", key)
	return nil
}

func runMaskList(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	masks, err := LoadMasks(maskFilePath(ctx))
	if err != nil {
		return err
	}
	keys := MaskedKeys(masks)
	if len(keys) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no masked keys")
		return nil
	}
	fmt.Fprintln(cmd.OutOrStdout(), strings.Join(keys, "\n"))
	return nil
}
