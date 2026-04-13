package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage environment snapshots",
	Long:  "Create and restore point-in-time snapshots of your environment store.",
}

var snapshotCreateCmd = &cobra.Command{
	Use:   "create [label]",
	Short: "Create a snapshot of the current environment",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSnapshotCreate,
}

var snapshotListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all snapshots",
	RunE:  runSnapshotList,
}

var snapshotRestoreCmd = &cobra.Command{
	Use:   "restore <id>",
	Short: "Restore environment from a snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runSnapshotRestore,
}

func init() {
	snapshotCmd.AddCommand(snapshotCreateCmd)
	snapshotCmd.AddCommand(snapshotListCmd)
	snapshotCmd.AddCommand(snapshotRestoreCmd)
	RootCmd.AddCommand(snapshotCmd)
}

func runSnapshotCreate(cmd *cobra.Command, args []string) error {
	passphrase, err := resolvePassphrase()
	if err != nil {
		return err
	}
	ctx := ActiveContext()
	storePath := StorePathForContext(ctx)
	st, err := store.Load(storePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}
	label := ""
	if len(args) > 0 {
		label = args[0]
	}
	snap, err := CreateSnapshot(storePath, st.Entries(), label, passphrase)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot created: %s (%s)\n", snap.ID, snap.CreatedAt.Format(time.RFC3339))
	return nil
}

func runSnapshotList(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	storePath := StorePathForContext(ctx)
	snaps, err := LoadSnapshots(storePath)
	if err != nil {
		return fmt.Errorf("failed to load snapshots: %w", err)
	}
	if len(snaps) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No snapshots found.")
		return nil
	}
	for _, s := range snaps {
		label := s.Label
		if label == "" {
			label = "(no label)"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\t%d keys\n", s.ID, s.CreatedAt.Format(time.RFC3339), label, len(s.Data))
	}
	return nil
}

func runSnapshotRestore(cmd *cobra.Command, args []string) error {
	passphrase, err := resolvePassphrase()
	if err != nil {
		return err
	}
	ctx := ActiveContext()
	storePath := StorePathForContext(ctx)
	if err := RestoreSnapshot(storePath, args[0], passphrase); err != nil {
		return fmt.Errorf("failed to restore snapshot: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot %s restored successfully.\n", args[0])
	return nil
}
