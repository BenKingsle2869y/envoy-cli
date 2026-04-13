package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/env"
	"envoy-cli/internal/store"
)

var (
	exportFormat string
	exportOutput string
)

func init() {
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export env variables to a file or stdout",
		Long:  "Export all stored environment variables to a .env file or print them to stdout in the specified format.",
		RunE:  runExport,
	}

	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv", "Output format: dotenv or shell")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path (defaults to stdout)")

	RootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	passphrase, err := crypto.ResolvePassphrase()
	if err != nil {
		return fmt.Errorf("failed to resolve passphrase: %w", err)
	}

	s, err := store.Load(store.DefaultStorePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	var lines []string
	for k, v := range s.Data {
		switch exportFormat {
		case "shell":
			lines = append(lines, fmt.Sprintf("export %s=%q", k, v))
		default:
			lines = append(lines, fmt.Sprintf("%s=%s", k, env.QuoteValue(v)))
		}
	}

	output := strings.Join(lines, "\n") + "\n"

	if exportOutput == "" {
		fmt.Fprint(cmd.OutOrStdout(), output)
		return nil
	}

	if err := os.WriteFile(exportOutput, []byte(output), 0600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Exported %d variable(s) to %s\n", len(lines), exportOutput)
	return nil
}
