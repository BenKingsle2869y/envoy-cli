package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/remote"
	"envoy-cli/internal/store"
)

var pullCmd = &cobra.Command{
	Use:   "pull [environment]",
	Short: "Pull remote .env variables into the local store",
	Args:  cobra.ExactArgs(1),
	RunE:  runPull,
}

func init() {
	pullCmd.Flags().StringP("url", "u", "", "Remote server base URL (required)")
	_ = pullCmd.MarkFlagRequired("url")
	pullCmd.Flags().StringP("token", "t", "", "Bearer token for authentication")
	pullCmd.Flags().Bool("merge", false, "Merge remote vars into existing local env instead of replacing")
	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) error {
	env := args[0]
	baseURL, _ := cmd.Flags().GetString("url")
	token, _ := cmd.Flags().GetString("token")
	merge, _ := cmd.Flags().GetBool("merge")

	passphrase, err := crypto.ResolvePassphrase()
	if err != nil {
		return fmt.Errorf("resolving passphrase: %w", err)
	}

	client := remote.NewClient(baseURL, token, nil)
	vars, err := client.Pull(env)
	if err != nil {
		return fmt.Errorf("pulling environment: %w", err)
	}

	st, err := store.Load(store.DefaultStorePath(), passphrase)
	if err != nil {
		return fmt.Errorf("loading store: %w", err)
	}

	if merge {
		if existing, ok := st.Envs[env]; ok {
			for k, v := range vars {
				existing[k] = v
			}
			vars = existing
		}
	}
	st.Envs[env] = vars

	if err := store.Save(store.DefaultStorePath(), st, passphrase); err != nil {
		return fmt.Errorf("saving store: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Successfully pulled environment %q from %s\n", env, baseURL)
	return nil
}
