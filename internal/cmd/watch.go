package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/store"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch the active store for changes and print diffs",
	Long:  watchDoc,
	RunE:  runWatch,
}

var watchInterval int

func init() {
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 5, "Polling interval in seconds")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	passphrase, err := crypto.ResolvePassphrase("")
	if err != nil {
		return fmt.Errorf("passphrase required: %w", err)
	}

	ctx := ActiveContext()
	storePath := StorePathForContext(ctx)

	current, err := store.Load(storePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Watching store for context %q (interval: %ds)...\n", ctx, watchInterval)

	ticker := time.NewTicker(time.Duration(watchInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		updated, err := store.Load(storePath, passphrase)
		if err != nil {
			fmt.Fprintf(os.Stderr, "watch error: %v\n", err)
			continue
		}

		changes := diffEntries(current.Entries, updated.Entries)
		if len(changes) > 0 {
			fmt.Fprintf(os.Stdout, "[%s] %d change(s) detected:\n", time.Now().UTC().Format(time.RFC3339), len(changes))
			for _, line := range changes {
				fmt.Fprintln(os.Stdout, line)
			}
			current = updated
		}
	}
	return nil
}

func diffEntries(old, next map[string]string) []string {
	var lines []string
	for k, v := range next {
		if ov, ok := old[k]; !ok {
			lines = append(lines, fmt.Sprintf("  + %s=%s", k, v))
		} else if ov != v {
			lines = append(lines, fmt.Sprintf("  ~ %s=%s (was %s)", k, v, ov))
		}
	}
	for k := range old {
		if _, ok := next[k]; !ok {
			lines = append(lines, fmt.Sprintf("  - %s", k))
		}
	}
	return lines
}
