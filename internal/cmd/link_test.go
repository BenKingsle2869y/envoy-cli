package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func setupLinkEnv(t *testing.T) (string, func()) {
	t.Helper()
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	os.MkdirAll(filepath.Join(tmpHome, ".envoy"), 0700)
	return tmpHome, func() { os.Setenv("HOME", origHome) }
}

func execLinkCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envoy"}
	root.AddCommand(linkCmd)
	buf := &bytes.Buffer{}
	linkCmd.ResetFlags()
	linkAddCmd.ResetFlags()
	linkRemoveCmd.ResetFlags()
	linkListCmd.ResetFlags()
	root.SetOut(buf)
	linkCmd.SetOut(buf)
	linkAddCmd.SetOut(buf)
	linkRemoveCmd.SetOut(buf)
	linkListCmd.SetOut(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestLinkAddCmd_AddsLink(t *testing.T) {
	_, cleanup := setupLinkEnv(t)
	defer cleanup()

	out, err := execLinkCmd(t, "link", "add", "DB_URL", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "linked DB_URL -> production:DB_URL") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestLinkListCmd_Empty(t *testing.T) {
	_, cleanup := setupLinkEnv(t)
	defer cleanup()

	out, err := execLinkCmd(t, "link", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no links defined") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestLinkRemoveCmd_RemovesLink(t *testing.T) {
	_, cleanup := setupLinkEnv(t)
	defer cleanup()

	ctx := ActiveContext()
	path := linkFilePath(ctx)
	links := Links{"API_KEY": "production"}
	if err := SaveLinks(path, links); err != nil {
		t.Fatalf("setup: %v", err)
	}

	out, err := execLinkCmd(t, "link", "remove", "API_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "removed link for API_KEY") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestLinkAddCmd_SameContextFails(t *testing.T) {
	_, cleanup := setupLinkEnv(t)
	defer cleanup()

	ctx := ActiveContext()
	_, err := execLinkCmd(t, "link", "add", "DB_URL", ctx)
	if err == nil {
		t.Error("expected error when linking to same context")
	}
}
