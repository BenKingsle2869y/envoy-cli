package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func setupTemplateStore(t *testing.T) (storePath, passphrase string) {
	t.Helper()
	dir := t.TempDir()
	passphrase = "template-pass"
	storePath = filepath.Join(dir, "default.env.enc")

	t.Setenv("ENVOY_HOME", dir)

	st := map[string]string{
		"APP_NAME": "envoy",
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
	}
	if err := storeSave(storePath, passphrase, st); err != nil {
		t.Fatalf("storeSave: %v", err)
	}
	return storePath, passphrase
}

func execTemplateCmd(t *testing.T, args []string) (string, error) {
	t.Helper()
	cmd := &cobra.Command{Use: "root"}
	templateCmd.ResetFlags()
	init() // re-register flags
	cmd.AddCommand(templateCmd)

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(append([]string{"template"}, args...))
	err := cmd.Execute()
	return buf.String(), err
}

func TestTemplateCmd_RendersPlaceholders(t *testing.T) {
	_, pass := setupTemplateStore(t)

	tmplFile := filepath.Join(t.TempDir(), "config.tmpl")
	content := "app={{APP_NAME}} env={{APP_ENV}} db={{DB_HOST}}"
	if err := os.WriteFile(tmplFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	out, err := execTemplateCmd(t, []string{tmplFile, "--passphrase", pass})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "app=envoy") {
		t.Errorf("expected APP_NAME substituted, got: %s", out)
	}
	if !strings.Contains(out, "env=production") {
		t.Errorf("expected APP_ENV substituted, got: %s", out)
	}
}

func TestRenderTemplate_LeavesUnknownPlaceholders(t *testing.T) {
	entries := map[string]string{"FOO": "bar"}
	result := renderTemplate("{{FOO}} and {{UNKNOWN}}", entries)
	if result != "bar and {{UNKNOWN}}" {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestTemplateCmd_FailsWithoutPassphrase(t *testing.T) {
	setupTemplateStore(t)
	os.Unsetenv("ENVOY_PASSPHRASE")

	tmplFile := filepath.Join(t.TempDir(), "t.tmpl")
	_ = os.WriteFile(tmplFile, []byte("hello"), 0600)

	_, err := execTemplateCmd(t, []string{tmplFile})
	if err == nil {
		t.Error("expected error when passphrase missing")
	}
}
