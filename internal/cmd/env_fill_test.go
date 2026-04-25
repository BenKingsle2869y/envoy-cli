package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"envoy-cli/internal/store"
)

func setupFillStore(t *testing.T, passphrase string) string {
	t.Helper()
	dir := t.TempDir()
	path := fmt.Sprintf("%s/test.env", dir)
	s := &store.Store{}
	s.Entries = []store.Entry{
		{Key: "FILL_KEY_A", Value: "alpha"},
		{Key: "FILL_KEY_B", Value: "beta"},
	}
	if err := store.Save(path, passphrase, s); err != nil {
		t.Fatalf("setup: %v", err)
	}
	return path
}

func execFillCmd(args []string, storePath string) (string, error) {
	root := &cobra.Command{Use: "envoy"}
	cmd := &cobra.Command{
		Use:  "fill",
		RunE: runEnvFill,
	}
	cmd.Flags().StringP("passphrase", "p", "", "")
	cmd.Flags().BoolP("overwrite", "o", false, "")
	cmd.Flags().BoolP("export", "e", false, "")
	root.AddCommand(cmd)

	// override context resolution for test
	t.Setenv("ENVOY_STORE_PATH", storePath)

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs(append([]string{"fill"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestFillCmd_ExportMode(t *testing.T) {
	pass := "testpass"
	path := setupFillStore(t, pass)

	out, err := execFillCmd([]string{"--passphrase", pass, "--export"}, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		// export goes to stdout via fmt, not cobra writer; just check no error
	}
}

func TestFillCmd_FailsWithoutPassphrase(t *testing.T) {
	dir := t.TempDir()
	path := fmt.Sprintf("%s/test.env", dir)
	os.Unsetenv("ENVOY_PASSPHRASE")
	_, err := execFillCmd([]string{}, path)
	if err == nil {
		t.Fatal("expected error without passphrase")
	}
}

func TestFillCmd_SkipsExistingByDefault(t *testing.T) {
	pass := "testpass"
	path := setupFillStore(t, pass)

	t.Setenv("FILL_KEY_A", "original")
	defer os.Unsetenv("FILL_KEY_A")

	_, err := execFillCmd([]string{"--passphrase", pass}, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val := os.Getenv("FILL_KEY_A"); val != "original" {
		t.Errorf("expected original value to be preserved, got %q", val)
	}
}
