package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage key groups",
}

var groupAddCmd = &cobra.Command{
	Use:   "add <group> <key>",
	Short: "Add a key to a group",
	Args:  cobra.ExactArgs(2),
	RunE:  runGroupAdd,
}

var groupRemoveCmd = &cobra.Command{
	Use:   "remove <group> <key>",
	Short: "Remove a key from a group",
	Args:  cobra.ExactArgs(2),
	RunE:  runGroupRemove,
}

var groupListCmd = &cobra.Command{
	Use:   "list <group>",
	Short: "List keys in a group",
	Args:  cobra.ExactArgs(1),
	RunE:  runGroupList,
}

func init() {
	groupCmd.AddCommand(groupAddCmd, groupRemoveCmd, groupListCmd)
	groupCmd.PersistentFlags().String("passphrase", "", "Passphrase for the store")
	RootCmd.AddCommand(groupCmd)
}

func runGroupAdd(cmd *cobra.Command, args []string) error {
	group, key := args[0], args[1]
	path := groupFilePath(ActiveContext())
	groups, err := LoadGroups(path)
	if err != nil {
		return err
	}
	if err := AddToGroup(groups, group, key); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Added %q to group %q\n", key, group)
	return SaveGroups(path, groups)
}

func runGroupRemove(cmd *cobra.Command, args []string) error {
	group, key := args[0], args[1]
	path := groupFilePath(ActiveContext())
	groups, err := LoadGroups(path)
	if err != nil {
		return err
	}
	if err := RemoveFromGroup(groups, group, key); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Removed %q from group %q\n", key, group)
	return SaveGroups(path, groups)
}

func runGroupList(cmd *cobra.Command, args []string) error {
	group := args[0]
	path := groupFilePath(ActiveContext())
	groups, err := LoadGroups(path)
	if err != nil {
		return err
	}
	keys := groups[group]
	sort.Strings(keys)
	if len(keys) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "No keys in group %q\n", group)
		return nil
	}
	for _, k := range keys {
		fmt.Fprintln(cmd.OutOrStdout(), k)
	}
	return nil
}
