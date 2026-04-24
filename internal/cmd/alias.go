package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage key aliases within a context",
}

var aliasAddCmd = &cobra.Command{
	Use:   "add <alias> <target-key>",
	Short: "Create an alias pointing to an existing key",
	Args:  cobra.ExactArgs(2),
	RunE:  runAliasAdd,
}

var aliasRemoveCmd = &cobra.Command{
	Use:   "remove <alias>",
	Short: "Remove an alias",
	Args:  cobra.ExactArgs(1),
	RunE:  runAliasRemove,
}

var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases in the current context",
	Args:  cobra.NoArgs,
	RunE:  runAliasList,
}

func init() {
	aliasCmd.AddCommand(aliasAddCmd, aliasRemoveCmd, aliasListCmd)
	RootCmd.AddCommand(aliasCmd)
}

func runAliasAdd(cmd *cobra.Command, args []string) error {
	alias, target := args[0], args[1]
	path := aliasFilePath(ActiveContext())
	aliases, err := LoadAliases(path)
	if err != nil {
		return err
	}
	if err := AddAlias(aliases, alias, target); err != nil {
		return err
	}
	if err := SaveAliases(path, aliases); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "alias %q -> %q added\n", alias, target)
	return nil
}

func runAliasRemove(cmd *cobra.Command, args []string) error {
	alias := args[0]
	path := aliasFilePath(ActiveContext())
	aliases, err := LoadAliases(path)
	if err != nil {
		return err
	}
	if err := RemoveAlias(aliases, alias); err != nil {
		return err
	}
	if err := SaveAliases(path, aliases); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "alias %q removed\n", alias)
	return nil
}

func runAliasList(cmd *cobra.Command, args []string) error {
	path := aliasFilePath(ActiveContext())
	aliases, err := LoadAliases(path)
	if err != nil {
		return err
	}
	if len(aliases) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no aliases defined")
		return nil
	}
	for alias, target := range aliases {
		fmt.Fprintf(cmd.OutOrStdout(), "%s -> %s\n", alias, target)
	}
	return nil
}
