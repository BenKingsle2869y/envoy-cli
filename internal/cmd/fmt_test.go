package cmd_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"

	iCmd "envoy-cli/internal/cmd"
	"envoy-cli/internal/store"
)

func setupFmtStore(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("ENVOY_PASSPHRASE", "fmt-secret")

	ctx := iCmd.ActiveContext()
	path := iCmd.StorePathForContext(ctx)

	s := &store.Store{Vars: map[string]string{
		"ZEBRA": "last",
		"ALPHA": "first",
		"MANGO": "middle",
	}}
	if err := store.Save(path, "fmt-secret", s); err != nil {
		t.Fatalf("setup: save store: %v", err)
	}
	return path, func() { os.RemoveAll(tmpDir) }
}

func execFmtCmd(t *testing.T, extraArgs ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envoy"}
	iCmd.RegisterFmtCmd(root)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	args := append([]string{"fmt"}, extraArgs...)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestFmtCmd_SortsKeys(t *testing.T) {
	path, cleanup := setupFmtStore(t)
	defer cleanup()

	_, err := execFmtCmd(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s, err := store.Load(path, "fmt-secret")
	if err != nil {
		t.Fatalf("load after fmt: %v", err)
	}

	expected := []string{"ALPHA", "MANGO", "ZEBRA"}
	i := 0
	for k := range s.Vars {
		if k != expected[i] {
			t.Errorf("key[%d] = %q, want %q", i, k, expected[i])
		}
		i++
	}
}

func TestFmtCmd_CheckPassesWhenFormatted(t *testing.T) {
	setupFmtStore(t)

	// format first
	_, err := execFmtCmd(t)
	if err != nil {
		t.Fatalf("initial fmt: %v", err)
	}

	// check should pass now
	_, err = execFmtCmd(t, "--check")
	if err != nil {
		t.Errorf("check after fmt should pass, got: %v", err)
	}
}

func TestFmtCmd_MissingPassphrase(t *testing.T) {
	t.Setenv("ENVOY_PASSPHRASE", "")
	_, err := execFmtCmd(t)
	if err == nil {
		t.Error("expected error when passphrase is missing")
	}
}
