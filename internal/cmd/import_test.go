package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(p, []byte(content), 0600))
	return p
}

func TestImportCmd_ImportsVariables(t *testing.T) {
	dir := t.TempDir()
	storePath, passphrase := setupKVStore(t, dir)
	t.Setenv("ENVOY_STORE", storePath)
	t.Setenv("ENVOY_PASSPHRASE", passphrase)

	envFile := writeEnvFile(t, dir, ".env", "FOO=bar\nBAZ=qux\n")

	out, err := executeCommand(RootCmd, "import", envFile)
	require.NoError(t, err)
	assert.Contains(t, out, "Imported 2 variable(s)")
	assert.Contains(t, out, "skipped 0")
}

func TestImportCmd_SkipsExistingKeys(t *testing.T) {
	dir := t.TempDir()
	storePath, passphrase := setupKVStore(t, dir)
	t.Setenv("ENVOY_STORE", storePath)
	t.Setenv("ENVOY_PASSPHRASE", passphrase)

	// Pre-set a key
	_, err := executeCommand(RootCmd, "set", "FOO=original")
	require.NoError(t, err)

	envFile := writeEnvFile(t, dir, ".env", "FOO=overwritten\nNEW=value\n")

	out, err := executeCommand(RootCmd, "import", envFile)
	require.NoError(t, err)
	assert.Contains(t, out, "Imported 1 variable(s)")
	assert.Contains(t, out, "skipped 1")

	// FOO should still be original
	out, err = executeCommand(RootCmd, "get", "FOO")
	require.NoError(t, err)
	assert.Contains(t, out, "original")
}

func TestImportCmd_OverwriteFlag(t *testing.T) {
	dir := t.TempDir()
	storePath, passphrase := setupKVStore(t, dir)
	t.Setenv("ENVOY_STORE", storePath)
	t.Setenv("ENVOY_PASSPHRASE", passphrase)

	_, err := executeCommand(RootCmd, "set", "FOO=original")
	require.NoError(t, err)

	envFile := writeEnvFile(t, dir, ".env", "FOO=overwritten\n")

	out, err := executeCommand(RootCmd, "import", "--overwrite", envFile)
	require.NoError(t, err)
	assert.Contains(t, out, "Imported 1 variable(s)")

	out, err = executeCommand(RootCmd, "get", "FOO")
	require.NoError(t, err)
	assert.Contains(t, out, "overwritten")
}

func TestImportCmd_FailsOnMissingFile(t *testing.T) {
	dir := t.TempDir()
	storePath, passphrase := setupKVStore(t, dir)
	t.Setenv("ENVOY_STORE", storePath)
	t.Setenv("ENVOY_PASSPHRASE", passphrase)

	_, err := executeCommand(RootCmd, "import", "/nonexistent/.env")
	assert.Error(t, err)
}
