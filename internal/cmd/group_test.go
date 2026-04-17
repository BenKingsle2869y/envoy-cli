package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func setupGroupEnv(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	old, _ := os.UserHomeDir()
	os.Setenv("HOME", dir)
	return dir, func() { os.Setenv("HOME", old) }
}

func execGroupCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envoy"}
	root.AddCommand(groupCmd)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestGroupAddCmd_AddsKey(t *testing.T) {
	dir, cleanup := setupGroupEnv(t)
	defer cleanup()
	_ = os.MkdirAll(filepath.Join(dir, ".envoy"), 0700)

	out, err := execGroupCmd(t, "group", "add", "backend", "DB_URL")
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Errorf("expected output")
	}
}

func TestGroupListCmd_ShowsKeys(t *testing.T) {
	dir, cleanup := setupGroupEnv(t)
	defer cleanup()
	_ = os.MkdirAll(filepath.Join(dir, ".envoy"), 0700)

	path := groupFilePath(ActiveContext())
	groups := map[string][]string{"frontend": {"API_URL", "APP_NAME"}}
	if err := SaveGroups(path, groups); err != nil {
		t.Fatal(err)
	}

	out, err := execGroupCmd(t, "group", "list", "frontend")
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Errorf("expected keys in output")
	}
}

func TestGroupListCmd_EmptyGroup(t *testing.T) {
	_, cleanup := setupGroupEnv(t)
	defer cleanup()

	out, err := execGroupCmd(t, "group", "list", "nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Errorf("expected empty message")
	}
}

func TestGroupRemoveCmd_RemovesKey(t *testing.T) {
	dir, cleanup := setupGroupEnv(t)
	defer cleanup()
	_ = os.MkdirAll(filepath.Join(dir, ".envoy"), 0700)

	path := groupFilePath(ActiveContext())
	groups := map[string][]string{"ops": {"SECRET_KEY"}}
	_ = SaveGroups(path, groups)

	_, err := execGroupCmd(t, "group", "remove", "ops", "SECRET_KEY")
	if err != nil {
		t.Fatal(err)
	}
}
