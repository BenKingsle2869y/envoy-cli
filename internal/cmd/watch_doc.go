package cmd

const watchDoc = `Watch the active environment store for changes and print a diff
whenever keys are added, modified, or removed.

The command polls the store file at a configurable interval and
prints a summary of changes to stdout. Press Ctrl+C to stop.

Examples:

  # Watch with default 5-second interval
  envoy watch

  # Watch with a custom 10-second interval
  envoy watch --interval 10

Change indicators:
  +  key was added
  ~  key value was changed
  -  key was removed

The passphrase is resolved from the ENVOY_PASSPHRASE environment
variable or the configured keyfile.
`
