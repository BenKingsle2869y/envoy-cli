package env

import "strings"

// QuoteValue returns a quoted version of the value if it contains
// spaces, special characters, or is empty. Otherwise returns the raw value.
func QuoteValue(v string) string {
	if v == "" {
		return `""`
	}

	needsQuoting := strings.ContainsAny(v, " \t\n\r#$\"'\\`+"`"+'!')
	if !needsQuoting {
		return v
	}

	// Prefer double-quote wrapping; escape inner double quotes.
	escaped := strings.ReplaceAll(v, `"`, `\"`)
	return `"` + escaped + `"`
}
