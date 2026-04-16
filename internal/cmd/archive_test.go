package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/envoy-cli/envoy/internal/store"
)

func setupArchiveStore(t *testing.T) (string, func()) {
	t.Helper()
	tmp, err := os.MkdirTemp("", "envoy-archivecmd-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", tmp)
	t.Setenv("ENVOY_PASSPHRASE", "testpass")

	ctx := "default"
	path := StorePathForContext(ctx)
	s := store.New()
	s.Entries = []store.Entry{{Key: "KEY1", Value: "val1"}}
	if err := store.Save(path, "testpass", s.Entries); err != nil {
		t.Fatal(err)
	}
	return tmp, func() { os.RemoveAll(tmp) }
}

func execArchiveCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestArchiveListCmd_Empty(t *testing.T) {
	_, cleanup := setupArchiveStore(t)
	defer cleanup()

	out, err := execArchiveCmd(t, "archive", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected some output")
	}
}

func TestArchiveCreateCmd_CreatesArchive(t *testing.T) {
	_, cleanup := setupArchiveStore(t)
	defer cleanup()

	out, err := execArchiveCmd(t, "archive", "create", "--passphrase", "testpass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected confirmation output")
	}

	archives, _ := LoadArchives("default")
	if len(archives) != 1 {
		t.Fatalf("expected 1 archive, got %d", len(archives))
	}
}

func TestArchiveCreateCmd_FailsWithoutPassphrase(t *testing.T) {
	_, cleanup := setupArchiveStore(t)
	defer cleanup()
	os.Unsetenv("ENVOY_PASSPHRASE")

	_, err := execArchiveCmd(t, "archive", "create")
	if err == nil {
		t.Fatal("expected error without passphrase")
	}
}
