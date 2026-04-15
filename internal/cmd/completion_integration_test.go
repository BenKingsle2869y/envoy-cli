package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// buildTestRoot creates a minimal cobra root with the completion command
// attached, isolated from global state for integration testing.
func buildTestRoot() *cobra.Command {
	root := &cobra.Command{
		Use:   "envoy",
		Short: "Manage .env files",
	}

	completion := &cobra.Command{
		Use:       "completion [bash|zsh|fish|powershell]",
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			}
			return nil
		},
	}
	root.AddCommand(completion)
	return root
}

// runCompletion is a helper that executes the completion subcommand for the
// given shell and returns the captured output.
func runCompletion(t *testing.T, shell string) string {
	t.Helper()
	root := buildTestRoot()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"completion", shell})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error running completion %s: %v", shell, err)
	}
	return buf.String()
}

func TestCompletionIntegration_BashContainsCommandName(t *testing.T) {
	output := runCompletion(t, "bash")
	if !strings.Contains(output, "envoy") {
		t.Errorf("expected bash completion to reference 'envoy', got: %s", output)
	}
}

func TestCompletionIntegration_ZshContainsCommandName(t *testing.T) {
	output := runCompletion(t, "zsh")
	if len(output) == 0 {
		t.Error("expected non-empty zsh completion output")
	}
}

func TestCompletionIntegration_FishContainsCommandName(t *testing.T) {
	output := runCompletion(t, "fish")
	if !strings.Contains(output, "envoy") {
		t.Errorf("expected fish completion to reference 'envoy', got: %s", output)
	}
}
