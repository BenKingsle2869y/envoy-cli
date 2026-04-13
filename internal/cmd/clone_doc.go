package cmd

// Clone command documentation and examples.
//
// Usage:
//
//	envoy clone <source-context> <dest-context> [flags]
//
// Flags:
//
//	-o, --overwrite   Overwrite destination if it already exists
//
// Examples:
//
//	# Clone the "staging" context into a new "production" context
//	envoy clone staging production
//
//	# Clone and overwrite an existing destination context
//	envoy clone staging production --overwrite
//
// Notes:
//
//	The ENVOY_PASSPHRASE environment variable must be set. The same
//	passphrase is used for both reading the source and writing the
//	destination store, so both contexts share the same encryption key
//	after cloning. Use `envoy rotate` on the destination afterwards
//	if independent passphrases are required.
