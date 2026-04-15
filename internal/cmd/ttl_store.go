package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// TTLMap maps key names to their expiry times.
type TTLMap map[string]time.Time

func ttlFilePath(ctx string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", ctx+".ttl.json")
}

// LoadTTLs reads the TTL file for the given context.
func LoadTTLs(path string) (TTLMap, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return TTLMap{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read ttl file: %w", err)
	}

	raw := map[string]string{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse ttl file: %w", err)
	}

	ttls := make(TTLMap, len(raw))
	for k, v := range raw {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return nil, fmt.Errorf("invalid time for key %q: %w", k, err)
		}
		ttls[k] = t
	}
	return ttls, nil
}

func saveTTLs(path string, ttls TTLMap) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	raw := make(map[string]string, len(ttls))
	for k, v := range ttls {
		raw[k] = v.UTC().Format(time.RFC3339)
	}
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// SetTTL persists an expiry time for the given key.
func SetTTL(path, key string, expiry time.Time) error {
	ttls, err := LoadTTLs(path)
	if err != nil {
		return err
	}
	ttls[key] = expiry.UTC()
	return saveTTLs(path, ttls)
}

// ClearTTL removes the expiry for the given key.
func ClearTTL(path, key string) error {
	ttls, err := LoadTTLs(path)
	if err != nil {
		return err
	}
	delete(ttls, key)
	return saveTTLs(path, ttls)
}

// ExpiredKeys returns all keys whose TTL has passed.
func ExpiredKeys(ttls TTLMap) []string {
	now := time.Now().UTC()
	var expired []string
	for k, exp := range ttls {
		if now.After(exp) {
			expired = append(expired, k)
		}
	}
	return expired
}

// parseTTLDuration parses durations like "24h", "7d", "30m".
func parseTTLDuration(s string) (time.Duration, error) {
	if len(s) > 1 && s[len(s)-1] == 'd' {
		days := s[:len(s)-1]
		var n int
		if _, err := fmt.Sscanf(days, "%d", &n); err != nil {
			return 0, fmt.Errorf("invalid day count: %s", days)
		}
		return time.Duration(n) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}
