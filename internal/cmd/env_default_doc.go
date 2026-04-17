package cmd

// env default sets a key to a given value only if the key does not already
// exist in the active context store. This is useful for initialising
// environments with safe fallback values without overwriting intentional
// overrides.
//
// Usage:
//
//	envoy env default <key> --value <default-value> [--passphrase <pass>]
//
// Flags:
//
//	--value        The default value to assign when the key is absent (required)
//	--passphrase   Passphrase used to decrypt/encrypt the store
//
// Examples:
//
//	# Set DATABASE_URL only if it is not already present
//	envoy env default DATABASE_URL --value postgres://localhost/dev
//
//	# Use an explicit passphrase
//	envoy env default LOG_LEVEL --value info --passphrase mysecret
