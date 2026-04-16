package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func setupLockEnv(t *testing.T) func() {
	t.Helper()
	tmpHome := t.TempDir()
	old := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	os.MkdirAll(filepath.Join(tmpHome, ".envoy"), 0700)
	return func() { os.Setenv("HOME", old) }
}

func execLockCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestLockCmd_LocksContext(t *testing.T) {
	cleanup := setupLockEnv(t)
	defer cleanup()

	out, err := execLockCmd(t, "lock")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(out, "locked") {
		t.Errorf("expected 'locked' in output, got: %s", out)
	}
}

func TestLockCmd_FailsWhenAlreadyLocked(t *testing.T) {
	cleanup := setupLockEnv(t)
	defer cleanup()

	execLockCmd(t, "lock")
	_, err := execLockCmd(t, "lock")
	if err == nil {
		t.Fatal("expected error when locking an already locked context")
	}
}

func TestUnlockCmd_UnlocksContext(t *testing.T) {
	cleanup := setupLockEnv(t)
	defer cleanup()

	execLockCmd(t, "lock")
	out, err := execLockCmd(t, "unlock")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(out, "unlocked") {
		t.Errorf("expected 'unlocked' in output, got: %s", out)
	}
}

func TestUnlockCmd_FailsWhenNotLocked(t *testing.T) {
	cleanup := setupLockEnv(t)
	defer cleanup()

	_, err := execLockCmd(t, "unlock")
	if err == nil {
		t.Fatal("expected error when unlocking a non-locked context")
	}
}

func TestIsLocked_ReturnsFalseWhenNotLocked(t *testing.T) {
	cleanup := setupLockEnv(t)
	defer cleanup()

	if IsLocked(ActiveContext()) {
		t.Error("expected context to not be locked")
	}
}

func TestIsLocked_ReturnsTrueWhenLocked(t *testing.T) {
	cleanup := setupLockEnv(t)
	defer cleanup()

	execLockCmd(t, "lock")
	if !IsLocked(ActiveContext()) {
		t.Error("expected context to be locked")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
