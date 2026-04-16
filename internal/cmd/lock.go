package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock the active context to prevent modifications",
	RunE:  runLock,
}

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock the active context to allow modifications",
	RunE:  runUnlock,
}

var lockStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show lock status of the active context",
	RunE:  runLockStatus,
}

func init() {
	lockCmd.AddCommand(lockStatusCmd)
	rootCmd.AddCommand(lockCmd)
	rootCmd.AddCommand(unlockCmd)
}

func lockFilePath() string {
	ctx := ActiveContext()
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", ctx+".lock")
}

func runLock(cmd *cobra.Command, args []string) error {
	path := lockFilePath()
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("context '%s' is already locked", ActiveContext())
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create lock file: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s\n", time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "context '%s' locked\n", ActiveContext())
	return nil
}

func runUnlock(cmd *cobra.Command, args []string) error {
	path := lockFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("context '%s' is not locked", ActiveContext())
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove lock file: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "context '%s' unlocked\n", ActiveContext())
	return nil
}

func runLockStatus(cmd *cobra.Command, args []string) error {
	path := lockFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(cmd.OutOrStdout(), "context '%s' is unlocked\n", ActiveContext())
		return nil
	}
	data, _ := os.ReadFile(path)
	fmt.Fprintf(cmd.OutOrStdout(), "context '%s' is locked (since %s)\n", ActiveContext(), string(data))
	return nil
}

// IsLocked returns true if the given context is currently locked.
func IsLocked(ctx string) bool {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".envoy", ctx+".lock")
	_, err := os.Stat(path)
	return err == nil
}
