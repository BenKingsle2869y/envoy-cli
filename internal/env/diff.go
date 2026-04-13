package env

// DiffResult holds the comparison between two sets of env variables.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // [0] = old, [1] = new
	Unchanged map[string]string
}

// Diff compares two env maps (base vs target) and returns a DiffResult.
// base is typically the local env, target is the remote/stored env.
func Diff(base, target map[string]string) DiffResult {
	result := DiffResult{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Changed:   make(map[string][2]string),
		Unchanged: make(map[string]string),
	}

	// Find added and changed keys (present in target)
	for k, targetVal := range target {
		baseVal, exists := base[k]
		if !exists {
			result.Added[k] = targetVal
		} else if baseVal != targetVal {
			result.Changed[k] = [2]string{baseVal, targetVal}
		} else {
			result.Unchanged[k] = baseVal
		}
	}

	// Find removed keys (present in base but not in target)
	for k, baseVal := range base {
		if _, exists := target[k]; !exists {
			result.Removed[k] = baseVal
		}
	}

	return result
}

// HasChanges returns true if the diff contains any additions, removals, or changes.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Apply merges the target values into base according to the diff,
// returning a new map with the result.
func Apply(base, target map[string]string) map[string]string {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range target {
		result[k] = v
	}
	return result
}
