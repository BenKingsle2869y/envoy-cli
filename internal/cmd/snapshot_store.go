package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/your-org/envoy-cli/internal/store"
)

// Snapshot represents a point-in-time capture of an environment store.
type Snapshot struct {
	ID        string            `json:"id"`
	Label     string            `json:"label,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	Context   string            `json:"context"`
	Data      map[string]string `json:"data"`
}

func snapshotDir(storePath string) string {
	base := strings.TrimSuffix(storePath, filepath.Ext(storePath))
	return base + ".snapshots"
}

func snapshotFilePath(storePath, id string) string {
	return filepath.Join(snapshotDir(storePath), id+".json")
}

// CreateSnapshot captures the current store entries and writes a snapshot file.
func CreateSnapshot(storePath string, entries map[string]string, label, passphrase string) (*Snapshot, error) {
	if err := os.MkdirAll(snapshotDir(storePath), 0700); err != nil {
		return nil, err
	}
	snap := &Snapshot{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Context:   ActiveContext(),
		Data:      entries,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(snapshotFilePath(storePath, snap.ID), data, 0600); err != nil {
		return nil, err
	}
	return snap, nil
}

// LoadSnapshots reads all snapshots for a given store path, sorted by creation time.
func LoadSnapshots(storePath string) ([]Snapshot, error) {
	dir := snapshotDir(storePath)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var snaps []Snapshot
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		var s Snapshot
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, err
		}
		snaps = append(snaps, s)
	}
	return snaps, nil
}

// RestoreSnapshot loads a snapshot by ID and overwrites the current store.
func RestoreSnapshot(storePath, id, passphrase string) error {
	data, err := os.ReadFile(snapshotFilePath(storePath, id))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("snapshot %q not found", id)
		}
		return err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return err
	}
	st, err := store.Load(storePath, passphrase)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for k, v := range snap.Data {
		st.Set(k, v)
	}
	return store.Save(storePath, st, passphrase)
}
