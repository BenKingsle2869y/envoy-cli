package crypto

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	envKeyVar    = "ENVOY_PASSPHRASE"
	keyFileName  = ".envoy_key"
)

// ErrNoPassphrase is returned when no passphrase can be resolved.
var ErrNoPassphrase = errors.New("no passphrase found: set ENVOY_PASSPHRASE env var or create a .envoy_key file")

// ResolvePassphrase attempts to load the passphrase from:
// 1. ENVOY_PASSPHRASE environment variable
// 2. .envoy_key file in the current or home directory
func ResolvePassphrase() (string, error) {
	if val := os.Getenv(envKeyVar); val != "" {
		return strings.TrimSpace(val), nil
	}

	// Try current directory first
	if pass, err := readKeyFile("."); err == nil {
		return pass, nil
	}

	// Fallback to home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", ErrNoPassphrase
	}

	if pass, err := readKeyFile(home); err == nil {
		return pass, nil
	}

	return "", ErrNoPassphrase
}

// SavePassphraseToFile writes the passphrase to a .envoy_key file in the given directory.
func SavePassphraseToFile(dir, passphrase string) error {
	path := filepath.Join(dir, keyFileName)
	return os.WriteFile(path, []byte(passphrase+"\n"), 0600)
}

func readKeyFile(dir string) (string, error) {
	path := filepath.Join(dir, keyFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
