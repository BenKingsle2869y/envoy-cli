package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupEnvDefaultStore(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	passphrase := "test-passphrase"
	t.Setenv("ENVOY_PASSPHRASE", passphrase)
	return dir, passphrase
}

func execEnvDefaultCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	envDefaultCmd.SetOut(buf)
	envDefaultCmd.SetErr(buf)
	envDefaultCmd.SetArgs(args)
	err := envDefaultCmd.Execute()
	envDefaultCmd.SetArgs(nil)
	return buf.String(), err
}

func TestEnvDefaultCmd_SetsWhenMissing(t *testing.T) {
	_, pass := setupEnvDefaultStore(t)

	// initialise a store first
	err := runInitWithPassphrase(pass)
	require.NoError(t, err)

	out, err := execEnvDefaultCmd("MY_KEY", "--value", "fallback", "--passphrase", pass)
	require.NoError(t, err)
	assert.Contains(t, out, "default set")
	assert.Contains(t, out, "MY_KEY")
}

func TestEnvDefaultCmd_SkipsWhenPresent(t *testing.T) {
	_, pass := setupEnvDefaultStore(t)

	err := runInitWithPassphrase(pass)
	require.NoError(t, err)

	// pre-set the key
	_, err = execEnvDefaultCmd("MY_KEY", "--value", "original", "--passphrase", pass)
	require.NoError(t, err)

	// attempt to default again
	out, err := execEnvDefaultCmd("MY_KEY", "--value", "override", "--passphrase", pass)
	require.NoError(t, err)
	assert.Contains(t, out, "already set")
}

func TestEnvDefaultCmd_FailsWithoutPassphrase(t *testing.T) {
	setupEnvDefaultStore(t)
	t.Setenv("ENVOY_PASSPHRASE", "")

	_, err := execEnvDefaultCmd("MY_KEY", "--value", "val")
	assert.Error(t, err)
}
