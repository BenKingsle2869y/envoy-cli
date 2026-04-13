package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddTag_AddsSuccessfully(t *testing.T) {
	tags := TagMap{}
	err := AddTag(tags, "DB_URL", "database")
	require.NoError(t, err)
	assert.Equal(t, []string{"database"}, tags["DB_URL"])
}

func TestAddTag_RejectsDuplicate(t *testing.T) {
	tags := TagMap{"DB_URL": {"database"}}
	err := AddTag(tags, "DB_URL", "database")
	assert.EqualError(t, err, `tag "database" already exists on key "DB_URL"`)
}

func TestAddTag_NilMapReturnsError(t *testing.T) {
	err := AddTag(nil, "KEY", "tag")
	assert.Error(t, err)
}

func TestRemoveTag_RemovesExisting(t *testing.T) {
	tags := TagMap{"API_KEY": {"secret", "prod"}}
	RemoveTag(tags, "API_KEY", "secret")
	assert.Equal(t, []string{"prod"}, tags["API_KEY"])
}

func TestRemoveTag_RemovesKeyWhenEmpty(t *testing.T) {
	tags := TagMap{"API_KEY": {"secret"}}
	RemoveTag(tags, "API_KEY", "secret")
	_, exists := tags["API_KEY"]
	assert.False(t, exists)
}

func TestRemoveTag_NoopWhenMissing(t *testing.T) {
	tags := TagMap{"API_KEY": {"prod"}}
	RemoveTag(tags, "API_KEY", "nonexistent")
	assert.Equal(t, []string{"prod"}, tags["API_KEY"])
}

func TestRemoveTag_NilMapNoopPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		RemoveTag(nil, "KEY", "tag")
	})
}

func TestKeysWithTag_ReturnsMatchingKeys(t *testing.T) {
	tags := TagMap{
		"DB_URL":  {"database", "prod"},
		"DB_PASS": {"database", "secret"},
		"API_KEY": {"secret"},
	}
	keys := KeysWithTag(tags, "database")
	assert.ElementsMatch(t, []string{"DB_URL", "DB_PASS"}, keys)
}

func TestKeysWithTag_EmptyWhenNoMatch(t *testing.T) {
	tags := TagMap{"KEY": {"prod"}}
	keys := KeysWithTag(tags, "staging")
	assert.Empty(t, keys)
}
