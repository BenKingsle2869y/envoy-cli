package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy <source-context> <dest-context>",
	Short: "Copy all env variables from one context to another",
	Args:  cobra.ExactArgs(2),
	RunE:  runCopy,
}

func init() {
	copyCmd.Flags().Bool("overwrite", false, "Overwrite existing keys in the destination context")
	rootCmd.AddCommand(copyCmd)
}

func runCopy(cmd *cobra.Command, args []string) error {
	srcCtx := args[0]
	dstCtx := args[1]

	if srcCtx == dstCtx {
		return fmt.Errorf("source and destination contexts must be different")
	}

	overwrite, _ := cmd.Flags().GetBool("overwrite")

	passphrase, err := resolvePassphraseForCmd(cmd)
	if err != nil {
		return err
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

	copied := 0
	skipped := 0
	for k, v := range srcStore.Vars {
		if _, exists := dstStore.Vars[k]; exists && !overwrite {
			skipped++
			continue
		}
		dstStore.Vars[k] = v
		copied++
	}

	if err := store.Save(dstPath, dstStore, passphrase); err != nil {
		return fmt.Errorf("failed to save destination context: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Copied %d variable(s) from %q to %q (%d skipped).\n", copied, srcCtx, dstCtx, skipped)
	return nil
}
