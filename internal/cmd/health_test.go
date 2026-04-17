package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func setupHealthStore(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	passphrase := "health-pass"
	storePath := filepath.Join(dir, ".envoy", "default.env")
	if err := os.MkdirAll(filepath.Dir(storePath), 0700); err != nil {
		t.Fatal(err)
	}
	s, err := loadStoreWithPassphrase(storePath, passphrase)
	if err != nil {
		t.Fatal(err)
	}
	s.Entries = map[string]string{"KEY": "value"}
	if err := saveStoreWithPassphrase(storePath, passphrase, s); err != nil {
		t.Fatal(err)
	}
	return storePath, passphrase
}

func execHealthCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	root := &cobra.Command{Use: "envoy"}
	healthCmd.ResetFlags()
	healthCmd.Flags().String("passphrase", "", "")
	root.AddCommand(healthCmd)
	root.SetOut(buf)
	healthCmd.SetOut(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestHealthCmd_Success(t *testing.T) {
	_, passphrase := setupHealthStore(t)
	_, err := execHealthCmd(t, "health", "--passphrase", passphrase)
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
}

func TestHealthCmd_FailsWithoutPassphrase(t *testing.T) {
	setupHealthStore(t)
	os.Unsetenv("ENVOY_PASSPHRASE")
	_, err := execHealthCmd(t, "health")
	if err == nil {
		t.Fatal("expected error when passphrase missing")
	}
}

func TestHealthCmd_FailsWhenStoreMissing(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	_, err := execHealthCmd(t, "health", "--passphrase", "any")
	if err == nil {
		t.Fatal("expected error when store does not exist")
	}
}
