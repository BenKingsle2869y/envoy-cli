package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var envExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all variables from the active context as shell export statements",
	Long: `Prints all key=value pairs from the active context as shell-compatible
export statements. Useful for sourcing into a shell session.

Example:
  eval $(envoy env export)
  envoy env export --context production > prod.env`,
	RunE: runEnvExport,
}

var envExportContext string
var envExportNoBlanks bool

func init() {
	envExportCmd.Flags().StringVar(&envExportContext, "context", "", "Context to export from (defaults to active context)")
	envExportCmd.Flags().BoolVar(&envExportNoBlanks, "no-blanks", false, "Skip keys with empty values")
	envCmd.AddCommand(envExportCmd)
}

func runEnvExport(cmd *cobra.Command, args []string) error {
	ctx := envExportContext
	if ctx == "" {
		ctx = ActiveContext()
	}

	passphrase := os.Getenv("ENVOY_PASSPHRASE")
	if passphrase == "" {
		return fmt.Errorf("ENVOY_PASSPHRASE is not set")
	}

	path := StorePathForContext(ctx)
	store, err := loadStoreWithPassphrase(path, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load context %q: %w", ctx, err)
	}

	keys := make([]string, 0, len(store))
	for k := range store {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := store[k]
		if envExportNoBlanks && v == "" {
			continue
		}
		fmt.Fprintf(&sb, "export %s=%s\n", k, QuoteValue(v))
	}

	fmt.Fprint(cmd.OutOrStdout(), sb.String())
	return nil
}
