// Package remote provides functionality for pushing and pulling
// encrypted .env stores to and from remote HTTP endpoints.
package remote

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 15 * time.Second

// Client is a lightweight HTTP client for syncing env stores.
type Client struct {
	BaseURL    string
	AuthToken  string
	HTTPClient *http.Client
}

// NewClient creates a new remote Client with the given base URL and auth token.
func NewClient(baseURL, authToken string) *Client {
	return &Client{
		BaseURL:   baseURL,
		AuthToken: authToken,
		HTTPClient: &http.Client{Timeout: defaultTimeout},
	}
}

// Push uploads the encrypted payload to the remote endpoint for the given environment.
func (c *Client) Push(env string, payload []byte) error {
	url := fmt.Sprintf("%s/envs/%s", c.BaseURL, env)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("remote: build push request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("remote: push request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("remote: push returned unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Pull downloads the encrypted payload from the remote endpoint for the given environment.
func (c *Client) Pull(env string) ([]byte, error) {
	url := fmt.Sprintf("%s/envs/%s", c.BaseURL, env)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("remote: build pull request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("remote: pull request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("remote: environment %q not found on server", env)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote: pull returned unexpected status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("remote: read response body: %w", err)
	}
	return data, nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/octet-stream")
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}
}
