package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/store"
)

// injectCmd runs a subprocess with the current context's env vars injected
// into its environment, without writing them to disk or exposing them in shell.
var injectCmd = &cobra.Command{
	Use:   "inject -- <command> [args...]",
	Short: "Run a command with env vars injected from the active context",
	Long: `Inject loads the active context's environment variables and passes them
to the specified subprocess. Variables are injected into the process environment
only — they are never written to disk or printed to stdout.

Example:
  envoy inject -- node server.js
  envoy inject -- python manage.py runserver
  envoy inject --context staging -- ./deploy.sh`,
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: false,
	RunE:               runInject,
}

func init() {
	injectCmd.Flags().String("passphrase", "", "Passphrase to decrypt the store")
	injectCmd.Flags().String("context", "", "Context to inject from (defaults to active context)")
	injectCmd.Flags().Bool("override", false, "Override existing environment variables with store values")
	rootCmd.AddCommand(injectCmd)
}

func runInject(cmd *cobra.Command, args []string) error {
	passphrase, err := crypto.ResolvePassphrase(cmd)
	if err != nil {
		return fmt.Errorf("passphrase required: %w", err)
	}

	ctxName, _ := cmd.Flags().GetString("context")
	if ctxName == "" {
		ctxName = ActiveContext()
	}

	storePath := StorePathForContext(ctxName)
	s, err := store.Load(storePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store for context %q: %w", ctxName, err)
	}

	override, _ := cmd.Flags().GetBool("override")

	// Build the environment for the subprocess.
	// Start with the current process environment.
	baseEnv := os.Environ()
	envMap := make(map[string]string, len(baseEnv))
	for _, e := range baseEnv {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	// Inject store entries, respecting the override flag.
	injected := 0
	for _, entry := range s.Entries {
		if _, exists := envMap[entry.Key]; exists && !override {
			continue
		}
		envMap[entry.Key] = entry.Value
		injected++
	}

	// Reconstruct the environment slice.
	env := make([]string, 0, len(envMap))
	for k, v := range envMap {
		env = append(env, k+"="+v)
	}

	// Locate the binary.
	binary, err := exec.LookPath(args[0])
	if err != nil {
		return fmt.Errorf("command not found: %s", args[0])
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "injecting %d variable(s) into %s\n", injected, args[0])

	subCmd := exec.Command(binary, args[1:]...)
	subCmd.Env = env
	subCmd.Stdin = os.Stdin
	subCmd.Stdout = os.Stdout
	subCmd.Stderr = os.Stderr

	if err := subCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}
