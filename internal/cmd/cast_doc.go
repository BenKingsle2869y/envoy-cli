package cmd

// Cast command documentation.
//
// Usage:
//
//	envoy cast [flags]
//
// Flags:
//
//	-p, --passphrase string   Passphrase to decrypt the store
//	-c, --context string      Context (environment) to use
//
// Description:
//
//	The cast command reads all key-value pairs from the active (or specified)
//	context store and displays each entry alongside its inferred type.
//
//	Supported types:
//	  - bool   (true, false, 1, 0)
//	  - int    (whole numbers)
//	  - float  (decimal numbers)
//	  - list   (values wrapped in square brackets)
//	  - string (everything else)
//
// Examples:
//
//	# Show typed values for the active context
//	envoy cast --passphrase mysecret
//
//	# Show typed values for a specific context
//	envoy cast --context production --passphrase mysecret
