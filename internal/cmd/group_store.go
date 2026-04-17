package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func groupFilePath(context string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", context+".groups.json")
}

func LoadGroups(path string) (map[string][]string, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string][]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	var groups map[string][]string
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func SaveGroups(path string, groups map[string][]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func AddToGroup(groups map[string][]string, group, key string) error {
	for _, k := range groups[group] {
		if k == key {
			return fmt.Errorf("key %q already in group %q", key, group)
		}
	}
	groups[group] = append(groups[group], key)
	return nil
}

func RemoveFromGroup(groups map[string][]string, group, key string) error {
	keys, ok := groups[group]
	if !ok {
		return fmt.Errorf("group %q not found", group)
	}
	for i, k := range keys {
		if k == key {
			groups[group] = append(keys[:i], keys[i+1:]...)
			if len(groups[group]) == 0 {
				delete(groups, group)
			}
			return nil
		}
	}
	return fmt.Errorf("key %q not found in group %q", key, group)
}

func KeysInGroup(groups map[string][]string, group string) []string {
	return groups[group]
}
