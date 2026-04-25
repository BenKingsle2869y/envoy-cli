package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var placeholderCmd = &cobra.Command{
	Use:   "placeholder",
	Short: "Manage placeholder keys with default values",
}

var placeholderSetCmd = &cobra.Command{
	Use:   "set <key> <default>",
	Short: "Register a placeholder key with a default value",
	Args:  cobra.ExactArgs(2),
	RunE:  runPlaceholderSet,
}

var placeholderListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all placeholder keys and their defaults",
	RunE:  runPlaceholderList,
}

var placeholderApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply placeholder defaults for any missing keys in the store",
	RunE:  runPlaceholderApply,
}

func init() {
	placeholderCmd.AddCommand(placeholderSetCmd)
	placeholderCmd.AddCommand(placeholderListCmd)
	placeholderCmd.AddCommand(placeholderApplyCmd)
	placeholderCmd.PersistentFlags().String("passphrase", "", "Passphrase to decrypt the store")
	rootCmd.AddCommand(placeholderCmd)
}

func runPlaceholderSet(cmd *cobra.Command, args []string) error {
	key := strings.TrimSpace(args[0])
	defaultVal := args[1]
	ph, err := LoadPlaceholders(placeholderFilePath())
	if err != nil {
		return fmt.Errorf("load placeholders: %w", err)
	}
	ph[key] = defaultVal
	return SavePlaceholders(placeholderFilePath(), ph)
}

func runPlaceholderList(cmd *cobra.Command, args []string) error {
	ph, err := LoadPlaceholders(placeholderFilePath())
	if err != nil {
		return fmt.Errorf("load placeholders: %w", err)
	}
	if len(ph) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No placeholders defined.")
		return nil
	}
	for k, v := range ph {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
	}
	return nil
}

func runPlaceholderApply(cmd *cobra.Command, args []string) error {
	pass, _ := cmd.Flags().GetString("passphrase")
	if pass == "" {
		pass = os.Getenv("ENVOY_PASSPHRASE")
	}
	if pass == "" {
		return fmt.Errorf("passphrase required")
	}
	ctx := ActiveContext()
	path := StorePathForContext(ctx)
	st, err := store.Load(path, pass)
	if err != nil {
		return fmt.Errorf("load store: %w", err)
	}
	ph, err := LoadPlaceholders(placeholderFilePath())
	if err != nil {
		return fmt.Errorf("load placeholders: %w", err)
	}
	applied := 0
	for k, def := range ph {
		if _, ok := st.Entries[k]; !ok {
			st.Entries[k] = store.Entry{Value: def}
			applied++
		}
	}
	if applied == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No placeholders applied.")
		return nil
	}
	if err := store.Save(path, pass, st); err != nil {
		return fmt.Errorf("save store: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Applied %d placeholder(s).\n", applied)
	return nil
}
