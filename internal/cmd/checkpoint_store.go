package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yourusername/envoy-cli/internal/store"
)

type Checkpoint struct {
	Name      string                 `json:"name"`
	CreatedAt time.Time              `json:"created_at"`
	Entries   map[string]store.Entry `json:"entries"`
}

func checkpointDir(context string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", "checkpoints", context)
}

func checkpointFilePath(context, name string) string {
	return filepath.Join(checkpointDir(context), name+".json")
}

func CreateCheckpoint(context, name string, entries map[string]store.Entry) error {
	if name == "" {
		return errors.New("checkpoint name must not be empty")
	}
	dir := checkpointDir(context)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	cp := Checkpoint{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
	}
	data, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(checkpointFilePath(context, name), data, 0600)
}

func LoadCheckpoints(context string) ([]Checkpoint, error) {
	dir := checkpointDir(context)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var cps []Checkpoint
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		var cp Checkpoint
		if err := json.Unmarshal(data, &cp); err != nil {
			return nil, err
		}
		cps = append(cps, cp)
	}
	return cps, nil
}

func RestoreCheckpoint(context, name string) (map[string]store.Entry, error) {
	path := checkpointFilePath(context, name)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("checkpoint %q not found", name)
	}
	if err != nil {
		return nil, err
	}
	var cp Checkpoint
	if err := json.Unmarshal(data, &cp); err != nil {
		return nil, err
	}
	return cp.Entries, nil
}
