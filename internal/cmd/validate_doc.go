package cmd

// validate command documentation
//
// The validate command checks a .env file against a set of required keys
// and optionally enforces strict mode where only declared keys are allowed.
//
// Usage:
//
//	envoy validate <file> [flags]
//
// Flags:
//
//	-r, --require strings   Comma-separated list of required keys
//	-s, --strict            Fail if file contains keys not listed in --require
//
// Examples:
//
//	# Check that DB_HOST and DB_PORT are present
//	envoy validate .env --require DB_HOST,DB_PORT
//
//	# Strict mode: only the listed keys are permitted
//	envoy validate .env --require DB_HOST,DB_PORT --strict
//
// Exit codes:
//
//	0  All required keys are present (and no extra keys in strict mode)
//	1  One or more validation issues were found
