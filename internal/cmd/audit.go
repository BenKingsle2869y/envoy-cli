package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// AuditEntry records a single operation performed on the store.
type AuditEntry struct {
	Timestamp time.Time
	Context   string
	Action    string
	Key       string
}

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Show audit log of recent changes to the active environment store",
	RunE:  runAudit,
}

func init() {
	auditCmd.Flags().IntP("limit", "n", 20, "Maximum number of entries to show")
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, _ []string) error {
	limit, _ := cmd.Flags().GetInt("limit")

	ctx := ActiveContext()
	logPath := auditLogPath(ctx)

	entries, err := LoadAuditLog(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(cmd.OutOrStdout(), "No audit log found for context:", ctx)
			return nil
		}
		return fmt.Errorf("failed to read audit log: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Audit log is empty.")
		return nil
	}

	if limit > 0 && len(entries) > limit {
		entries = entries[len(entries)-limit:]
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tCONTEXT\tACTION\tKEY")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.Timestamp.Format(time.RFC3339),
			e.Context,
			e.Action,
			e.Key,
		)
	}
	return w.Flush()
}
