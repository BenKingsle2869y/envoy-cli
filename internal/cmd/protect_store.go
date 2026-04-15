package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ProtectedStore struct {
	Keys map[string]bool `json:"keys"`
}

func protectedFilePath(storePath string) string {
	dir := strings.TrimSuffix(storePath, filepath.Ext(storePath))
	return dir + ".protected.json"
}

func LoadProtected(storePath string) (*ProtectedStore, error) {
	path := protectedFilePath(storePath)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &ProtectedStore{Keys: make(map[string]bool)}, nil
	}
	if err != nil {
		return nil, err
	}
	var ps ProtectedStore
	if err := json.Unmarshal(data, &ps); err != nil {
		return nil, err
	}
	if ps.Keys == nil {
		ps.Keys = make(map[string]bool)
	}
	return &ps, nil
}

func SaveProtected(storePath string, ps *ProtectedStore) error {
	path := protectedFilePath(storePath)
	data, err := json.MarshalIndent(ps, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func MarkProtected(ps *ProtectedStore, key string) error {
	if ps.Keys == nil {
		return errors.New("protected store not initialised")
	}
	if ps.Keys[key] {
		return fmt.Errorf("key %q is already protected", key)
	}
	ps.Keys[key] = true
	return nil
}

func UnmarkProtected(ps *ProtectedStore, key string) error {
	if ps.Keys == nil || !ps.Keys[key] {
		return fmt.Errorf("key %q is not protected", key)
	}
	delete(ps.Keys, key)
	return nil
}

func IsProtected(ps *ProtectedStore, key string) bool {
	return ps != nil && ps.Keys[key]
}

func ProtectedKeys(ps *ProtectedStore) []string {
	if ps == nil {
		return nil
	}
	out := make([]string, 0, len(ps.Keys))
	for k := range ps.Keys {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
