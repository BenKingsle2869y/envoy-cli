package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/cmd"
	"envoy-cli/internal/env"
	"envoy-cli/internal/store"
)

func init() {
	fmtCmd := &cobra.Command{
		Use:   "fmt",
		Short: "Format and sort the active env store",
		Long:  "Normalizes the active env store by sorting keys alphabetically and re-serializing entries with consistent quoting.",
		RunE:  runFmt,
	}
	fmtCmd.Flags().Bool("check", false, "Exit with non-zero status if store is not already formatted, without writing changes")
	rootCmd.AddCommand(fmtCmd)
}

func runFmt(cmd *cobra.Command, args []string) error {
	passphrase := os.Getenv("ENVOY_PASSPHRASE")
	if passphrase == "" {
		return fmt.Errorf("ENVOY_PASSPHRASE is not set")
	}

	ctx := ActiveContext()
	storePath := StorePathForContext(ctx)

	s, err := store.Load(storePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	keys := make([]string, 0, len(s.Vars))
	for k := range s.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		lines = append(lines, k+"="+env.QuoteValue(s.Vars[k]))
	}
	formatted := strings.Join(lines, "\n") + "\n"

	checkOnly, _ := cmd.Flags().GetBool("check")
	if checkOnly {
		current := env.Serialize(s.Vars)
		if current != formatted {
			return fmt.Errorf("store is not formatted; run 'envoy fmt' to fix")
		}
		fmt.Println("store is already formatted")
		return nil
	}

	orderedVars := make(map[string]string, len(keys))
	for _, k := range keys {
		orderedVars[k] = s.Vars[k]
	}
	s.Vars = orderedVars

	if err := store.Save(storePath, passphrase, s); err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}

	fmt.Printf("formatted %d keys in context %q\n", len(keys), ctx)
	return nil
}
