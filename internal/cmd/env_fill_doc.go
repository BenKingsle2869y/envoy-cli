package cmd

const envFillDoc = `Fill missing environment variables from the active envoy store into
the current process environment.

By default, only variables not already set in the environment are filled.
Use --overwrite to replace existing values.

Use --export to print 'export KEY=VALUE' statements to stdout instead of
applying them directly. This is useful for shell integration:

  eval $(envoy fill --export)

Examples:

  # Fill missing vars silently
  envoy fill --passphrase mysecret

  # Overwrite existing vars
  envoy fill --passphrase mysecret --overwrite

  # Print export statements for shell eval
  envoy fill --passphrase mysecret --export
`
