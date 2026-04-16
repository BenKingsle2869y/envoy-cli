package cmd

// Secret command documentation
//
// The `secret` command group allows users to mark specific environment
// variable keys as sensitive/secret. Marked keys are masked with asterisks
// when displayed in list or export output, preventing accidental exposure
// of credentials in terminal sessions or logs.
//
// Secret metadata is stored in a sidecar file (`secrets.json`) alongside
// the encrypted store file. The sidecar itself is not encrypted but contains
// only key names, not values.
//
// Usage:
//
//	envoy secret mark <key>    Mark a key as secret
//	envoy secret unmark <key>  Remove secret marking from a key
//	envoy secret list          List all keys marked as secret
//
// Example:
//
//	envoy secret mark DATABASE_PASSWORD
//	envoy list
//	# DATABASE_PASSWORD = ********
//
// Notes:
//   - Marking a key as secret does not affect encryption; all values
//     are encrypted regardless.
//   - Secret metadata persists across rotations and snapshots.
