package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/store"
)

var (
	searchKeys   bool
	searchValues bool
	searchExact  bool
)

var searchCmd = &cobra.Command{
	Use:   "search <pattern>",
	Short: "Search for keys or values matching a pattern",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().BoolVarP(&searchKeys, "keys", "k", false, "Search only in keys")
	searchCmd.Flags().BoolVarP(&searchValues, "values", "v", false, "Search only in values")
	searchCmd.Flags().BoolVarP(&searchExact, "exact", "e", false, "Match exact string (case-sensitive)")
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	pattern := args[0]

	passphrase, err := crypto.ResolvePassphrase()
	if err != nil {
		return fmt.Errorf("passphrase error: %w", err)
	}

	storePath := store.DefaultStorePath(ActiveContext())
	s, err := store.Load(storePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	// Default: search both keys and values
	if !searchKeys && !searchValues {
		searchKeys = true
		searchValues = true
	}

	matches := 0
	for k, v := range s.Vars {
		keyMatch := searchKeys && matchPattern(k, pattern, searchExact)
		valMatch := searchValues && matchPattern(v, pattern, searchExact)
		if keyMatch || valMatch {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
			matches++
		}
	}

	if matches == 0 {
		fmt.Fprintf(os.Stdout, "no matches found for %q\n", pattern)
	}

	return nil
}

func matchPattern(s, pattern string, exact bool) bool {
	if exact {
		return s == pattern
	}
	return strings.Contains(strings.ToLower(s), strings.ToLower(pattern))
}
