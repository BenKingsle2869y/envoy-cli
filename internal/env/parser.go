package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a collection of key-value environment variable pairs.
type EnvMap map[string]string

// ParseFile reads and parses a .env file into an EnvMap.
// It supports comments (lines starting with #) and blank lines.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open env file: %w", err)
	}
	defer f.Close()

	return parse(bufio.NewScanner(f))
}

// ParseString parses a raw .env string into an EnvMap.
func ParseString(content string) (EnvMap, error) {
	return parse(bufio.NewScanner(strings.NewReader(content)))
}

func parse(scanner *bufio.Scanner) (EnvMap, error) {
	envMap := make(EnvMap)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, err
		}

		envMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading env content: %w", err)
	}

	return envMap, nil
}

func parseLine(line string) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid line format (expected KEY=VALUE): %q", line)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	// Strip surrounding quotes if present
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
	}

	if key == "" {
		return "", "", fmt.Errorf("empty key in line: %q", line)
	}

	return key, value, nil
}

// Serialize converts an EnvMap back into a .env formatted string.
func Serialize(envMap EnvMap) string {
	var sb strings.Builder
	for k, v := range envMap {
		if strings.ContainsAny(v, " \t") {
			fmt.Fprintf(&sb, "%s=\"%s\"\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return sb.String()
}
