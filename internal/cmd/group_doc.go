package cmd

// group command documentation
//
// The group command allows users to organize environment variable keys
// into named logical groups. This is useful for large .env stores where
// keys belong to different subsystems (e.g., "backend", "frontend", "infra").
//
// Groups are stored in a JSON sidecar file alongside the encrypted store:
//   ~/.envoy/<context>.groups.json
//
// Usage:
//
//	envoy group add <group> <key>     — Add a key to a group
//	envoy group remove <group> <key>  — Remove a key from a group
//	envoy group list <group>          — List all keys in a group
//
// Examples:
//
//	envoy group add backend DB_URL
//	envoy group add backend API_SECRET
//	envoy group list backend
//	envoy group remove backend DB_URL
//
// Notes:
//   - A key can belong to multiple groups.
//   - Groups are deleted automatically when all keys are removed.
//   - Group membership is not enforced against the store; keys may be
//     added to a group before or after being set in the store.
const groupDocString = "group: organize keys into named logical groups"
