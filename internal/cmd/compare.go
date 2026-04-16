package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envoy-cli/internal/cmd"
	"envoy-cli/internal/store"
)

var compareCmd = &cobra.Command{
	Use:   "compare <context-a> <context-b>",
	Short: "Compare keys between two contexts",
	Args:  cobra.ExactArgs(2),
	RunE:  runCompare,
}

func init() {
	compareCmd.Flags().String("passphrase", "", "Passphrase for both stores")
	compareCmd.Flags().String("passphrase-a", "", "Passphrase for context A")
	compareCmd.Flags().String("passphrase-b", "", "Passphrase for context B")
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	ctxA, ctxB := args[0], args[1]
	if ctxA == ctxB {
		return fmt.Errorf("contexts must be different")
	}

	shared, _ := cmd.Flags().GetString("passphrase")
	passphraseA, _ := cmd.Flags().GetString("passphrase-a")
	passphraseB, _ := cmd.Flags().GetString("passphrase-b")
	if passphraseA == "" {
		passphraseA = shared
	}
	if passphraseB == "" {
		passphraseB = shared
	}

	pathA := StorePathForContext(ctxA)
	pathB := StorePathForContext(ctxB)

	storeA, err := store.Load(pathA, passphraseA)
	if err != nil {
		return fmt.Errorf("failed to load context %q: %w", ctxA, err)
	}
	storeB, err := store.Load(pathB, passphraseB)
	if err != nil {
		return fmt.Errorf("failed to load context %q: %w", ctxB, err)
	}

	result := compareStores(storeA.Entries, storeB.Entries)

	if len(result.OnlyInA)+len(result.OnlyInB)+len(result.Different) == 0 {
		fmt.Fprintln(os.Stdout, "No differences found.")
		return nil
	}

	for _, k := range sorted(result.OnlyInA) {
		fmt.Fprintf(os.Stdout, "< only in %s: %s\n", ctxA, k)
	}
	for _, k := range sorted(result.OnlyInB) {
		fmt.Fprintf(os.Stdout, "> only in %s: %s\n", ctxB, k)
	}
	for _, k := range sorted(result.Different) {
		fmt.Fprintf(os.Stdout, "~ differs: %s\n", k)
	}
	return nil
}

func sorted(keys []string) []string {
	sort.Strings(keys)
	return keys
}
