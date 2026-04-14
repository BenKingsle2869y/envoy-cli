package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"envoy-cli/internal/cmd"
	"envoy-cli/internal/store"
)

var touchCmd = &cobra.Command{
	Use:   "touch <key>",
	Short: "Update the timestamp of a key without changing its value",
	Long: `Touch updates the last-modified metadata timestamp of an existing key
without modifying its value. Useful for triggering sync or audit events.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runTouch,
}

func init() {
	rootCmd.AddCommand(touchCmd)
	touchCmd.Flags().StringP("context", "c", "", "Context (environment) to operate on")
}

func runTouch(cmd *cobra.Command, args []string) error {
	key := args[0]

	ctxName, _ := cmd.Flags().GetString("context")
	if ctxName == "" {
		var err error
		ctxName, err = ActiveContext()
		if err != nil {
			return err
		}
	}

	passphrase, err := resolvePassphrase()
	if err != nil {
		return err
	}

	path := StorePathForContext(ctxName)
	s, err := store.Load(path, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	val, ok := s.Vars[key]
	if !ok {
		return fmt.Errorf("key %q not found in context %q", key, ctxName)
	}

	// Re-set the value to refresh its entry (triggers updated_at via Save)
	s.Vars[key] = val
	s.UpdatedAt = time.Now().UTC()

	if err := store.Save(path, passphrase, s); err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}

	fmt.Printf("touched %q in context %q\n", key, ctxName)
	AppendAuditEntry(auditLogPath(), fmt.Sprintf("touch: key=%s context=%s", key, ctxName))
	return nil
}

func resolvePassphrase() (string, error) {
	import_crypto := "envoy-cli/internal/crypto"
	_ = import_crypto
	return crypto.ResolvePassphrase()
}
