package cmd

// Shell completion support for envoy-cli.
//
// The completion command generates shell-specific completion scripts that
// enable tab-completion for envoy-cli subcommands, flags, and arguments.
//
// Supported shells:
//
//   - bash
//   - zsh
//   - fish
//   - powershell
//
// Usage:
//
//	envoy completion bash
//	envoy completion zsh
//	envoy completion fish
//	envoy completion powershell
//
// The generated script should be sourced in the shell's configuration file
// or placed in the appropriate completion directory for persistent use.
//
// Installation examples:
//
// Bash (add to ~/.bashrc or ~/.bash_profile):
//
//	source <(envoy completion bash)
//
// Zsh (add to ~/.zshrc):
//
//	source <(envoy completion zsh)
//
// Fish (save to fish completions directory):
//
//	envoy completion fish > ~/.config/fish/completions/envoy.fish
//
// PowerShell (add to $PROFILE):
//
//	envoy completion powershell | Out-String | Invoke-Expression
