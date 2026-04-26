package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTransformStore(t *testing.T, entries map[string]string) (string, string) {
	t.Helper()
	dir := t.TempDir()
	path, passphrase := dir+"/test.env.enc", "transform-pass"
	s := newTestStore(t, path, passphrase)
	for k, v := range entries {
		s.Set(k, v)
	}
	require.NoError(t, s.Save(path, passphrase))
	return path, passphrase
}

func execTransformCmd(t *testing.T, path, passphrase string, args ...string) (string, error) {
	t.Helper()
	root := newTestRoot()
	var buf bytes.Buffer
	root.SetOut(&buf)
	fullArgs := append([]string{"transform"}, args...)
	fullArgs = append(fullArgs,
		"--passphrase", passphrase,
		"--context", contextNameFromPath(path),
	)
	root.SetArgs(fullArgs)
	err := root.Execute()
	return buf.String(), err
}

func TestTransformCmd_Upper(t *testing.T) {
	path, pass := setupTransformStore(t, map[string]string{"APP_ENV": "production"})
	out, err := execTransformCmd(t, path, pass, "upper", "APP_ENV")
	require.NoError(t, err)
	assert.Contains(t, out, "transformed 1 key(s)")

	s := loadForTest(t, path, pass)
	val, ok := s.Get("APP_ENV")
	require.True(t, ok)
	assert.Equal(t, "PRODUCTION", val)
}

func TestTransformCmd_Lower(t *testing.T) {
	path, pass := setupTransformStore(t, map[string]string{"DB_HOST": "LOCALHOST"})
	out, err := execTransformCmd(t, path, pass, "lower", "DB_HOST")
	require.NoError(t, err)
	assert.Contains(t, out, "transformed 1 key(s)")

	s := loadForTest(t, path, pass)
	val, ok := s.Get("DB_HOST")
	require.True(t, ok)
	assert.Equal(t, "localhost", val)
}

func TestTransformCmd_Trim(t *testing.T) {
	path, pass := setupTransformStore(t, map[string]string{"SECRET": "  abc  "})
	out, err := execTransformCmd(t, path, pass, "trim", "SECRET")
	require.NoError(t, err)
	assert.Contains(t, out, "transformed 1 key(s)")

	s := loadForTest(t, path, pass)
	val, ok := s.Get("SECRET")
	require.True(t, ok)
	assert.Equal(t, "abc", val)
}

func TestTransformCmd_MissingKey(t *testing.T) {
	path, pass := setupTransformStore(t, map[string]string{})
	_, err := execTransformCmd(t, path, pass, "upper", "NONEXISTENT")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestTransformCmd_NoArgs(t *testing.T) {
	path, pass := setupTransformStore(t, map[string]string{})
	_, err := execTransformCmd(t, path, pass, "upper")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one key")
}
