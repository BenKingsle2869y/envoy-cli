package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func execEnvCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	root := &cobra.Command{Use: "envoy"}
	envCmd.ResetCommands()
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envUseCmd)
	root.AddCommand(envCmd)
	root.SetOut(&buf)
	envCmd.SetOut(&buf)
	envListCmd.SetOut(&buf)
	envUseCmd.SetOut(&buf)
	root.SetArgs(append([]string{"env"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestEnvListCmd_Empty(t *testing.T) {
	_, cleanup := withTempHome(t)
	defer cleanup()

	out, err := execEnvCmd(t, "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No contexts found") {
		t.Errorf("expected empty message, got: %q", out)
	}
}

func TestEnvListCmd_ShowsContexts(t *testing.T) {
	tmp, cleanup := withTempHome(t)
	defer cleanup()

	dir := filepath.Join(tmp, ".envoy")
	os.MkdirAll(dir, 0700)
	os.WriteFile(filepath.Join(dir, "staging.enc"), []byte("x"), 0600)
	os.WriteFile(filepath.Join(dir, "context"), []byte("staging\n"), 0600)

	out, err := execEnvCmd(t, "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "staging") {
		t.Errorf("expected staging in output, got: %q", out)
	}
	if !strings.Contains(out, "*") {
		t.Errorf("expected active marker '*' in output, got: %q", out)
	}
}

func TestEnvUseCmd_SwitchesContext(t *testing.T) {
	_, cleanup := withTempHome(t)
	defer cleanup()

	out, err := execEnvCmd(t, "use", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "production") {
		t.Errorf("expected production in output, got: %q", out)
	}
	ctx, _ := ActiveContext()
	if ctx != "production" {
		t.Errorf("expected active context to be production, got %q", ctx)
	}
}
