package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenameCmd_Success(t *testing.T) {
	st, storePath, passphrase := setupKVStore(t)
	st.Entries["OLD_KEY"] = "hello"
	err := saveStore(storePath, passphrase, st)
	require.NoError(t, err)

	t.Setenv("ENVOY_PASSPHRASE", passphrase)
	t.Setenv("ENVOY_STORE", storePath)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"rename", "OLD_KEY", "NEW_KEY"})

	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "OLD_KEY")
	assert.Contains(t, buf.String(), "NEW_KEY")

	reloaded, err := loadStore(storePath, passphrase)
	require.NoError(t, err)
	_, oldExists := reloaded.Entries["OLD_KEY"]
	val, newExists := reloaded.Entries["NEW_KEY"]
	assert.False(t, oldExists, "old key should be removed")
	assert.True(t, newExists, "new key should exist")
	assert.Equal(t, "hello", val)
}

func TestRenameCmd_SameKeyFails(t *testing.T) {
	_, storePath, passphrase := setupKVStore(t)
	t.Setenv("ENVOY_PASSPHRASE", passphrase)
	t.Setenv("ENVOY_STORE", storePath)

	rootCmd.SetArgs([]string{"rename", "KEY", "KEY"})
	err := rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "same")
}

func TestRenameCmd_MissingKeyFails(t *testing.T) {
	_, storePath, passphrase := setupKVStore(t)
	t.Setenv("ENVOY_PASSPHRASE", passphrase)
	t.Setenv("ENVOY_STORE", storePath)

	rootCmd.SetArgs([]string{"rename", "DOES_NOT_EXIST", "NEW_KEY"})
	err := rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRenameCmd_DestinationExistsFails(t *testing.T) {
	st, storePath, passphrase := setupKVStore(t)
	st.Entries["KEY_A"] = "val_a"
	st.Entries["KEY_B"] = "val_b"
	err := saveStore(storePath, passphrase, st)
	require.NoError(t, err)

	t.Setenv("ENVOY_PASSPHRASE", passphrase)
	t.Setenv("ENVOY_STORE", storePath)

	rootCmd.SetArgs([]string{"rename", "KEY_A", "KEY_B"})
	err = rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}
