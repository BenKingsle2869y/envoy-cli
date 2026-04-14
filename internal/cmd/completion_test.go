package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func execCompletionCmd(t *testing.T, shell string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"completion", shell})
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestCompletionCmd_Bash(t *testing.T) {
	out, err := execCompletionCmd(t, "bash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "bash") && len(out) == 0 {
		t.Error("expected non-empty bash completion output")
	}
}

func TestCompletionCmd_Zsh(t *testing.T) {
	out, err := execCompletionCmd(t, "zsh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) == 0 {
		t.Error("expected non-empty zsh completion output")
	}
}

func TestCompletionCmd_Fish(t *testing.T) {
	out, err := execCompletionCmd(t, "fish")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) == 0 {
		t.Error("expected non-empty fish completion output")
	}
}

func TestCompletionCmd_PowerShell(t *testing.T) {
	out, err := execCompletionCmd(t, "powershell")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) == 0 {
		t.Error("expected non-empty powershell completion output")
	}
}

func TestCompletionCmd_InvalidShell(t *testing.T) {
	rootCmd.SetArgs([]string{"completion", "unknownshell"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for invalid shell argument")
	}
}

func TestCompletionCmd_NoArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"completion"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no shell argument is provided")
	}
}
