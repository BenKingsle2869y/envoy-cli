package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envoy-cli/internal/store"
)

func TestExportCmd_PrintsToStdout(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "store.env")
	passphrase := "export-test-pass"

	t.Setenv("ENVOY_STORE_PATH", path)
	t.Setenv("ENVOY_PASSPHRASE", passphrase)

	s, err := store.Load(path, passphrase)
	if err != nil {
		t.Fatalf("failed to load store: %v", err)
	}
	s.Data["FOO"] = "bar"
	s.Data["BAZ"] = "qux"
	if err := store.Save(path, passphrase, s); err != nil {
		t.Fatalf("failed to save store: %v", err)
	}

	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetArgs([]string{"export"})

	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("export command failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", out)
	}
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got: %s", out)
	}
}

func TestExportCmd_WritesToFile(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "store.env")
	outPath := filepath.Join(dir, "output.env")
	passphrase := "export-file-pass"

	t.Setenv("ENVOY_STORE_PATH", storePath)
	t.Setenv("ENVOY_PASSPHRASE", passphrase)

	s, err := store.Load(storePath, passphrase)
	if err != nil {
		t.Fatalf("failed to load store: %v", err)
	}
	s.Data["KEY"] = "value"
	if err := store.Save(storePath, passphrase, s); err != nil {
		t.Fatalf("failed to save store: %v", err)
	}

	RootCmd.SetArgs([]string{"export", "--output", outPath})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("export command failed: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !strings.Contains(string(data), "KEY=value") {
		t.Errorf("expected KEY=value in file, got: %s", string(data))
	}
}

func TestExportCmd_ShellFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "store.env")
	passphrase := "export-shell-pass"

	t.Setenv("ENVOY_STORE_PATH", path)
	t.Setenv("ENVOY_PASSPHRASE", passphrase)

	s, err := store.Load(path, passphrase)
	if err != nil {
		t.Fatalf("failed to load store: %v", err)
	}
	s.Data["MY_VAR"] = "hello"
	if err := store.Save(path, passphrase, s); err != nil {
		t.Fatalf("failed to save store: %v", err)
	}

	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetArgs([]string{"export", "--format", "shell"})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("export command failed: %v", err)
	}

	if !strings.Contains(buf.String(), "export MY_VAR=") {
		t.Errorf("expected shell export syntax, got: %s", buf.String())
	}
}
