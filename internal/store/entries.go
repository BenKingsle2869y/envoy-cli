package store

import "sort"

// SortedKeys returns the store's keys in alphabetical order.
func (s *Store) SortedKeys() []string {
	keys := make([]string, 0, len(s.Entries))
	for k := range s.Entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Len returns the number of entries in the store.
func (s *Store) Len() int {
	return len(s.Entries)
}
