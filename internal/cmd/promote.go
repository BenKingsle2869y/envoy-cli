package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envoy-cli/internal/cmd"
	"envoy-cli/internal/store"
)

var promoteCmd = &cobra.Command{
	Use:   "promote <src-context> <dst-context>",
	Short: "Promote env variables from one context to another",
	Long:  promoteDoc,
	Args:  cobra.ExactArgs(2),
	RunE:  runPromote,
}

func init() {
	promoteCmd.Flags().Bool("overwrite", false, "Overwrite existing keys in destination")
	promoteCmd.Flags().String("passphrase", "", "Passphrase for both stores (or set ENVOY_PASSPHRASE)")
	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, args []string) error {
	src, dst := args[0], args[1]
	if src == dst {
		return fmt.Errorf("source and destination contexts must differ")
	}

	overwrite, _ := cmd.Flags().GetBool("overwrite")
	passphrase, _ := cmd.Flags().GetString("passphrase")
	if passphrase == "" {
		var err error
		passphrase, err = resolvePassphrase()
		if err != nil {
			return err
		}
	}

	srcPath := StorePathForContext(src)
	dstPath := StorePathForContext(dst)

	srcStore, err := store.Load(srcPath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load source context %q: %w", src, err)
	}

	dstStore, err := store.Load(dstPath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load destination context %q: %w", dst, err)
	}

	promoted := 0
	skipped := 0
	for k, v := range srcStore.Vars {
		if _, exists := dstStore.Vars[k]; exists && !overwrite {
			skipped++
			continue
		}
		dstStore.Vars[k] = v
		promoted++
	}

	if err := store.Save(dstPath, passphrase, dstStore); err != nil {
		return fmt.Errorf("failed to save destination context: %w", err)
	}

	fmt.Printf("Promoted %d key(s) from %q to %q (%d skipped).\n", promoted, src, dst, skipped)
	return nil
}
