package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Pin or unpin environment variables to prevent accidental modification",
}

var pinAddCmd = &cobra.Command{
	Use:   "add <key>",
	Short: "Pin a key so it cannot be overwritten by import, merge, or copy",
	Args:  cobra.ExactArgs(1),
	RunE:  runPinAdd,
}

var pinRemoveCmd = &cobra.Command{
	Use:   "remove <key>",
	Short: "Unpin a key to allow modifications",
	Args:  cobra.ExactArgs(1),
	RunE:  runPinRemove,
}

var pinListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all pinned keys in the current context",
	RunE:  runPinList,
}

func init() {
	pinCmd.AddCommand(pinAddCmd)
	pinCmd.AddCommand(pinRemoveCmd)
	pinCmd.AddCommand(pinListCmd)
	rootCmd.AddCommand(pinCmd)
}

func runPinAdd(cmd *cobra.Command, args []string) error {
	key := strings.TrimSpace(args[0])
	ctx := ActiveContext()
	path := StorePathForContext(ctx)
	passphrase, err := resolvePassphrase(path)
	if err != nil {
		return err
	}
	st, err := store.Load(path, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}
	if _, ok := st.Vars[key]; !ok {
		return fmt.Errorf("key %q not found in store", key)
	}
	if err := AddPin(st.Tags, key); err != nil {
		return err
	}
	if err := store.Save(path, passphrase, st); err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "pinned %q\n", key)
	return nil
}

func runPinRemove(cmd *cobra.Command, args []string) error {
	key := strings.TrimSpace(args[0])
	ctx := ActiveContext()
	path := StorePathForContext(ctx)
	passphrase, err := resolvePassphrase(path)
	if err != nil {
		return err
	}
	st, err := store.Load(path, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}
	if err := RemovePin(st.Tags, key); err != nil {
		return err
	}
	if err := store.Save(path, passphrase, st); err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "unpinned %q\n", key)
	return nil
}

func runPinList(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	path := StorePathForContext(ctx)
	passphrase, err := resolvePassphrase(path)
	if err != nil {
		return err
	}
	st, err := store.Load(path, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}
	keys := PinnedKeys(st.Tags)
	if len(keys) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no pinned keys")
		return nil
	}
	for _, k := range keys {
		fmt.Fprintln(cmd.OutOrStdout(), k)
	}
	return nil
}
