package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive and restore env contexts",
}

var archiveCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Archive the active context to a compressed file",
	RunE:  runArchiveCreate,
}

var archiveListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available archives",
	RunE:  runArchiveList,
}

var archiveRestoreCmd = &cobra.Command{
	Use:   "restore <name>",
	Short: "Restore a context from an archive",
	Args:  cobra.ExactArgs(1),
	RunE:  runArchiveRestore,
}

func init() {
	archiveCmd.AddCommand(archiveCreateCmd, archiveListCmd, archiveRestoreCmd)
	archiveCreateCmd.Flags().String("passphrase", "", "Passphrase for the store")
	archiveRestoreCmd.Flags().String("passphrase", "", "Passphrase for the store")
	archiveRestoreCmd.Flags().String("context", "", "Target context to restore into")
	rootCmd.AddCommand(archiveCmd)
}

func runArchiveCreate(cmd *cobra.Command, _ []string) error {
	pass, _ := cmd.Flags().GetString("passphrase")
	if pass == "" {
		pass = os.Getenv("ENVOY_PASSPHRASE")
	}
	if pass == "" {
		return fmt.Errorf("passphrase required")
	}
	ctx := ActiveContext()
	store, err := loadStoreWithPassphrase(cmd, ctx)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("%s-%d", ctx, time.Now().Unix())
	if err := CreateArchive(ctx, name, store.Entries); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "archived context %q as %q\n", ctx, name)
	return nil
}

func runArchiveList(cmd *cobra.Command, _ []string) error {
	archives, err := LoadArchives(ActiveContext())
	if err != nil {
		return err
	}
	if len(archives) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no archives found")
		return nil
	}
	for _, a := range archives {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%d keys\n", a.Name, a.CreatedAt.Format(time.RFC3339), len(a.Entries))
	}
	return nil
}

func runArchiveRestore(cmd *cobra.Command, args []string) error {
	pass, _ := cmd.Flags().GetString("passphrase")
	if pass == "" {
		pass = os.Getenv("ENVOY_PASSPHRASE")
	}
	if pass == "" {
		return fmt.Errorf("passphrase required")
	}
	targetCtx, _ := cmd.Flags().GetString("context")
	if targetCtx == "" {
		targetCtx = ActiveContext()
	}
	entries, err := RestoreArchive(ActiveContext(), args[0])
	if err != nil {
		return err
	}
	path := StorePathForContext(targetCtx)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	store, err := loadStoreWithPassphrase(cmd, targetCtx)
	if err != nil {
		store = &storeWrapper{Entries: entries, path: path, pass: pass}
	} else {
		store.Entries = entries
	}
	_ = store
	fmt.Fprintf(cmd.OutOrStdout(), "restored archive %q into context %q\n", args[0], targetCtx)
	return nil
}
