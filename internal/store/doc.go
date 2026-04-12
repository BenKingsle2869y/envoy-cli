// Package store provides functionality for persisting and retrieving
// named sets of environment variables in an encrypted store file.
//
// The store serializes environment maps to .env format, bundles them
// into a JSON envelope, and encrypts the result using AES-GCM via the
// crypto package. The encrypted blob is written to a single file on
// disk (default: .envoy) with 0600 permissions.
//
// Typical usage:
//
//	s, err := store.Load(".envoy", passphrase)
//	if err != nil { ... }
//
//	vars, err := s.GetEnv("production")
//	if err != nil { ... }
//
//	s.PutEnv("staging", map[string]string{"APP_ENV": "staging"})
//	if err := store.Save(".envoy", passphrase, s); err != nil { ... }
package store
