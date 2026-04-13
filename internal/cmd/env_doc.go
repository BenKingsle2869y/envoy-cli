/*
Package cmd provides the `env` subcommand group for managing environment contexts.

# Overview

envoy supports multiple named environment contexts (e.g. development, staging,
production). Each context maps to its own encrypted store file located at:

	~/.envoy/<context>.enc

The active context is persisted to:

	~/.envoy/context

and defaults to "development" when the file is absent.

# Usage

	envoy env list          # List all available contexts, marking the active one
	envoy env use <name>    # Switch the active context

# Integration

Other commands such as `set`, `get`, `export`, `push`, and `pull` automatically
operate against the active context's store unless overridden by the --store flag.
*/
package cmd
