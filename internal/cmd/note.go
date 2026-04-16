package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage notes attached to environment keys",
}

var noteSetCmd = &cobra.Command{
	Use:   "set <key> <note>",
	Short: "Attach a note to a key",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runNoteSet,
}

var noteGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Show the note attached to a key",
	Args:  cobra.ExactArgs(1),
	RunE:  runNoteGet,
}

var noteClearCmd = &cobra.Command{
	Use:   "clear <key>",
	Short: "Remove the note attached to a key",
	Args:  cobra.ExactArgs(1),
	RunE:  runNoteClear,
}

func init() {
	noteCmd.AddCommand(noteSetCmd, noteGetCmd, noteClearCmd)
	noteCmd.PersistentFlags().String("passphrase", "", "Passphrase for the store")
	rootCmd.AddCommand(noteCmd)
}

func runNoteSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	note := strings.Join(args[1:], " ")
	pass, _ := cmd.Flags().GetString("passphrase")
	ctx := ActiveContext()
	notes, err := LoadNotes(ctx)
	if err != nil {
		return err
	}
	notes[key] = note
	if err := SaveNotes(ctx, notes); err != nil {
		return err
	}
	_ = pass
	fmt.Fprintf(cmd.OutOrStdout(), "Note set for %q\n", key)
	return nil
}

func runNoteGet(cmd *cobra.Command, args []string) error {
	key := args[0]
	ctx := ActiveContext()
	notes, err := LoadNotes(ctx)
	if err != nil {
		return err
	}
	n, ok := notes[key]
	if !ok {
		return fmt.Errorf("no note found for key %q", key)
	}
	fmt.Fprintln(cmd.OutOrStdout(), n)
	return nil
}

func runNoteClear(cmd *cobra.Command, args []string) error {
	key := args[0]
	ctx := ActiveContext()
	notes, err := LoadNotes(ctx)
	if err != nil {
		return err
	}
	if _, ok := notes[key]; !ok {
		return fmt.Errorf("no note found for key %q", key)
	}
	delete(notes, key)
	if err := SaveNotes(ctx, notes); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Note cleared for %q\n", key)
	return nil
}
