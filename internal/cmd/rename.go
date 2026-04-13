package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename <old-key> <new-key>",
	Short: "Rename an environment variable key",
	Long: `Rename renames an existing key in the active environment store.

The value associated with the old key is preserved under the new key name.
If the old key does not exist or the new key already exists, the command fails.`,
	Args: cobra.ExactArgs(2),
	RunE: runRename,
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	oldKey := args[0]
	newKey := args[1]

	if oldKey == newKey {
		return fmt.Errorf("old key and new key are the same: %q", oldKey)
	}

	passphrase, err := resolvePassphrase()
	if err != nil {
		return fmt.Errorf("passphrase error: %w", err)
	}

	storePath := StorePathForContext(ActiveContext())
	st, err := loadStore(storePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	val, ok := st.Entries[oldKey]
	if !ok {
		return fmt.Errorf("key %q not found", oldKey)
	}

	if _, exists := st.Entries[newKey]; exists {
		return fmt.Errorf("key %q already exists; use --force to overwrite", newKey)
	}

	st.Entries[newKey] = val
	delete(st.Entries, oldKey)

	if err := saveStore(storePath, passphrase, st); err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "renamed %q → %q\n", oldKey, newKey)

	_ = AppendAuditEntry(AuditEntry{
		Action:  "rename",
		Context: ActiveContext(),
		Key:     fmt.Sprintf("%s->%s", oldKey, newKey),
	})

	return nil
}
