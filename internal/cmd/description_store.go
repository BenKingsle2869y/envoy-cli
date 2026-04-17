package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func descriptionFilePath(context string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", context+".descriptions.json")
}

// LoadDescriptions reads key->description mapping from disk.
// Returns an empty map if the file does not exist.
func LoadDescriptions(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	var descs map[string]string
	if err := json.Unmarshal(data, &descs); err != nil {
		return nil, err
	}
	return descs, nil
}

// SaveDescriptions writes the key->description mapping to disk.
func SaveDescriptions(path string, descs map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(descs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
