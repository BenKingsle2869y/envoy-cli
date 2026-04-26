package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func maskFilePath(context string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", context+".masks.json")
}

// LoadMasks reads the masked-keys set from disk.
// Returns an empty map if the file does not exist.
func LoadMasks(path string) (map[string]bool, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]bool{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read masks: %w", err)
	}
	var m map[string]bool
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse masks: %w", err)
	}
	return m, nil
}

// SaveMasks writes the masked-keys set to disk.
func SaveMasks(path string, masks map[string]bool) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(masks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// MarkMasked adds key to the masks set.
func MarkMasked(masks map[string]bool, key string) error {
	if masks == nil {
		return fmt.Errorf("masks map is nil")
	}
	if masks[key] {
		return fmt.Errorf("key %q is already masked", key)
	}
	masks[key] = true
	return nil
}

// UnmarkMasked removes key from the masks set.
func UnmarkMasked(masks map[string]bool, key string) error {
	if !masks[key] {
		return fmt.Errorf("key %q is not masked", key)
	}
	delete(masks, key)
	return nil
}

// IsMasked reports whether key is in the masks set.
func IsMasked(masks map[string]bool, key string) bool {
	return masks[key]
}

// MaskedKeys returns a sorted slice of all masked keys.
func MaskedKeys(masks map[string]bool) []string {
	keys := make([]string, 0, len(masks))
	for k := range masks {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
