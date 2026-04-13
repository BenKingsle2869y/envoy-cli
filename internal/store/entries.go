package store

// Entries returns a shallow copy of all key-value pairs in the store.
// This is used by features like snapshots that need a read-only view
// of the current environment data.
func (s *Store) Entries() map[string]string {
	copy := make(map[string]string, len(s.data))
	for k, v := range s.data {
		copy[k] = v
	}
	return copy
}

// Keys returns a sorted list of all keys in the store.
func (s *Store) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
