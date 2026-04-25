package cmd

const resolveDoc = `Resolve variable references within the active store.

Values that contain references in the form ${KEY} or $KEY are expanded
using other keys present in the same store. References to unknown keys
are left unchanged.

Examples:

  # Show resolved values without modifying the store
  envoy resolve --passphrase secret

  # Resolve references and write the expanded values back
  envoy resolve --passphrase secret --in-place

Notes:
  - Only single-level expansion is performed; circular references are
    not detected and will be left partially unresolved.
  - The --in-place flag overwrites the original (possibly templated)
    values, so use it with care.
`
