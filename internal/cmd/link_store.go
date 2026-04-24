package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// links maps key -> target context name
type Links map[string]string

func linkFilePath(ctx string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", ctx+".links.json")
}

func LoadLinks(path string) (Links, error) {
	links := make(Links)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return links, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read links: %w", err)
	}
	if err := json.Unmarshal(data, &links); err != nil {
		return nil, fmt.Errorf("parse links: %w", err)
	}
	return links, nil
}

func SaveLinks(path string, links Links) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(links, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func AddLink(links Links, key, targetCtx string) error {
	if links == nil {
		return fmt.Errorf("links map is nil")
	}
	if _, exists := links[key]; exists {
		return fmt.Errorf("link for key %q already exists", key)
	}
	links[key] = targetCtx
	return nil
}

func RemoveLink(links Links, key string) error {
	if _, exists := links[key]; !exists {
		return fmt.Errorf("no link found for key %q", key)
	}
	delete(links, key)
	return nil
}

func ResolveLink(links Links, key string) (string, bool) {
	target, ok := links[key]
	return target, ok
}
