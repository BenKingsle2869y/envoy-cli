package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func aliasFilePath(context string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", context+".aliases.json")
}

// LoadAliases reads the alias map from disk. Returns an empty map if the file
// does not exist yet.
func LoadAliases(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load aliases: %w", err)
	}
	var aliases map[string]string
	if err := json.Unmarshal(data, &aliases); err != nil {
		return nil, fmt.Errorf("parse aliases: %w", err)
	}
	return aliases, nil
}

// SaveAliases writes the alias map to disk.
func SaveAliases(path string, aliases map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("save aliases: %w", err)
	}
	data, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal aliases: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

// AddAlias registers a new alias. Returns an error if the alias already exists.
func AddAlias(aliases map[string]string, alias, target string) error {
	if aliases == nil {
		return errors.New("alias map is nil")
	}
	if _, exists := aliases[alias]; exists {
		return fmt.Errorf("alias %q already exists", alias)
	}
	aliases[alias] = target
	return nil
}

// RemoveAlias deletes an alias. Returns an error if the alias does not exist.
func RemoveAlias(aliases map[string]string, alias string) error {
	if _, exists := aliases[alias]; !exists {
		return fmt.Errorf("alias %q not found", alias)
	}
	delete(aliases, alias)
	return nil
}

// ResolveAlias returns the target key for the given alias, or the original key
// if no alias is registered.
func ResolveAlias(aliases map[string]string, key string) string {
	if target, ok := aliases[key]; ok {
		return target
	}
	return key
}
