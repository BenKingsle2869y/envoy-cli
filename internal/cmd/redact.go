package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var redactCmd = &cobra.Command{
	Use:   "redact",
	Short: "Print env vars with secret values masked",
	Long:  "Display all keys in the current context with secret values replaced by ****.",
	RunE:  runRedact,
}

func init() {
	redactCmd.Flags().StringP("passphrase", "p", "", "Passphrase to decrypt the store")
	redactCmd.Flags().StringP("context", "c", "", "Context (environment) to use")
	rootCmd.AddCommand(redactCmd)
}

func runRedact(cmd *cobra.Command, args []string) error {
	passphrase, _ := cmd.Flags().GetString("passphrase")
	ctxName, _ := cmd.Flags().GetString("context")

	if passphrase == "" {
		var err error
		passphrase, err = ResolvePassphrase(cmd)
		if err != nil {
			return err
		}
	}

	if ctxName == "" {
		ctxName = ActiveContext()
	}

	path := StorePathForContext(ctxName)
	store, err := loadStoreWithPassphrase(path, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	secrets, err := LoadSecrets(ctxName)
	if err != nil {
		secrets = map[string]bool{}
	}

	for _, e := range store.Entries {
		value := e.Value
		if secrets[e.Key] || looksSecret(e.Key) {
			value = strings.Repeat("*", 8)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", e.Key, value)
	}
	return nil
}

func looksSecret(key string) bool {
	upper := strings.ToUpper(key)
	for _, kw := range []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "PRIVATE_KEY", "API_KEY"} {
		if strings.Contains(upper, kw) {
			return true
		}
	}
	return false
}
