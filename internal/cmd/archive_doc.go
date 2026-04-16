package cmd

// Archive command documentation
//
// Usage:
//   envoy archive create           Archive the active context
//   envoy archive list             List all archives for the active context
//   envoy archive restore <name>   Restore a named archive into a context
//
// Flags (create):
//   --passphrase   Passphrase to decrypt the store (or ENVOY_PASSPHRASE env var)
//
// Flags (restore):
//   --passphrase   Passphrase to encrypt the restored store
//   --context      Target context to restore into (defaults to active context)
//
// Examples:
//   envoy archive create --passphrase secret
//   envoy archive list
//   envoy archive restore production-1700000000 --context staging --passphrase secret
//
// Archives are stored under ~/.envoy/archives/<context>/ as JSON files.
// Each archive captures a full snapshot of all key-value entries at the
// time of creation and can be restored into any context.
