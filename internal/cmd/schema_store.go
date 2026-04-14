package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// LoadSchema reads schema entries from a JSON file.
// Returns an empty slice if the file does not exist.
func LoadSchema(path string) ([]SchemaEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []SchemaEntry{}, nil
		}
		return nil, err
	}
	var entries []SchemaEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// SaveSchema writes schema entries to a JSON file, creating parent directories
// as needed.
func SaveSchema(path string, entries []SchemaEntry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// ValidateAgainstSchema checks that all required keys defined in the schema
// are present in the provided env map. It returns a list of missing keys.
func ValidateAgainstSchema(schema []SchemaEntry, env map[string]string) []string {
	var missing []string
	for _, entry := range schema {
		if !entry.Required {
			continue
		}
		val, ok := env[entry.Key]
		if !ok || val == "" {
			missing = append(missing, entry.Key)
		}
	}
	return missing
}
