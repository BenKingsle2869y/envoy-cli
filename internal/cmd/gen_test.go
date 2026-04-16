package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSecret_Length(t *testing.T) {
	v := generateSecret(24, true, false)
	assert.Len(t, v, 24)
}

func TestGenerateSecret_SymbolsIncluded(t *testing.T) {
	found := false
	for i := 0; i < 50; i++ {
		v := generateSecret(32, false, true)
		if strings.ContainsAny(v, charsetSymbols) {
			found = true
			break
		}
	}
	assert.True(t, found, "expected symbols to appear in generated value")
}

func TestGenerateSecret_NoSymbols(t *testing.T) {
	for i := 0; i < 20; i++ {
		v := generateSecret(32, true, false)
		assert.False(t, strings.ContainsAny(v, charsetSymbols), "unexpected symbol in value: %s", v)
	}
}

func TestGenerateSecret_Uniqueness(t *testing.T) {
	a := generateSecret(32, true, true)
	b := generateSecret(32, true, true)
	assert.NotEqual(t, a, b)
}

func setupGenStore(t *testing.T) (string, func()) {
	t.Helper()
	return setupKVStore(t)
}

func execGenCmd(t *testing.T, passphrase string, extraArgs ...string) (string, error) {
	t.Helper()
	args := append([]string{"gen", "--passphrase", passphrase}, extraArgs...)
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestGenCmd_StoresKey(t *testing.T) {
	passphrase, cleanup := setupGenStore(t)
	defer cleanup()

	out, err := execGenCmd(t, passphrase, "MY_SECRET", "--print")
	require.NoError(t, err)
	assert.NotEmpty(t, strings.TrimSpace(out))
	assert.Len(t, strings.TrimSpace(out), 32)
}

func TestGenCmd_SilentByDefault(t *testing.T) {
	passphrase, cleanup := setupGenStore(t)
	defer cleanup()

	out, err := execGenCmd(t, passphrase, "MY_SECRET")
	require.NoError(t, err)
	assert.Contains(t, out, "MY_SECRET")
	assert.NotContains(t, out, charsetAlpha[:5])
}

func TestGenCmd_FailsWithoutPassphrase(t *testing.T) {
	_, cleanup := setupGenStore(t)
	defer cleanup()

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"gen", "MY_SECRET"})
	err := rootCmd.Execute()
	assert.Error(t, err)
}
