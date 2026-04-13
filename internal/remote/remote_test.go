package remote_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/remote"
)

func TestPush_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/envs/production" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := remote.NewClient(server.URL, "test-token")
	err := client.Push("production", []byte("encrypted-payload"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestPush_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := remote.NewClient(server.URL, "")
	err := client.Push("staging", []byte("data"))
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestPull_Success(t *testing.T) {
	expected := []byte("secret-encrypted-data")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(expected)
	}))
	defer server.Close()

	client := remote.NewClient(server.URL, "token-abc")
	data, err := client.Pull("production")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if string(data) != string(expected) {
		t.Errorf("expected %q, got %q", expected, data)
	}
}

func TestPull_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := remote.NewClient(server.URL, "")
	_, err := client.Pull("unknown-env")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestNewClient_DefaultsHTTPClient(t *testing.T) {
	client := remote.NewClient("https://example.com", "tok")
	if client.HTTPClient == nil {
		t.Error("expected HTTPClient to be initialised")
	}
	if client.BaseURL != "https://example.com" {
		t.Errorf("unexpected BaseURL: %s", client.BaseURL)
	}
}
