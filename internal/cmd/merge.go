package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/envoy-cli/internal/cmd"
	"github.com/envoy-cli/internal/crypto"
	"github.com/envoy-cli/internal/store"
)

var mergeOverwrite bool

var mergeCmd = &cobra.Command{
	Use:   "merge <source-context> <dest-context>",
	Short: "Merge variables from one context into another",
	Long: `Merge all key-value pairs from the source context into the destination context.

By default, existing keys in the destination are preserved. Use --overwrite to
replace conflicting keys with values from the source context.`,
	Args: cobra.ExactArgs(2),
	RunE: runMerge,
}

func init() {
	mergeCmd.Flags().BoolVar(&mergeOverwrite, "overwrite", false, "Overwrite existing keys in destination")
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	srcCtx := args[0]
	dstCtx := args[1]

	if srcCtx == dstCtx {
		return fmt.Errorf("source and destination contexts must be different")
	}

	passphrase, err := crypto.ResolvePassphrase("")
	if err != nil {
		return fmt.Errorf("could not resolve passphrase: %w", err)
	}

	srcPath := StorePathForContext(srcCtx)
	srcStore, err := store.Load(srcPath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load source context %q: %w", srcCtx, err)
	}

	dstPath := StorePathForContext(dstCtx)
	dstStore, err := store.Load(dstPath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load destination context %q: %w", dstCtx, err)
	}

	merged := 0
	skipped := 0

	for key, entry := range srcStore.Entries {
		if _, exists := dstStore.Entries[key]; exists && !mergeOverwrite {
			skipped++
			continue
		}
		dstStore.Entries[key] = entry
		merged++
	}

	if err := store.Save(dstPath, passphrase, dstStore); err != nil {
		return fmt.Errorf("failed to save destination context: %w", err)
	}

	fmt.Printf("Merged %d key(s) from %q into %q (%d skipped).\n", merged, srcCtx, dstCtx, skipped)
	return nil
}
