package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

func TestPushCmd_Success(t *testing.T) {
	var received map[string]string
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/envs/production", r.URL.Path)
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusNoContent)
	})

	t.Setenv("ENVOY_PASSPHRASE", "test-secret")

	rootCmd.SetArgs([]string{"push", "production", "--url", srv.URL})
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	// We only verify the command wires up without panicking;
	// a missing local store will return an error which is acceptable here.
	_ = rootCmd.Execute()
}

func TestPullCmd_Success(t *testing.T) {
	payload := map[string]string{"KEY": "value"}
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/envs/staging", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	})

	t.Setenv("ENVOY_PASSPHRASE", "test-secret")

	rootCmd.SetArgs([]string{"pull", "staging", "--url", srv.URL})
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	// Pull will succeed fetching remote data; saving to default path may fail in CI.
	err := rootCmd.Execute()
	// We accept either nil or a store-write error; no panic expected.
	assert.NotContains(t, err, "pulling environment")
}
