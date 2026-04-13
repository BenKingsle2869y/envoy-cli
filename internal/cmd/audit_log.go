package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const auditLogFileName = "audit.log"

func auditLogPath(context string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", context+"_"+auditLogFileName)
}

// AppendAuditEntry writes a new entry to the audit log for the given context.
func AppendAuditEntry(context, action, key string) error {
	path := auditLogPath(context)

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("audit: create dir: %w", err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		Context:   context,
		Action:    action,
		Key:       key,
	}

	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}

// LoadAuditLog reads all entries from the audit log at path.
func LoadAuditLog(path string) ([]AuditEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entries []AuditEntry
	for _, raw := range splitLines(data) {
		if len(raw) == 0 {
			continue
		}
		var e AuditEntry
		if err := json.Unmarshal(raw, &e); err != nil {
			continue // skip malformed lines
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
