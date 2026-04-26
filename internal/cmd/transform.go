package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "Apply a transformation to one or more env values",
	Long:  transformDoc,
}

var transformUpper = &cobra.Command{
	Use:   "upper [key...]",
	Short: "Convert values to uppercase",
	RunE:  runTransform(strings.ToUpper),
}

var transformLower = &cobra.Command{
	Use:   "lower [key...]",
	Short: "Convert values to lowercase",
	RunE:  runTransform(strings.ToLower),
}

var transformTrim = &cobra.Command{
	Use:   "trim [key...]",
	Short: "Trim whitespace from values",
	RunE:  runTransform(strings.TrimSpace),
}

func init() {
	transformCmd.PersistentFlags().String("passphrase", "", "Passphrase to decrypt the store")
	transformCmd.PersistentFlags().String("context", "", "Context (environment) name")
	transformCmd.AddCommand(transformUpper, transformLower, transformTrim)
	rootCmd.AddCommand(transformCmd)
}

func runTransform(fn func(string) string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("at least one key is required")
		}

		passphrase, _ := cmd.Flags().GetString("passphrase")
		ctxName, _ := cmd.Flags().GetString("context")
		if passphrase == "" {
			passphrase = resolvePassphrase()
		}
		if passphrase == "" {
			return fmt.Errorf("passphrase is required")
		}
		if ctxName == "" {
			ctxName = ActiveContext()
		}

		path := StorePathForContext(ctxName)
		s, err := loadStoreWithPassphrase(path, passphrase)
		if err != nil {
			return fmt.Errorf("failed to load store: %w", err)
		}

		changed := 0
		for _, key := range args {
			val, ok := s.Get(key)
			if !ok {
				return fmt.Errorf("key %q not found", key)
			}
			newVal := fn(val)
			s.Set(key, newVal)
			changed++
		}

		if err := s.Save(path, passphrase); err != nil {
			return fmt.Errorf("failed to save store: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "transformed %d key(s)\n", changed)
		return nil
	}
}
