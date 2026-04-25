package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/cmd"
	"envoy-cli/internal/crypto"
	"envoy-cli/internal/store"
)

var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve variable references within the store",
	Long:  resolveDoc,
	RunE:  runResolve,
}

func init() {
	resolveCmd.Flags().String("passphrase", "", "Passphrase to decrypt the store")
	resolveCmd.Flags().Bool("in-place", false, "Write resolved values back to the store")
	rootCmd.AddCommand(resolveCmd)
}

func runResolve(cmd *cobra.Command, args []string) error {
	passphrase, err := crypto.ResolvePassphrase(cmd)
	if err != nil {
		return err
	}

	ctx := ActiveContext()
	path := StorePathForContext(ctx)

	s, err := store.Load(path, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	inPlace, _ := cmd.Flags().GetBool("in-place")
	resolved, changed := resolveReferences(s.Entries)

	if len(changed) == 0 {
		fmt.Println("No variable references found.")
		return nil
	}

	for k, v := range changed {
		fmt.Printf("%s = %s\n", k, v)
	}

	if inPlace {
		for k, v := range changed {
			s.Entries[k] = store.Entry{Value: v}
		}
		if err := store.Save(path, passphrase, s); err != nil {
			return fmt.Errorf("failed to save store: %w", err)
		}
		fmt.Printf("Resolved %d reference(s) in-place.\n", len(changed))
	}

	return nil
}
