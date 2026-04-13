package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment contexts (e.g. development, staging, production)",
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available environment contexts",
	RunE:  runEnvList,
}

var envUseCmd = &cobra.Command{
	Use:   "use <context>",
	Short: "Switch the active environment context",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvUse,
}

func init() {
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envUseCmd)
	rootCmd.AddCommand(envCmd)
}

func runEnvList(cmd *cobra.Command, args []string) error {
	ctxs, err := ListContexts()
	if err != nil {
		return fmt.Errorf("failed to list contexts: %w", err)
	}
	active, _ := ActiveContext()
	if len(ctxs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No contexts found.")
		return nil
	}
	for _, c := range ctxs {
		marker := "  "
		if strings.EqualFold(c, active) {
			marker = "* "
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s%s\n", marker, c)
	}
	return nil
}

func runEnvUse(cmd *cobra.Command, args []string) error {
	ctx := args[0]
	if err := SetActiveContext(ctx); err != nil {
		return fmt.Errorf("failed to switch context: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Switched to context %q\n", ctx)
	return nil
}

// contextFilePath returns the path to the context config file.
func contextFilePath() string {
	home, _ := os.UserHomeDir()
	return home + "/.envoy/context"
}
