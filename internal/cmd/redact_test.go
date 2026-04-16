package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func setupRedactStore(t *testing.T, passphrase string) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	path := StorePathForContext("default")
	store, err := createOrLoadStore(path, passphrase)
	if err != nil {
		t.Fatalf("setup store: %v", err)
	}
	store.Entries = append(store.Entries,
		newEntry("APP_NAME", "myapp"),
		newEntry("DB_PASSWORD", "s3cr3t"),
		newEntry("API_KEY", "abc123"),
		newEntry("PORT", "8080"),
	)
	if err := saveStore(path, passphrase, store); err != nil {
		t.Fatalf("save store: %v", err)
	}
	return dir
}

func execRedactCmd(t *testing.T, passphrase string, extraArgs ...string) string {
	t.Helper()
	cmd := &cobra.Command{Use: "root"}
	redact := *redactCmd
	cmd.AddCommand(&redact)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	redact.SetOut(&buf)
	args := append([]string{"redact", "--passphrase", passphrase}, extraArgs...)
	cmd.SetArgs(args)
	_ = cmd.Execute()
	return buf.String()
}

func TestRedactCmd_MasksSecretKeys(t *testing.T) {
	passphrase := "redactpass"
	setupRedactStore(t, passphrase)

	out := execRedactCmd(t, passphrase)

	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME to be unmasked, got: %s", out)
	}
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("expected DB_PASSWORD value to be masked, got: %s", out)
	}
	if strings.Contains(out, "abc123") {
		t.Errorf("expected API_KEY value to be masked, got: %s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT to be unmasked, got: %s", out)
	}
}

func TestLooksSecret_DetectsKeywords(t *testing.T) {
	cases := []struct {
		key    string
		want   bool
	}{
		{"DB_PASSWORD", true},
		{"AUTH_TOKEN", true},
		{"API_KEY", true},
		{"APP_NAME", false},
		{"PORT", false},
		{"PRIVATE_KEY_PATH", true},
	}
	for _, tc := range cases {
		if got := looksSecret(tc.key); got != tc.want {
			t.Errorf("looksSecret(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}
