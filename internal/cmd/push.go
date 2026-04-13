package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/remote"
	"envoy-cli/internal/store"
)

var pushCmd = &cobra.Command{
	Use:   "push [environment]",
	Short: "Push local .env store to a remote server",
	Args:  cobra.ExactArgs(1),
	RunE:  runPush,
}

func init() {
	pushCmd.Flags().StringP("url", "u", "", "Remote server base URL (required)")
	_ = pushCmd.MarkFlagRequired("url")
	pushCmd.Flags().StringP("token", "t", "", "Bearer token for authentication")
	rootCmd.AddCommand(pushCmd)
}

func runPush(cmd *cobra.Command, args []string) error {
	env := args[0]
	baseURL, _ := cmd.Flags().GetString("url")
	token, _ := cmd.Flags().GetString("token")

	passphrase, err := crypto.ResolvePassphrase()
	if err != nil {
		return fmt.Errorf("resolving passphrase: %w", err)
	}

	st, err := store.Load(store.DefaultStorePath(), passphrase)
	if err != nil {
		return fmt.Errorf("loading store: %w", err)
	}

	vars, ok := st.Envs[env]
	if !ok {
		return fmt.Errorf("environment %q not found in local store", env)
	}

	client := remote.NewClient(baseURL, token, nil)
	if err := client.Push(env, vars); err != nil {
		return fmt.Errorf("pushing environment: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Successfully pushed environment %q to %s\n", env, baseURL)
	return nil
}
