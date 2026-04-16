package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

// CastResult holds a key, its raw value, and inferred type.
type CastResult struct {
	Key   string
	Value string
	Type  string
}

// CastEntries returns CastResult for each entry in the provided map.
func CastEntries(entries map[string]string) []CastResult {
	results := make([]CastResult, 0, len(entries))
	for k, v := range entries {
		results = append(results, CastResult{
			Key:   k,
			Value: v,
			Type:  inferType(v),
		})
	}
	return results
}

// CoerceToInt attempts to parse a string value as int64.
func CoerceToInt(value string) (int64, error) {
	v, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot cast %q to int: %w", value, err)
	}
	return v, nil
}

// CoerceToBool attempts to parse a string value as bool.
func CoerceToBool(value string) (bool, error) {
	v, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return false, fmt.Errorf("cannot cast %q to bool: %w", value, err)
	}
	return v, nil
}

// CoerceToFloat attempts to parse a string value as float64.
func CoerceToFloat(value string) (float64, error) {
	v, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0, fmt.Errorf("cannot cast %q to float: %w", value, err)
	}
	return v, nil
}
