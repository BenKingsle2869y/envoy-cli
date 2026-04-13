package cmd

// Snapshot commands provide point-in-time backup and restore functionality
// for environment stores managed by envoy-cli.
//
// Usage:
//
//	envoy snapshot create [label]   Create a new snapshot of the current context
//	envoy snapshot list             List all available snapshots
//	envoy snapshot restore <id>     Restore a snapshot by its ID
//
// Snapshots are stored as JSON files alongside the encrypted store file,
// in a directory named <store>.snapshots/. Each snapshot captures the full
// set of key-value pairs at the time of creation.
//
// Snapshot IDs are Unix nanosecond timestamps, ensuring uniqueness and
// natural chronological ordering.
//
// Examples:
//
//	# Create a labeled snapshot before a deploy
//	envoy snapshot create pre-deploy
//
//	# List all snapshots
//	envoy snapshot list
//
//	# Restore to a previous state
//	envoy snapshot restore 1718000000000000000
