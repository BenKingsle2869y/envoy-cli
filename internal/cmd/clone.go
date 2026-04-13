package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/store"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <source-context> <dest-context>",
	Short: "Clone an environment context into a new one",
	Long: `Clone creates a full copy of an existing environment context
into a new context name. The destination must not already exist
unless --overwrite is specified.`,
	Args: cobra.ExactArgs(2),
	RunE: runClone,
}

func init() {
	cloneCmd.Flags().BoolP("overwrite", "o", false, "Overwrite destination if it already exists")
	RootCmd.AddCommand(cloneCmd)
}

func runClone(cmd *cobra.Command, args []string) error {
	srcContext := args[0]
	dstContext := args[1]

	if srcContext == dstContext {
		return fmt.Errorf("source and destination context must differ")
	}

	overwrite, _ := cmd.Flags().GetBool("overwrite")

	passphrase := os.Getenv("ENVOY_PASSPHRASE")
	if passphrase == "" {
		return fmt.Errorf("ENVOY_PASSPHRASE environment variable is required")
	}

	srcPath := StorePathForContext(srcContext)
	dstPath := StorePathForContext(dstContext)

	if _, err := os.Stat(dstPath); err == nil && !overwrite {
		return fmt.Errorf("destination context %q already exists; use --overwrite to replace it", dstContext)
	}

	data, err := store.Load(srcPath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load source context %q: %w", srcContext, err)
	}

	if err := store.Save(dstPath, passphrase, data); err != nil {
		return fmt.Errorf("failed to save destination context %q: %w", dstContext, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "cloned context %q → %q (%d keys)\n", srcContext, dstContext, len(data))
	return nil
}
