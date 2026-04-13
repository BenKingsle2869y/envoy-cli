// Package cmd provides the CLI commands for envoy-cli.
//
// # rotate
//
// The rotate command re-encrypts the local store file with a new passphrase
// without changing any of the stored environment variables. This is useful
// when you need to cycle credentials after a potential passphrase exposure.
//
// Usage:
//
//	envoy rotate [flags]
//
// Flags:
//
//	--store string          Path to the encrypted store file (default: ~/.envoy/store.enc)
//	--new-passphrase string New passphrase used to re-encrypt the store.
//	                        Can also be provided via the ENVOY_NEW_PASSPHRASE
//	                        environment variable.
//
// The current passphrase is resolved via the ENVOY_PASSPHRASE environment
// variable or from the keyfile written by `envoy init`.
//
// Example:
//
//	ENVOY_PASSPHRASE=old ENVOY_NEW_PASSPHRASE=new envoy rotate
//	envoy rotate --new-passphrase mynewsecret
package cmd
