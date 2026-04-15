package cmd

import (
	"fmt"
	"sort"
)

const pinnedTag = "__pinned__"

// AddPin marks a key as pinned by adding the internal pinned tag.
func AddPin(tags map[string][]string, key string) error {
	if tags == nil {
		return fmt.Errorf("tags map is nil")
	}
	for _, existing := range tags[pinnedTag] {
		if existing == key {
			return fmt.Errorf("key %q is already pinned", key)
		}
	}
	tags[pinnedTag] = append(tags[pinnedTag], key)
	return nil
}

// RemovePin removes the pinned marker from a key.
func RemovePin(tags map[string][]string, key string) error {
	if tags == nil {
		return fmt.Errorf("tags map is nil")
	}
	list := tags[pinnedTag]
	for i, k := range list {
		if k == key {
			tags[pinnedTag] = append(list[:i], list[i+1:]...)
			if len(tags[pinnedTag]) == 0 {
				delete(tags, pinnedTag)
			}
			return nil
		}
	}
	return fmt.Errorf("key %q is not pinned", key)
}

// IsPinned reports whether the given key is pinned.
func IsPinned(tags map[string][]string, key string) bool {
	if tags == nil {
		return false
	}
	for _, k := range tags[pinnedTag] {
		if k == key {
			return true
		}
	}
	return false
}

// PinnedKeys returns a sorted list of all pinned keys.
func PinnedKeys(tags map[string][]string) []string {
	if tags == nil {
		return nil
	}
	keys := make([]string, len(tags[pinnedTag]))
	copy(keys, tags[pinnedTag])
	sort.Strings(keys)
	return keys
}
