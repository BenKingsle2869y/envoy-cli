package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func noteFilePath(context string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", context+".notes.json")
}

func LoadNotes(context string) (map[string]string, error) {
	path := noteFilePath(context)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	var notes map[string]string
	if err := json.Unmarshal(data, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

func SaveNotes(context string, notes map[string]string) error {
	path := noteFilePath(context)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
