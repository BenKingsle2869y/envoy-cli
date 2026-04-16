package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var expireCmd = &cobra.Command{
	Use:   "expire",
	Short: "List keys that have expired or are expiring soon",
	Long:  `Scans the active context store and reports keys whose TTL has passed or will expire within a given threshold.`,
	RunE:  runExpire,
}

func init() {
	expireCmd.Flags().StringP("passphrase", "p", "", "Passphrase to decrypt the store")
	expireCmd.Flags().IntP("within", "w", 0, "Warn about keys expiring within N hours")
	rootCmd.AddCommand(expireCmd)
}

func runExpire(cmd *cobra.Command, args []string) error {
	passphrase, _ := cmd.Flags().GetString("passphrase")
	withinHours, _ := cmd.Flags().GetInt("within")

	if passphrase == "" {
		var err error
		passphrase, err = ResolvePassphrase(cmd)
		if err != nil {
			return err
		}
	}

	ctx := ActiveContext()
	storePath := StorePathForContext(ctx)

	ttls, err := LoadTTLs(storePath)
	if err != nil {
		return fmt.Errorf("failed to load TTLs: %w", err)
	}

	now := time.Now().UTC()
	threshold := now.Add(time.Duration(withinHours) * time.Hour)

	found := false
	for key, expiry := range ttls {
		if expiry.Before(now) {
			fmt.Fprintf(cmd.OutOrStdout(), "EXPIRED   %s (expired %s ago)\n", key, now.Sub(expiry).Round(time.Second))
			found = true
		} else if withinHours > 0 && expiry.Before(threshold) {
			fmt.Fprintf(cmd.OutOrStdout(), "EXPIRING  %s (expires in %s)\n", key, expiry.Sub(now).Round(time.Second))
			found = true
		}
	}

	if !found {
		fmt.Fprintln(cmd.OutOrStdout(), "No expired or expiring keys found.")
	}

	return nil
}
