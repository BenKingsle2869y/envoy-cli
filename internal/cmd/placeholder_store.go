package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func placeholderFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", "placeholders.json")
}

// LoadPlaceholders reads the placeholder map from disk.
// Returns an empty map if the file does not exist.
func LoadPlaceholders(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	var ph map[string]string
	if err := json.Unmarshal(data, &ph); err != nil {
		return nil, err
	}
	return ph, nil
}

// SavePlaceholders writes the placeholder map to disk.
func SavePlaceholders(path string, ph map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(ph, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// RemovePlaceholder deletes a single key from the placeholder map.
func RemovePlaceholder(path, key string) error {
	ph, err := LoadPlaceholders(path)
	if err != nil {
		return err
	}
	if _, ok := ph[key]; !ok {
		return errors.New("placeholder not found: " + key)
	}
	delete(ph, key)
	return SavePlaceholders(path, ph)
}
