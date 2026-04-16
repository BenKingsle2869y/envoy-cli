package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var castCmd = &cobra.Command{
	Use:   "cast",
	Short: "Cast env values to typed representations",
	Long:  "Display env values with their inferred types (string, int, bool, float).",
	RunE:  runCast,
}

func init() {
	castCmd.Flags().StringP("passphrase", "p", "", "Passphrase to decrypt the store")
	castCmd.Flags().StringP("context", "c", "", "Context (environment) to use")
	rootCmd.AddCommand(castCmd)
}

func runCast(cmd *cobra.Command, args []string) error {
	passphrase, _ := cmd.Flags().GetString("passphrase")
	ctx, _ := cmd.Flags().GetString("context")

	if passphrase == "" {
		var err error
		passphrase, err = resolvePassphrase(cmd)
		if err != nil {
			return err
		}
	}

	if ctx == "" {
		ctx = ActiveContext()
	}

	path := StorePathForContext(ctx)
	st, err := storeLoad(path, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	for _, e := range st.Entries {
		typed := inferType(e.Value)
		fmt.Fprintf(cmd.OutOrStdout(), "%s = %s (%s)\n", e.Key, e.Value, typed)
	}
	return nil
}

func inferType(value string) string {
	if _, err := strconv.ParseBool(value); err == nil {
		return "bool"
	}
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return "int"
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "float"
	}
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		return "list"
	}
	return "string"
}
