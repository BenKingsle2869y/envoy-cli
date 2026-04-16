package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template <file>",
	Short: "Render a template file by substituting env variables from the active store",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplate,
}

func init() {
	templateCmd.Flags().StringP("output", "o", "", "Write rendered output to file instead of stdout")
	templateCmd.Flags().StringP("passphrase", "p", "", "Passphrase to decrypt the store")
	RootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, args []string) error {
	tmplPath := args[0]
	outPath, _ := cmd.Flags().GetString("output")
	passphrase, _ := cmd.Flags().GetString("passphrase")

	if passphrase == "" {
		passphrase = os.Getenv("ENVOY_PASSPHRASE")
	}
	if passphrase == "" {
		return fmt.Errorf("passphrase required: use --passphrase or ENVOY_PASSPHRASE")
	}

	storePath := StorePathForContext(ActiveContext())
	st, err := storeLoad(storePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	result := renderTemplate(string(tmplBytes), st)

	if outPath != "" {
		if err := os.WriteFile(outPath, []byte(result), 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "rendered template written to %s\n", outPath)
	} else {
		fmt.Fprint(cmd.OutOrStdout(), result)
	}
	return nil
}

// renderTemplate replaces {{KEY}} placeholders with values from the store.
func renderTemplate(tmpl string, entries map[string]string) string {
	for k, v := range entries {
		placeholder := "{{" + k + "}}"
		tmpl = strings.ReplaceAll(tmpl, placeholder, v)
	}
	return tmpl
}
