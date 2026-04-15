package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var ttlCmd = &cobra.Command{
	Use:   "ttl",
	Short: "Manage time-to-live expiry on env keys",
}

var ttlSetCmd = &cobra.Command{
	Use:   "set <key> <duration>",
	Short: "Set an expiry duration on a key (e.g. 24h, 7d)",
	Args:  cobra.ExactArgs(2),
	RunE:  runTTLSet,
}

var ttlShowCmd = &cobra.Command{
	Use:   "show <key>",
	Short: "Show the expiry time for a key",
	Args:  cobra.ExactArgs(1),
	RunE:  runTTLShow,
}

var ttlClearCmd = &cobra.Command{
	Use:   "clear <key>",
	Short: "Remove the expiry on a key",
	Args:  cobra.ExactArgs(1),
	RunE:  runTTLClear,
}

func init() {
	ttlCmd.AddCommand(ttlSetCmd, ttlShowCmd, ttlClearCmd)
	ttlCmd.PersistentFlags().String("passphrase", "", "Passphrase to decrypt the store")
	RootCmd.AddCommand(ttlCmd)
}

func runTTLSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	durStr := args[1]

	dur, err := parseTTLDuration(durStr)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", durStr, err)
	}

	pass, _ := cmd.Flags().GetString("passphrase")
	ctx := ActiveContext()
	store, err := loadStoreWithPassphrase(ctx, pass)
	if err != nil {
		return err
	}

	if _, ok := store.Get(key); !ok {
		return fmt.Errorf("key %q not found", key)
	}

	expiry := time.Now().UTC().Add(dur)
	if err := SetTTL(ttlFilePath(ctx), key, expiry); err != nil {
		return fmt.Errorf("failed to save TTL: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "TTL set: %s expires at %s\n", key, expiry.Format(time.RFC3339))
	return nil
}

func runTTLShow(cmd *cobra.Command, args []string) error {
	key := args[0]
	ctx := ActiveContext()

	ttls, err := LoadTTLs(ttlFilePath(ctx))
	if err != nil {
		return err
	}

	expiry, ok := ttls[key]
	if !ok {
		fmt.Fprintf(cmd.OutOrStdout(), "No TTL set for %q\n", key)
		return nil
	}

	if time.Now().UTC().After(expiry) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s expired at %s\n", key, expiry.Format(time.RFC3339))
	} else {
		remaining := time.Until(expiry).Round(time.Second)
		fmt.Fprintf(cmd.OutOrStdout(), "%s expires at %s (%s remaining)\n", key, expiry.Format(time.RFC3339), remaining)
	}
	return nil
}

func runTTLClear(cmd *cobra.Command, args []string) error {
	key := args[0]
	ctx := ActiveContext()

	if err := ClearTTL(ttlFilePath(ctx), key); err != nil {
		return fmt.Errorf("failed to clear TTL: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "TTL cleared for %q\n", key)
	return nil
}

func loadStoreWithPassphrase(ctx, pass string) (interface{ Get(string) (string, bool) }, error) {
	if pass == "" {
		pass = os.Getenv("ENVOY_PASSPHRASE")
	}
	if pass == "" {
		return nil, fmt.Errorf("passphrase required: use --passphrase or ENVOY_PASSPHRASE")
	}
	path := StorePathForContext(ctx)
	return loadStore(path, pass)
}
