package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint [file]",
	Short: "Lint a .env file for common issues",
	Long: `Lint parses a .env file and reports issues such as:
  - Lines missing an '=' separator
  - Keys with invalid characters
  - Empty keys
  - Duplicate keys`,
	Args: cobra.ExactArgs(1),
	RunE: runLint,
}

func init() {
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file %q: %w", filePath, err)
	}

	issues := lintEnvContent(string(data))
	if len(issues) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No issues found.")
		return nil
	}

	for _, issue := range issues {
		fmt.Fprintln(cmd.OutOrStdout(), issue)
	}
	return fmt.Errorf("%d issue(s) found", len(issues))
}

func lintEnvContent(content string) []string {
	var issues []string
	seen := make(map[string]int)

	lines := strings.Split(content, "\n")
	for i, raw := range lines {
		lineNum := i + 1
		line := strings.TrimSpace(raw)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		eqIdx := strings.Index(line, "=")
		if eqIdx < 0 {
			issues = append(issues, fmt.Sprintf("line %d: missing '=' separator: %q", lineNum, line))
			continue
		}

		key := strings.TrimSpace(line[:eqIdx])
		if key == "" {
			issues = append(issues, fmt.Sprintf("line %d: empty key", lineNum))
			continue
		}

		if !isValidEnvKey(key) {
			issues = append(issues, fmt.Sprintf("line %d: invalid key %q (only letters, digits, and underscores allowed)", lineNum, key))
		}

		if prev, ok := seen[key]; ok {
			issues = append(issues, fmt.Sprintf("line %d: duplicate key %q (first seen on line %d)", lineNum, key, prev))
		} else {
			seen[key] = lineNum
		}
	}

	return issues
}

func isValidEnvKey(key string) bool {
	for i, ch := range key {
		switch {
		case ch >= 'A' && ch <= 'Z':
		case ch >= 'a' && ch <= 'z':
		case ch >= '0' && ch <= '9' && i > 0:
		case ch == '_':
		default:
			return false
		}
	}
	return true
}
