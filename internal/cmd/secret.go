package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage secret metadata for env keys",
}

var secretMarkCmd = &cobra.Command{
	Use:   "mark <key>",
	Short: "Mark a key as secret (masked in output)",
	Args:  cobra.ExactArgs(1),
	RunE:  runSecretMark,
}

var secretUnmarkCmd = &cobra.Command{
	Use:   "unmark <key>",
	Short: "Unmark a key as secret",
	Args:  cobra.ExactArgs(1),
	RunE:  runSecretUnmark,
}

var secretListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all keys marked as secret",
	RunE:  runSecretList,
}

func init() {
	secretCmd.AddCommand(secretMarkCmd, secretUnmarkCmd, secretListCmd)
	secretCmd.PersistentFlags().String("passphrase", "", "Passphrase for the store")
	rootCmd.AddCommand(secretCmd)
}

func runSecretMark(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	path := StorePathForContext(ctx)
	pass, _ := cmd.Flags().GetString("passphrase")
	if pass == "" {
		pass = resolvePassphrase(path)
	}
	sf := secretFilePath(path)
	secrets, err := LoadSecrets(sf)
	if err != nil {
		return err
	}
	if err := MarkSecret(secrets, args[0]); err != nil {
		return err
	}
	if err := SaveSecrets(sf, secrets); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Key %q marked as secret\n", args[0])
	return nil
}

func runSecretUnmark(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	path := StorePathForContext(ctx)
	sf := secretFilePath(path)
	secrets, err := LoadSecrets(sf)
	if err != nil {
		return err
	}
	if err := UnmarkSecret(secrets, args[0]); err != nil {
		return err
	}
	if err := SaveSecrets(sf, secrets); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Key %q unmarked as secret\n", args[0])
	return nil
}

func runSecretList(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	path := StorePathForContext(ctx)
	sf := secretFilePath(path)
	secrets, err := LoadSecrets(sf)
	if err != nil {
		return err
	}
	keys := SecretKeys(secrets)
	if len(keys) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No secret keys")
		return nil
	}
	fmt.Fprintln(cmd.OutOrStdout(), strings.Join(keys, "\n"))
	return nil
}
