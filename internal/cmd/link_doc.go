package cmd

const linkDoc = `Link keys across contexts so a key in one context can reference
the value from another context.

Links are stored per-context in ~/.envoy/<context>.links.json.

Examples:

  # Link DATABASE_URL in current context to the "production" context
  envoy link add DATABASE_URL production

  # List all links in the current context
  envoy link list

  # Remove a link
  envoy link remove DATABASE_URL

Links are advisory metadata — they do not automatically sync values.
Use 'envoy copy' or 'envoy promote' to propagate values between contexts.
`
