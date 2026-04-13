package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func withAuditHome(t *testing.T) (string, func()) {
	t.Helper()
	tmp := t.TempDir()
	orig, _ := os.UserHomeDir()
	t.Setenv("HOME", tmp)
	return tmp, func() { t.Setenv("HOME", orig) }
}

func TestAuditCmd_EmptyLog(t *testing.T) {
	withAuditHome(t)

	buf := &bytes.Buffer{}
	auditCmd.SetOut(buf)
	auditCmd.SetArgs([]string{})

	if err := auditCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No audit log found") {
		t.Errorf("expected empty log message, got: %s", buf.String())
	}
}

func TestAppendAndLoadAuditLog(t *testing.T) {
	tmp, cleanup := withAuditHome(t)
	defer cleanup()

	logPath := filepath.Join(tmp, ".envoy", "default_audit.log")

	actions := []struct{ action, key string }{
		{"set", "FOO"},
		{"set", "BAR"},
		{"unset", "FOO"},
	}

	for _, a := range actions {
		if err := AppendAuditEntry("default", a.action, a.key); err != nil {
			t.Fatalf("AppendAuditEntry: %v", err)
		}
	}

	entries, err := LoadAuditLog(logPath)
	if err != nil {
		t.Fatalf("LoadAuditLog: %v", err)
	}

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(
	if entries[2].Action != "unset" || entries[2].Key != "FOO" {
		t.Errorf("unexpected last entry: %+v", entries[2])
	}

	if entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestAuditCmd_L(t *testing.T) {
	withAuditHome(t)

	for i := 0; i < 5; i++ {
		_endAuditEntry("default", "set", "KEY")
	}

	buf := &bytes.Buffer{}
	auditCmd.SetOut(buf)
	auditCmd.SetArgs([]string{"-n", "2"})

	if err := auditCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 2 entries
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header+2), got %d: %v", len(lines), lines)
	}
}

func TestAuditEntry_TimestampIsUTC(t *testing.T) {
	withAuditHome(t)

	_ = AppendAuditEntry("default", "rotate", "")
	logPath := auditLogPath("default")

	entries, err := LoadAuditLog(logPath)
	if err != nil || len(entries) == 0 {
		t.Fatal("could not load entries")
	}

	if entries[0].Timestamp.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got %v", entries[0].Timestamp.Location())
	}
}
