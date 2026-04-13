package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envoy-cli/internal/store"
)

func setupCopyStores(t *testing.T) (srcCtx, dstCtx, passphrase string) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("ENVOY_HOME", dir)

	passphrase = "copy-test-pass"
	srcCtx = "source"
	dstCtx = "dest"

	// create source store with two vars
	srcPath := filepath.Join(dir, srcCtx+".env.enc")
	src := store.Store{Vars: map[string]string{"FOO": "bar", "BAZ": "qux"}}
	if err := store.Save(srcPath, src, passphrase); err != nil {
		t.Fatalf("setup: save src store: %v", err)
	}

	// create destination store with one pre-existing var
	dstPath := filepath.Join(dir, dstCtx+".env.enc")
	dst := store.Store{Vars: map[string]string{"FOO": "original"}}
	if err := store.Save(dstPath, dst, passphrase); err != nil {
		t.Fatalf("setup: save dst store: %v", err)
	}

	return srcCtx, dstCtx, passphrase
}

func TestCopyCmd_SkipsExistingByDefault(t *testing.T) {
	srcCtx, dstCtx, passphrase := setupCopyStores(t)
	os.Setenv("ENVOY_PASSPHRASE", passphrase)
	t.Cleanup(func() { os.Unsetenv("ENVOY_PASSPHRASE") })

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"copy", srcCtx, dstCtx})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("1 skipped")) {
		t.Errorf("expected 1 skipped, got: %s", out)
	}
}

func TestCopyCmd_OverwriteFlag(t *testing.T) {
	srcCtx, dstCtx, passphrase := setupCopyStores(t)
	os.Setenv("ENVOY_PASSPHRASE", passphrase)
	t.Cleanup(func() { os.Unsetenv("ENVOY_PASSPHRASE") })

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"copy", srcCtx, dstCtx, "--overwrite"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("2 variable(s)")) {
		t.Errorf("expected 2 variables copied, got: %s", out)
	}
}

func TestCopyCmd_SameContextFails(t *testing.T) {
	_, _, passphrase := setupCopyStores(t)
	os.Setenv("ENVOY_PASSPHRASE", passphrase)
	t.Cleanup(func() { os.Unsetenv("ENVOY_PASSPHRASE") })

	rootCmd.SetArgs([]string{"copy", "same", "same"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for same source and destination")
	}
}
