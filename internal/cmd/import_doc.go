package cmd

// Import command documentation.
//
// The import command reads a standard .env file from disk and loads its
// key-value pairs into the encrypted envoy store.
//
// Usage:
//
//	envoy import <file> [--overwrite]
//
// Arguments:
//
//	file  Path to the .env file to import.
//
// Flags:
//
//	--overwrite  When set, existing keys in the store are overwritten by
//	             values from the imported file. Without this flag, existing
//	             keys are preserved and only new keys are added.
//
// Examples:
//
//	# Import variables, keeping any existing values
//	envoy import .env.production
//
//	# Import and overwrite all matching keys
//	envoy import --overwrite .env.staging
//
// The passphrase is resolved via ENVOY_PASSPHRASE or a key file, consistent
// with all other envoy commands.
