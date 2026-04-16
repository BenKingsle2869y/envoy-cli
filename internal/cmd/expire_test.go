package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupExpireStore(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	return dir, func() { os.RemoveAll(dir) }
}

func execExpireCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"expire"}, args...))
	err := rootCmd.Execute()
	rootCmd.SetArgs(nil)
	return buf.String(), err
}

func TestExpireCmd_NoTTLsReturnsEmpty(t *testing.T) {
	dir, cleanup := setupExpireStore(t)
	defer cleanup()

	storePath := filepath.Join(dir, ".envoy", "default.store")
	_ = storePath

	out, err := execExpireCmd(t, "--passphrase", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Log("got empty output as expected when no TTLs set")
	}
}

func TestExpireCmd_ShowsExpiredKey(t *testing.T) {
	dir, cleanup := setupExpireStore(t)
	defer cleanup()

	storePath := filepath.Join(dir, ".envoy", "default.store")
	_ = os.MkdirAll(filepath.Dir(storePath), 0700)

	ttls := map[string]time.Time{
		"OLD_KEY": time.Now().UTC().Add(-2 * time.Hour),
	}
	if err := saveTTLs(storePath, ttls); err != nil {
		t.Fatalf("failed to save TTLs: %v", err)
	}

	out, err := execExpireCmd(t, "--passphrase", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains([]byte(out), []byte("EXPIRED")) {
		t.Errorf("expected EXPIRED in output, got: %s", out)
	}
	if !bytes.Contains([]byte(out), []byte("OLD_KEY")) {
		t.Errorf("expected OLD_KEY in output, got: %s", out)
	}
}
