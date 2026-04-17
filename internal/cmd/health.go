package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check the health of the current context store",
	Long:  "Verifies that the active context store exists, is readable, and can be decrypted with the provided passphrase.",
	RunE:  runHealth,
}

func init() {
	healthCmd.Flags().String("passphrase", "", "Passphrase to decrypt the store")
	rootCmd.AddCommand(healthCmd)
}

func runHealth(cmd *cobra.Command, args []string) error {
	passphrase, _ := cmd.Flags().GetString("passphrase")
	if passphrase == "" {
		passphrase = os.Getenv("ENVOY_PASSPHRASE")
	}
	if passphrase == "" {
		return fmt.Errorf("passphrase required: use --passphrase or set ENVOY_PASSPHRASE")
	}

	ctx := ActiveContext()
	storePath := StorePathForContext(ctx)

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		return fmt.Errorf("store not found for context %q at %s", ctx, storePath)
	}

	start := time.Now()
	_, err := loadStoreWithPassphrase(storePath, passphrase)
	if err != nil {
		return fmt.Errorf("store could not be decrypted: %w", err)
	}
	elapsed := time.Since(start)

	fmt.Printf("context:  %s\n", ctx)
	fmt.Printf("store:    %s\n", storePath)
	fmt.Printf("status:   OK\n")
	fmt.Printf("latency:  %s\n", elapsed.Round(time.Microsecond))
	return nil
}
