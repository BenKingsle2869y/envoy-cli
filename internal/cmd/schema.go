package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

// SchemaEntry describes a single env variable's schema definition.
type SchemaEntry struct {
	Key         string `json:"key"`
	Required    bool   `json:"required"`
	Description string `json:"description,omitempty"`
	Default     string `json:"default,omitempty"`
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Manage the env schema for the active context",
}

var schemaShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Print the current schema as JSON",
	RunE:  runSchemaShow,
}

var schemaAddCmd = &cobra.Command{
	Use:   "add <key>",
	Short: "Add or update a key in the schema",
	Args:  cobra.ExactArgs(1),
	RunE:  runSchemaAdd,
}

var (
	schemaRequired    bool
	schemaDescription string
	schemaDefault     string
)

func init() {
	schemaAddCmd.Flags().BoolVar(&schemaRequired, "required", false, "Mark key as required")
	schemaAddCmd.Flags().StringVar(&schemaDescription, "desc", "", "Description of the key")
	schemaAddCmd.Flags().StringVar(&schemaDefault, "default", "", "Default value for the key")
	schemaCmd.AddCommand(schemaShowCmd, schemaAddCmd)
	rootCmd.AddCommand(schemaCmd)
}

func runSchemaShow(cmd *cobra.Command, args []string) error {
	path := schemaFilePath()
	entries, err := LoadSchema(path)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "(no schema defined)")
		return nil
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func runSchemaAdd(cmd *cobra.Command, args []string) error {
	key := args[0]
	path := schemaFilePath()
	entries, err := LoadSchema(path)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}
	entries = upsertSchemaEntry(entries, SchemaEntry{
		Key:         key,
		Required:    schemaRequired,
		Description: schemaDescription,
		Default:     schemaDefault,
	})
	if err := SaveSchema(path, entries); err != nil {
		return fmt.Errorf("failed to save schema: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "schema entry %q saved\n", key)
	return nil
}

func schemaFilePath() string {
	ctx := ActiveContext()
	home, _ := os.UserHomeDir()
	return fmt.Sprintf("%s/.envoy/schema_%s.json", home, ctx)
}

func upsertSchemaEntry(entries []SchemaEntry, e SchemaEntry) []SchemaEntry {
	for i, existing := range entries {
		if existing.Key == e.Key {
			entries[i] = e
			return entries
		}
	}
	return append(entries, e)
}
