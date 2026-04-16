package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var checkpointCmd = &cobra.Command{
	Use:   "checkpoint",
	Short: "Manage named checkpoints for the active context",
}

var checkpointCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a named checkpoint of the current store",
	Args:  cobra.ExactArgs(1),
	RunE:  runCheckpointCreate,
}

var checkpointListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all checkpoints for the active context",
	RunE:  runCheckpointList,
}

var checkpointRestoreCmd = &cobra.Command{
	Use:   "restore <name>",
	Short: "Restore a named checkpoint",
	Args:  cobra.ExactArgs(1),
	RunE:  runCheckpointRestore,
}

func init() {
	checkpointCmd.AddCommand(checkpointCreateCmd)
	checkpointCmd.AddCommand(checkpointListCmd)
	checkpointCmd.AddCommand(checkpointRestoreCmd)
	checkpointCmd.PersistentFlags().String("passphrase", "", "Passphrase for the store")
	rootCmd.AddCommand(checkpointCmd)
}

func runCheckpointCreate(cmd *cobra.Command, args []string) error {
	name := args[0]
	pass, _ := cmd.Flags().GetString("passphrase")
	ctx := ActiveContext()
	store, err := loadStoreWithPassphrase(cmd, pass, ctx)
	if err != nil {
		return err
	}
	if err := CreateCheckpoint(ctx, name, store.Entries); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Checkpoint %q created for context %q\n", name, ctx)
	return nil
}

func runCheckpointList(cmd *cobra.Command, args []string) error {
	ctx := ActiveContext()
	cps, err := LoadCheckpoints(ctx)
	if err != nil {
		return err
	}
	if len(cps) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No checkpoints found.")
		return nil
	}
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tCREATED AT")
	for _, cp := range cps {
		fmt.Fprintf(w, "%s\t%s\n", cp.Name, cp.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
	}
	return w.Flush()
}

func runCheckpointRestore(cmd *cobra.Command, args []string) error {
	name := args[0]
	pass, _ := cmd.Flags().GetString("passphrase")
	ctx := ActiveContext()
	entries, err := RestoreCheckpoint(ctx, name)
	if err != nil {
		return err
	}
	store, err := loadStoreWithPassphrase(cmd, pass, ctx)
	if err != nil {
		return err
	}
	store.Entries = entries
	if err := store.Save(os.Getenv("ENVOY_PASSPHRASE")); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Restored checkpoint %q for context %q\n", name, ctx)
	return nil
}
