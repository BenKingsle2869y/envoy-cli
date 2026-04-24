package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Manage key links between contexts",
	Long:  linkDoc,
}

var linkAddCmd = &cobra.Command{
	Use:   "add <key> <target-context>",
	Short: "Link a key to the same key in another context",
	Args:  cobra.ExactArgs(2),
	RunE:  runLinkAdd,
}

var linkRemoveCmd = &cobra.Command{
	Use:   "remove <key>",
	Short: "Remove a link from a key",
	Args:  cobra.ExactArgs(1),
	RunE:  runLinkRemove,
}

var linkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all key links in the current context",
	Args:  cobra.NoArgs,
	RunE:  runLinkList,
}

func init() {
	linkCmd.AddCommand(linkAddCmd, linkRemoveCmd, linkListCmd)
	linkCmd.PersistentFlags().String("passphrase", "", "Passphrase for the store")
	RootCmd.AddCommand(linkCmd)
}

func runLinkAdd(cmd *cobra.Command, args []string) error {
	key, targetCtx := args[0], args[1]
	ctx := ActiveContext()
	if ctx == targetCtx {
		return fmt.Errorf("cannot link key to the same context")
	}
	path := linkFilePath(ctx)
	links, err := LoadLinks(path)
	if err != nil {
		return err
	}
	if err := AddLink(links, key, targetCtx); err != nil {
		return err
	}
	if err := SaveLinks(path, links); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "linked %s -> %s:%s\n", key, targetCtx, key)
	return nil
}

func runLinkRemove(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	path := linkFilePath(ctx)
	links, err := LoadLinks(path)
	if err != nil {
		return err
	}
	if err := RemoveLink(links, args[0]); err != nil {
		return err
	}
	if err := SaveLinks(path, links); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "removed link for %s\n", args[0])
	return nil
}

func runLinkList(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	links, err := LoadLinks(linkFilePath(ctx))
	if err != nil {
		return err
	}
	if len(links) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no links defined")
		return nil
	}
	for key, target := range links {
		fmt.Fprintf(cmd.OutOrStdout(), "%s -> %s:%s\n", key, target, key)
	}
	return nil
}
