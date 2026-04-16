package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type SecretStore struct {
	Keys map[string]bool `json:"keys"`
}

func secretFilePath(storePath string) string {
	dir := filepath.Dir(storePath)
	return filepath.Join(dir, "secrets.json")
}

func LoadSecrets(path string) (*SecretStore, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &SecretStore{Keys: map[string]bool{}}, nil
	}
	if err != nil {
		return nil, err
	}
	var s SecretStore
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if s.Keys == nil {
		s.Keys = map[string]bool{}
	}
	return &s, nil
}

func SaveSecrets(path string, s *SecretStore) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func MarkSecret(s *SecretStore, key string) error {
	if s.Keys == nil {
		return errors.New("secret store is nil")
	}
	if s.Keys[key] {
		return errors.New("key already marked as secret")
	}
	s.Keys[key] = true
	return nil
}

func UnmarkSecret(s *SecretStore, key string) error {
	if !s.Keys[key] {
		return errors.New("key is not marked as secret")
	}
	delete(s.Keys, key)
	return nil
}

func IsSecret(s *SecretStore, key string) bool {
	return s != nil && s.Keys[key]
}

func SecretKeys(s *SecretStore) []string {
	keys := make([]string, 0, len(s.Keys))
	for k := range s.Keys {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func MaskIfSecret(s *SecretStore, key, value string) string {
	if IsSecret(s, key) {
		return strings.Repeat("*", 8)
	}
	return value
}
