// Package remote implements a thin HTTP client used by envoy-cli to push and
// pull encrypted environment stores to and from a compatible remote server.
//
// # Overview
//
// The [Client] type wraps a standard net/http.Client and exposes two
// operations:
//
//   - Push – serialises and uploads an encrypted payload for a named
//     environment via HTTP PUT.
//   - Pull – downloads an encrypted payload for a named environment via
//     HTTP GET.
//
// Authentication is handled through a Bearer token set on every request.
// The payload is treated as raw bytes (application/octet-stream) so that
// encryption and serialisation remain the responsibility of the caller —
// typically [internal/store].
//
// # Usage
//
//	client := remote.NewClient("https://envoy.example.com", os.Getenv("ENVOY_TOKEN"))
//
//	// Upload
//	if err := client.Push("production", encryptedBytes); err != nil { ... }
//
//	// Download
//	data, err := client.Pull("production")
package remote
