package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/envoy-cli/envoy/internal/crypto"
	"github.com/envoy-cli/envoy/internal/env"
)

const defaultStoreFile = ".envoy"

// EnvStore represents an encrypted store of environment variable sets.
type EnvStore struct {
	Environments map[string]string `json:"environments"` // name -> encrypted payload (base64)
}

// Load reads and decrypts an EnvStore from disk.
func Load(path, passphrase string) (*EnvStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &EnvStore{Environments: make(map[string]string)}, nil
		}
		return nil, fmt.Errorf("reading store: %w", err)
	}

	key, err := crypto.DeriveKey(passphrase)
	if err != nil {
		return nil, fmt.Errorf("deriving key: %w", err)
	}

	plaintext, err := crypto.Decrypt(key, string(data))
	if err != nil {
		return nil, fmt.Errorf("decrypting store: %w", err)
	}

	var store EnvStore
	if err := json.Unmarshal([]byte(plaintext), &store); err != nil {
		return nil, fmt.Errorf("parsing store: %w", err)
	}
	return &store, nil
}

// Save encrypts and writes the EnvStore to disk.
func Save(path, passphrase string, store *EnvStore) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating store directory: %w", err)
	}

	data, err := json.Marshal(store)
	if err != nil {
		return fmt.Errorf("serializing store: %w", err)
	}

	key, err := crypto.DeriveKey(passphrase)
	if err != nil {
		return fmt.Errorf("deriving key: %w", err)
	}

	ciphertext, err := crypto.Encrypt(key, string(data))
	if err != nil {
		return fmt.Errorf("encrypting store: %w", err)
	}

	return os.WriteFile(path, []byte(ciphertext), 0600)
}

// PutEnv stores a parsed env map under the given environment name.
func (s *EnvStore) PutEnv(name string, vars map[string]string) {
	s.Environments[name] = env.Serialize(vars)
}

// GetEnv retrieves and parses the env vars for the given environment name.
func (s *EnvStore) GetEnv(name string) (map[string]string, error) {
	raw, ok := s.Environments[name]
	if !ok {
		return nil, fmt.Errorf("environment %q not found in store", name)
	}
	return env.ParseString(raw)
}

// DefaultStorePath returns the default path for the store file.
func DefaultStorePath() string {
	return defaultStoreFile
}
