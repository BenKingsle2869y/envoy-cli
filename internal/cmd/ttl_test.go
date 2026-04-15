package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTTLStore(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "envoy-ttl-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	orig, _ := os.UserHomeDir()
	t.Setenv("HOME", tmpDir)
	return filepath.Join(tmpDir, ".envoy", "default.ttl.json"), func() {
		t.Setenv("HOME", orig)
		os.RemoveAll(tmpDir)
	}
}

func TestParseTTLDuration_Hours(t *testing.T) {
	d, err := parseTTLDuration("2h")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 2*time.Hour {
		t.Errorf("expected 2h, got %v", d)
	}
}

func TestParseTTLDuration_Days(t *testing.T) {
	d, err := parseTTLDuration("3d")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 72*time.Hour {
		t.Errorf("expected 72h, got %v", d)
	}
}

func TestParseTTLDuration_Invalid(t *testing.T) {
	_, err := parseTTLDuration("banana")
	if err == nil {
		t.Fatal("expected error for invalid duration")
	}
}

func TestSetAndLoadTTL(t *testing.T) {
	path, cleanup := setupTTLStore(t)
	defer cleanup()

	expiry := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	if err := SetTTL(path, "MY_KEY", expiry); err != nil {
		t.Fatalf("SetTTL: %v", err)
	}

	ttls, err := LoadTTLs(path)
	if err != nil {
		t.Fatalf("LoadTTLs: %v", err)
	}

	got, ok := ttls["MY_KEY"]
	if !ok {
		t.Fatal("expected MY_KEY in TTL map")
	}
	if !got.Equal(expiry) {
		t.Errorf("expected %v, got %v", expiry, got)
	}
}

func TestClearTTL_RemovesKey(t *testing.T) {
	path, cleanup := setupTTLStore(t)
	defer cleanup()

	_ = SetTTL(path, "TEMP_KEY", time.Now().Add(time.Hour))
	if err := ClearTTL(path, "TEMP_KEY"); err != nil {
		t.Fatalf("ClearTTL: %v", err)
	}

	ttls, _ := LoadTTLs(path)
	if _, ok := ttls["TEMP_KEY"]; ok {
		t.Error("expected TEMP_KEY to be removed")
	}
}

func TestExpiredKeys_ReturnsExpired(t *testing.T) {
	ttls := TTLMap{
		"OLD_KEY": time.Now().UTC().Add(-time.Hour),
		"NEW_KEY": time.Now().UTC().Add(time.Hour),
	}
	expired := ExpiredKeys(ttls)
	if len(expired) != 1 || expired[0] != "OLD_KEY" {
		t.Errorf("expected [OLD_KEY], got %v", expired)
	}
}

func TestTTLShowCmd_NoTTLSet(t *testing.T) {
	_, cleanup := setupTTLStore(t)
	defer cleanup()

	cmd := ttlShowCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"MISSING_KEY"})
	if err := cmd.RunE(cmd, []string{"MISSING_KEY"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got == "" {
		t.Error("expected output for missing TTL key")
	}
}
