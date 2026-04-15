package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"envoy-cli/internal/crypto"
)

// Store holds environment key/value entries.
type Store struct {
	Entries map[string]string `json:"entries"`
	Tags    map[string][]string `json:"tags,omitempty"`
	Pins    map[string]bool     `json:"pins,omitempty"`
}

// DefaultStorePath returns the default path for the active store.
func DefaultStorePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", "default.enc")
}

// Load decrypts and deserialises a store from disk.
func Load(path, passphrase string) (*Store, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Store{Entries: map[string]string{}}, nil
		}
		return nil, fmt.Errorf("read store: %w", err)
	}

	key, err := crypto.DeriveKey(passphrase, nil)
	if err != nil {
		return nil, fmt.Errorf("derive key: %w", err)
	}

	plain, err := crypto.Decrypt(key, data)
	if err != nil {
		return nil, fmt.Errorf("decrypt store: %w", err)
	}

	var s Store
	if err := json.Unmarshal(plain, &s); err != nil {
		return nil, fmt.Errorf("unmarshal store: %w", err)
	}
	if s.Entries == nil {
		s.Entries = map[string]string{}
	}
	return &s, nil
}

// Save serialises and encrypts a store to disk.
func Save(path, passphrase string, s *Store) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal store: %w", err)
	}

	key, err := crypto.DeriveKey(passphrase, nil)
	if err != nil {
		return fmt.Errorf("derive key: %w", err)
	}

	cipher, err := crypto.Encrypt(key, data)
	if err != nil {
		return fmt.Errorf("encrypt store: %w", err)
	}

	if err := os.WriteFile(path, cipher, 0o600); err != nil {
		return fmt.Errorf("write store: %w", err)
	}
	return nil
}

// Get returns the value for a key, or an error if not found.
func (s *Store) Get(key string) (string, error) {
	v, ok := s.Entries[key]
	if !ok {
		return "", fmt.Errorf("key %q not found", key)
	}
	return v, nil
}

// Set inserts or updates a key.
func (s *Store) Set(key, value string) {
	if s.Entries == nil {
		s.Entries = map[string]string{}
	}
	s.Entries[key] = value
}

// Delete removes a key from the store.
func (s *Store) Delete(key string) {
	delete(s.Entries, key)
}
