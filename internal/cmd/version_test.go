package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCmd_DefaultValues(t *testing.T) {
	// Reset to known defaults
	origVersion, origCommit, origBuild := Version, Commit, BuildDate
	defer func() {
		Version = origVersion
		Commit = origCommit
		BuildDate = origBuild
	}()

	Version = "dev"
	Commit = "none"
	BuildDate = "unknown"

	buf := &bytes.Buffer{}
	versionCmd.SetOut(buf)
	versionCmd.SetErr(buf)
	versionCmd.Run(versionCmd, []string{})

	out := buf.String()
	if !strings.Contains(out, "envoy-cli version dev") {
		t.Errorf("expected version line, got: %s", out)
	}
	if !strings.Contains(out, "commit:     none") {
		t.Errorf("expected commit line, got: %s", out)
	}
	if !strings.Contains(out, "build date: unknown") {
		t.Errorf("expected build date line, got: %s", out)
	}
}

func TestVersionCmd_CustomValues(t *testing.T) {
	origVersion, origCommit, origBuild := Version, Commit, BuildDate
	defer func() {
		Version = origVersion
		Commit = origCommit
		BuildDate = origBuild
	}()

	Version = "1.2.3"
	Commit = "abc1234"
	BuildDate = "2024-06-01"

	buf := &bytes.Buffer{}
	versionCmd.SetOut(buf)
	versionCmd.SetErr(buf)
	versionCmd.Run(versionCmd, []string{})

	out := buf.String()
	if !strings.Contains(out, "envoy-cli version 1.2.3") {
		t.Errorf("expected version 1.2.3, got: %s", out)
	}
	if !strings.Contains(out, "commit:     abc1234") {
		t.Errorf("expected commit abc1234, got: %s", out)
	}
	if !strings.Contains(out, "build date: 2024-06-01") {
		t.Errorf("expected build date 2024-06-01, got: %s", out)
	}
}
