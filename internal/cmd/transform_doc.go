package cmd

const transformDoc = `Transform applies a built-in string transformation to one or more
env variable values in the active (or specified) context.

Available sub-commands:

  upper   Convert each value to UPPERCASE
  lower   Convert each value to lowercase
  trim    Strip leading and trailing whitespace from each value

Examples:

  # Uppercase a single key
  envoy transform upper APP_ENV

  # Lowercase multiple keys
  envoy transform lower DB_HOST DB_USER

  # Trim whitespace from a key that may have been set with extra spaces
  envoy transform trim SECRET_KEY

  # Operate on a specific context
  envoy transform upper --context staging APP_ENV

All operations are performed in-place: the store is decrypted, the
selected values are mutated, and the store is re-encrypted before
being written back to disk.
`
