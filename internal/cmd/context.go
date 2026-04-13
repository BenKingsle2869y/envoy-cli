package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultContext = "development"

// ActiveContext reads the currently active environment context.
func ActiveContext() (string, error) {
	data, err := os.ReadFile(contextFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			return defaultContext, nil
		}
		return "", fmt.Errorf("read context file: %w", err)
	}
	ctx := strings.TrimSpace(string(data))
	if ctx == "" {
		return defaultContext, nil
	}
	return ctx, nil
}

// SetActiveContext writes the given context name as the active context.
func SetActiveContext(ctx string) error {
	if ctx == "" {
		return fmt.Errorf("context name must not be empty")
	}
	path := contextFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create context dir: %w", err)
	}
	return os.WriteFile(path, []byte(ctx+"\n"), 0600)
}

// ListContexts returns all context names discovered from store files.
func ListContexts() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".envoy")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	var ctxs []string
	for _, e := range entries {
		if e.IsDir() || e.Name() == "context" {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".enc")
		if name != e.Name() {
			ctxs = append(ctxs, name)
		}
	}
	return ctxs, nil
}

// StorePathForContext returns the store file path for the given context.
func StorePathForContext(ctx string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", ctx+".enc")
}
