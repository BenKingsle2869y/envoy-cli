package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var descriptionCmd = &cobra.Command{
	Use:   "description",
	Short: "Manage key descriptions",
	Long:  "Set, get, or clear human-readable descriptions for environment variable keys.",
}

var descriptionSetCmd = &cobra.Command{
	Use:   "set <key> <description>",
	Short: "Set a description for a key",
	Args:  cobra.ExactArgs(2),
	RunE:  runDescriptionSet,
}

var descriptionGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get the description for a key",
	Args:  cobra.ExactArgs(1),
	RunE:  runDescriptionGet,
}

var descriptionClearCmd = &cobra.Command{
	Use:   "clear <key>",
	Short: "Clear the description for a key",
	Args:  cobra.ExactArgs(1),
	RunE:  runDescriptionClear,
}

func init() {
	descriptionCmd.AddCommand(descriptionSetCmd)
	descriptionCmd.AddCommand(descriptionGetCmd)
	descriptionCmd.AddCommand(descriptionClearCmd)
	RootCmd.AddCommand(descriptionCmd)
}

func runDescriptionSet(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	path := descriptionFilePath(ctx)
	descs, err := LoadDescriptions(path)
	if err != nil {
		return err
	}
	descs[args[0]] = args[1]
	if err := SaveDescriptions(path, descs); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Description set for %q\n", args[0])
	return nil
}

func runDescriptionGet(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	path := descriptionFilePath(ctx)
	descs, err := LoadDescriptions(path)
	if err != nil {
		return err
	}
	v, ok := descs[args[0]]
	if !ok {
		fmt.Fprintf(os.Stderr, "no description found for %q\n", args[0])
		return nil
	}
	fmt.Fprintln(cmd.OutOrStdout(), v)
	return nil
}

func runDescriptionClear(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	path := descriptionFilePath(ctx)
	descs, err := LoadDescriptions(path)
	if err != nil {
		return err
	}
	delete(descs, args[0])
	if err := SaveDescriptions(path, descs); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Description cleared for %q\n", args[0])
	return nil
}
