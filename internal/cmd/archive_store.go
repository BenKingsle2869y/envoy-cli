package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/envoy-cli/envoy/internal/store"
)

type Archive struct {
	Name      string            `json:"name"`
	Context   string            `json:"context"`
	CreatedAt time.Time         `json:"created_at"`
	Entries   []store.Entry     `json:"entries"`
}

func archiveDir(ctx string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", "archives", ctx)
}

func archiveFilePath(ctx, name string) string {
	return filepath.Join(archiveDir(ctx), name+".json")
}

func CreateArchive(ctx, name string, entries []store.Entry) error {
	dir := archiveDir(ctx)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create archive dir: %w", err)
	}
	a := Archive{
		Name:      name,
		Context:   ctx,
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
	}
	data, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("marshal archive: %w", err)
	}
	return os.WriteFile(archiveFilePath(ctx, name), data, 0600)
}

func LoadArchives(ctx string) ([]Archive, error) {
	dir := archiveDir(ctx)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var archives []Archive
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var a Archive
		if err := json.Unmarshal(data, &a); err == nil {
			archives = append(archives, a)
		}
	}
	return archives, nil
}

func RestoreArchive(ctx, name string) ([]store.Entry, error) {
	data, err := os.ReadFile(archiveFilePath(ctx, name))
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("archive %q not found", name)
	}
	if err != nil {
		return nil, err
	}
	var a Archive
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("parse archive: %w", err)
	}
	return a.Entries, nil
}
