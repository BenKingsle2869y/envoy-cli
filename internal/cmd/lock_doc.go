package cmd

// Lock command documentation.
//
// Usage:
//
//	envoy lock          Lock the active context
//	envoy unlock        Unlock the active context
//	envoy lock status   Show lock status
//
// Description:
//
//	The lock command prevents any modifications to the active context's
//	store. When a context is locked, commands that write to the store
//	(set, unset, import, rotate, merge, etc.) will refuse to proceed.
//
//	A lock file is stored at ~/.envoy/<context>.lock and contains the
//	UTC timestamp of when the lock was applied.
//
// Examples:
//
//	# Lock the current context
//	envoy lock
//
//	# Check lock status
//	envoy lock status
//
//	# Unlock the current context
//	envoy unlock
