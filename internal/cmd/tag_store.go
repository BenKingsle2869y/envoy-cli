package cmd

import "fmt"

// TagMap maps env key names to their associated tags.
type TagMap map[string][]string

// AddTag appends tag to the list for key, returning an error if it already exists.
func AddTag(tags TagMap, key, tag string) error {
	if tags == nil {
		return fmt.Errorf("tag map is nil")
	}
	for _, t := range tags[key] {
		if t == tag {
			return fmt.Errorf("tag %q already exists on key %q", tag, key)
		}
	}
	tags[key] = append(tags[key], tag)
	return nil
}

// RemoveTag removes tag from key's list. No-op if not present.
func RemoveTag(tags TagMap, key, tag string) {
	if tags == nil {
		return
	}
	current := tags[key]
	updated := current[:0]
	for _, t := range current {
		if t != tag {
			updated = append(updated, t)
		}
	}
	if len(updated) == 0 {
		delete(tags, key)
	} else {
		tags[key] = updated
	}
}

// KeysWithTag returns all keys that carry the given tag.
func KeysWithTag(tags TagMap, tag string) []string {
	var result []string
	for key, ts := range tags {
		for _, t := range ts {
			if t == tag {
				result = append(result, key)
				break
			}
		}
	}
	return result
}
