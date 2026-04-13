package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags on environment variables",
	Long:  "Add, remove, or list tags on stored environment variables for grouping and filtering.",
}

var tagAddCmd = &cobra.Command{
	Use:   "add <key> <tag>",
	Short: "Add a tag to a key",
	Args:  cobra.ExactArgs(2),
	RunE:  runTagAdd,
}

var tagRemoveCmd = &cobra.Command{
	Use:   "remove <key> <tag>",
	Short: "Remove a tag from a key",
	Args:  cobra.ExactArgs(2),
	RunE:  runTagRemove,
}

var tagListCmd = &cobra.Command{
	Use:   "list <key>",
	Short: "List tags on a key",
	Args:  cobra.ExactArgs(1),
	RunE:  runTagList,
}

func init() {
	tagCmd.AddCommand(tagAddCmd, tagRemoveCmd, tagListCmd)
	rootCmd.AddCommand(tagCmd)
}

func runTagAdd(cmd *cobra.Command, args []string) error {
	key, tag := args[0], args[1]
	passphrase, storePath, err := resolveStoreContext(cmd)
	if err != nil {
		return err
	}
	st, err := loadStore(storePath, passphrase)
	if err != nil {
		return err
	}
	if err := AddTag(st.Tags, key, tag); err != nil {
		return err
	}
	if err := saveStore(storePath, passphrase, st); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "tag %q added to %q\n", tag, key)
	return nil
}

func runTagRemove(cmd *cobra.Command, args []string) error {
	key, tag := args[0], args[1]
	passphrase, storePath, err := resolveStoreContext(cmd)
	if err != nil {
		return err
	}
	st, err := loadStore(storePath, passphrase)
	if err != nil {
		return err
	}
	RemoveTag(st.Tags, key, tag)
	if err := saveStore(storePath, passphrase, st); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "tag %q removed from %q\n", tag, key)
	return nil
}

func runTagList(cmd *cobra.Command, args []string) error {
	key := args[0]
	passphrase, storePath, err := resolveStoreContext(cmd)
	if err != nil {
		return err
	}
	st, err := loadStore(storePath, passphrase)
	if err != nil {
		return err
	}
	tags := st.Tags[key]
	if len(tags) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "no tags for %q\n", key)
		return nil
	}
	fmt.Fprintln(cmd.OutOrStdout(), strings.Join(tags, "\n"))
	return nil
}
