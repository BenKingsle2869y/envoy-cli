package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"envoy-cli/internal/crypto"
	"envoy-cli/internal/store"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environment variables in the store",
	Aliases: []string{"ls"},
	RunE:  runList,
}

func init() {
	listCmd.Flags().StringP("passphrase", "p", "", "Passphrase to decrypt the store")
	listCmd.Flags().StringP("file", "f", store.DefaultStorePath, "Path to the encrypted store file")
	listCmd.Flags().BoolP("keys-only", "k", false, "Print only key names without values")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	passphrase, err := crypto.ResolvePassphrase(cmd)
	if err != nil {
		return fmt.Errorf("passphrase error: %w", err)
	}

	filePath, _ := cmd.Flags().GetString("file")
	keysOnly, _ := cmd.Flags().GetBool("keys-only")

	s, err := store.Load(filePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	if len(s.Env) == 0 {
		fmt.Println("Store is empty.")
		return nil
	}

	keys := make([]string, 0, len(s.Env))
	for k := range s.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if keysOnly {
			fmt.Println(k)
		} else {
			fmt.Printf("%s=%s\n", k, s.Env[k])
		}
	}
	return nil
}
