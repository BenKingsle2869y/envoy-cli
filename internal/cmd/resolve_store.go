package cmd

import (
	"strings"

	"envoy-cli/internal/store"
)

// resolveReferences scans entries for values containing ${KEY} or $KEY
// references and substitutes them with values from the same store.
// Returns a map of keys whose values changed after resolution.
func resolveReferences(entries map[string]store.Entry) (map[string]string, map[string]string) {
	all := make(map[string]string, len(entries))
	for k, e := range entries {
		all[k] = e.Value
	}

	resolved := make(map[string]string)
	changed := make(map[string]string)

	for k, e := range entries {
		original := e.Value
		v := expandVars(original, all)
		resolved[k] = v
		if v != original {
			changed[k] = v
		}
	}

	return resolved, changed
}

// expandVars replaces ${VAR} and $VAR occurrences in s using the provided
// lookup map. Unknown references are left as-is.
func expandVars(s string, lookup map[string]string) string {
	result := s

	// Handle ${VAR} style
	for {
		start := strings.Index(result, "${") 
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "}")
		if end == -1 {
			break
		}
		end += start
		key := result[start+2 : end]
		if val, ok := lookup[key]; ok {
			result = result[:start] + val + result[end+1:]
		} else {
			break
		}
	}

	// Handle $VAR style (word boundary)
	for key, val := range lookup {
		token := "$" + key
		result = strings.ReplaceAll(result, token, val)
	}

	return result
}
