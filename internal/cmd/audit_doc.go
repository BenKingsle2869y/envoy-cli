/*
Package cmd — audit command

# audit

The audit command displays a chronological log of changes made to the
active environment store. Each entry records the timestamp, context name,
action performed, and the key that was affected.

## Usage

    envoy audit [flags]

## Flags

    -n, --limit int   Maximum number of entries to show (default 20)

## Actions logged

  - set    — a key was added or updated
  - unset  — a key was removed
  - import — keys were imported from a file
  - copy   — a key was copied from another context
  - rotate — the store passphrase was rotated

## Example

    $ envoy audit -n 5
    TIMESTAMP              CONTEXT   ACTION  KEY
    2024-06-01T10:00:00Z   default   set     DATABASE_URL
    2024-06-01T10:01:00Z   default   set     SECRET_KEY
    2024-06-01T10:05:00Z   default   unset   OLD_TOKEN
*/
package cmd
