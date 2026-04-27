package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func setupEnvExportStore(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("ENVOY_PASSPHRASE", "export-pass")

	ctxFile := filepath.Join(tmpDir, ".envoy_context")
	if err := os.WriteFile(ctxFile, []byte("default"), 0600); err != nil {
		t.Fatal(err)
	}
	return tmpDir, func() {}
}

func execEnvExportCmd(t *testing.T, extraArgs ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envoy"}
	root.AddCommand(envCmd)

	var buf bytes.Buffer
	root.SetOut(&buf)
	envExportCmd.ResetFlags()
	envExportCmd.Flags().StringVar(&envExportContext, "context", "", "")
	envExportCmd.Flags().BoolVar(&envExportNoBlanks, "no-blanks", false, "")

	args := append([]string{"env", "export"}, extraArgs...)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestEnvExportCmd_OutputsExportStatements(t *testing.T) {
	setupEnvExportStore(t)

	path := StorePathForContext(ActiveContext())
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}

	store := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	if err := saveTestStore(t, path, "export-pass", store); err != nil {
		t.Fatal(err)
	}

	out, err := execEnvExportCmd(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP_ENV=") {
		t.Errorf("expected APP_ENV in output, got: %s", out)
	}
	if !strings.Contains(out, "export PORT=") {
		t.Errorf("expected PORT in output, got: %s", out)
	}
}

func TestEnvExportCmd_NoBlanksFlag(t *testing.T) {
	setupEnvExportStore(t)

	path := StorePathForContext(ActiveContext())
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}

	store := map[string]string{"FILLED": "yes", "EMPTY": ""}
	if err := saveTestStore(t, path, "export-pass", store); err != nil {
		t.Fatal(err)
	}

	out, err := execEnvExportCmd(t, "--no-blanks")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "EMPTY") {
		t.Errorf("expected EMPTY to be skipped, got: %s", out)
	}
	if !strings.Contains(out, "FILLED") {
		t.Errorf("expected FILLED in output, got: %s", out)
	}
}

func TestEnvExportCmd_FailsWithoutPassphrase(t *testing.T) {
	setupEnvExportStore(t)
	os.Unsetenv("ENVOY_PASSPHRASE")

	_, err := execEnvExportCmd(t)
	if err == nil {
		t.Fatal("expected error when passphrase is missing")
	}
}

func saveTestStore(t *testing.T, path, passphrase string, data map[string]string) error {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	s, err := loadStoreWithPassphrase(path, passphrase)
	if err != nil {
		s = make(map[string]string)
	}
	for k, v := range data {
		s[k] = v
	}
	_ = s
	return nil
}
