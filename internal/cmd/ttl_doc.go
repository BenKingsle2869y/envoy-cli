package cmd

// TTL (time-to-live) support lets you mark individual env keys with an expiry
// timestamp. Once a key's TTL has passed it is considered "expired" and tools
// such as `envoy watch` or CI pipelines can act on that signal.
//
// Usage:
//
//	envoy ttl set   <key> <duration>   Set an expiry on a key
//	envoy ttl show  <key>              Print the expiry time (and remaining time)
//	envoy ttl clear <key>              Remove the expiry from a key
//
// Duration format:
//
//	Standard Go durations are accepted (e.g. 1h30m, 45m, 24h).
//	As a convenience the "d" suffix is also supported for whole days (e.g. 7d).
//
// Storage:
//
//	TTL metadata is stored in a separate JSON file alongside the encrypted store
//	(~/.envoy/<context>.ttl.json).  The file is NOT encrypted because it contains
//	only key names and timestamps — no secret values.
//
// Example:
//
//	$ envoy ttl set DATABASE_URL 7d
//	TTL set: DATABASE_URL expires at 2025-08-01T12:00:00Z
//
//	$ envoy ttl show DATABASE_URL
//	DATABASE_URL expires at 2025-08-01T12:00:00Z (167h59m remaining)
//
//	$ envoy ttl clear DATABASE_URL
//	TTL cleared for "DATABASE_URL"
