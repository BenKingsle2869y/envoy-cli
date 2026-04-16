package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func setupNoteEnv(t *testing.T) func() {
	t.Helper()
	tmp := t.TempDir()
	old, _ := os.UserHomeDir()
	t.Setenv("HOME", tmp)
	_ = os.MkdirAll(filepath.Join(tmp, ".envoy"), 0700)
	return func() { t.Setenv("HOME", old) }
}

func execNoteCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"note"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestNoteSetAndGet(t *testing.T) {
	defer setupNoteEnv(t)()
	ctx := ActiveContext()
	notes, _ := LoadNotes(ctx)
	notes["API_KEY"] = "used for payment provider"
	_ = SaveNotes(ctx, notes)

	loaded, err := LoadNotes(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if loaded["API_KEY"] != "used for payment provider" {
		t.Errorf("expected note, got %q", loaded["API_KEY"])
	}
}

func TestNoteClear_RemovesKey(t *testing.T) {
	defer setupNoteEnv(t)()
	ctx := ActiveContext()
	notes := map[string]string{"DB_URL": "primary database"}
	_ = SaveNotes(ctx, notes)

	delete(notes, "DB_URL")
	_ = SaveNotes(ctx, notes)

	loaded, _ := LoadNotes(ctx)
	if _, ok := loaded["DB_URL"]; ok {
		t.Error("expected key to be removed")
	}
}

func TestLoadNotes_EmptyWhenNotExist(t *testing.T) {
	defer setupNoteEnv(t)()
	notes, err := LoadNotes("nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 0 {
		t.Errorf("expected empty map, got %v", notes)
	}
}

func TestNoteGetCmd_FailsWhenMissing(t *testing.T) {
	defer setupNoteEnv(t)()
	_, err := execNoteCmd(t, "get", "MISSING_KEY")
	if err == nil {
		t.Error("expected error for missing key")
	}
}
