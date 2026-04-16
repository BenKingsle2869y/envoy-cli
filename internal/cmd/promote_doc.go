package cmd

const promoteDoc = `Promote copies environment variables from a source context into a
destination context.

By default, keys that already exist in the destination are skipped.
Use --overwrite to replace them.

Both contexts must be accessible with the same passphrase. If you need
different passphrases per context, export each store first and import
manually.

Examples:
  # Promote all vars from staging to production (skip existing)
  envoy promote staging production

  # Promote and overwrite any conflicting keys
  envoy promote staging production --overwrite

  # Provide passphrase inline
  envoy promote staging production --passphrase mysecret
`
