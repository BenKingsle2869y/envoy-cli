package cmd

// gen command documentation
//
// Usage:
//   envoy-cli gen KEY [flags]
//
// Description:
//   The gen command generates a cryptographically random secret value
//   and stores it under the specified key in the active context store.
//
//   By default it produces a 32-character alphanumeric string. Use flags
//   to control length, character set, and whether to print the value.
//
// Flags:
//   -l, --length int    Length of generated value (default 32)
//   -s, --symbols       Include symbol characters
//   -d, --digits        Include digit characters (default true)
//   -p, --print         Print generated value to stdout
//
// Examples:
//   # Generate and store a 32-char alphanumeric secret
//   envoy-cli gen DB_PASSWORD
//
//   # Generate a 64-char value with symbols and print it
//   envoy-cli gen API_SECRET --length 64 --symbols --print
//
//   # Generate and store without printing
//   envoy-cli gen SESSION_KEY --length 48 --digits
