package cmd

// copy command copies all environment variables from one named context
// (environment) to another. By default it will not overwrite keys that
// already exist in the destination context; pass --overwrite to change
// this behaviour.
//
// Both contexts must share the same passphrase because the same
// encryption key is used to load the source store and save the
// destination store.
//
// Example usage:
//
//	# Copy all vars from "staging" into "production" (skip existing keys)
//	envoy copy staging production
//
//	# Copy and overwrite any conflicting keys
//	envoy copy staging production --overwrite
//
// The command reports how many variables were copied and how many were
// skipped so you always have a clear audit trail.
