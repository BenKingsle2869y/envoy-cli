package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate a .env file against required keys in the active store",
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

func init() {
	validateCmd.Flags().StringSliceP("require", "r", nil, "Comma-separated list of required keys")
	validateCmd.Flags().BoolP("strict", "s", false, "Fail on keys present in file but missing from store")
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	requiredKeys, _ := cmd.Flags().GetStringSlice("require")
	strict, _ := cmd.Flags().GetBool("strict")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file %q: %w", filePath, err)
	}

	issue, err := validateEnvContent(string(data), requiredKeys, strict)
	if err != nil {
		return err
	}

	if len(issue) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✓ Validation passed")
		return nil
	}

	for _, i := range issue {
		fmt.Fprintln(cmd.OutOrStdout(), i)
	}
	return fmt.Errorf("validation failed with %d issue(s)", len(issue))
}

func validateEnvContent(content string, requiredKeys []string, strict bool) ([]string, error) {
	var issues []string

	parsed, err := parseEnvString(content)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	for _, key := range requiredKeys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, ok := parsed[key]; !ok {
			issues = append(issues, fmt.Sprintf("✗ missing required key: %s", key))
		}
	}

	if strict {
		for key := range parsed {
			found := false
			for _, rk := range requiredKeys {
				if strings.TrimSpace(rk) == key {
					found = true
					break
				}
			}
			if !found {
				issues = append(issues, fmt.Sprintf("⚠ unexpected key in strict mode: %s", key))
			}
		}
	}

	return issues, nil
}

func parseEnvString(content string) (map[string]string, error) {
	result := make(map[string]string)
	for i, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format %q", i+1, line)
		}
		result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return result, nil
}
