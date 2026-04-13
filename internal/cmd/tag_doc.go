package cmd

// Tag Command Documentation
//
// The `tag` command group allows users to annotate environment variable keys
// with arbitrary string labels. Tags are stored alongside the encrypted store
// and can be used to group, filter, or document variables by purpose or
// environment tier.
//
// Subcommands:
//
//   tag add <key> <tag>
//     Attaches the given tag to the specified key. Returns an error if the
//     tag already exists on that key.
//
//   tag remove <key> <tag>
//     Detaches the given tag from the specified key. No-op if the tag is
//     not present.
//
//   tag list <key>
//     Prints all tags currently associated with the specified key, one per
//     line. Prints a message if no tags are set.
//
// Examples:
//
//   envoy tag add DB_PASSWORD secret
//   envoy tag add DB_PASSWORD database
//   envoy tag list DB_PASSWORD
//   envoy tag remove DB_PASSWORD secret
