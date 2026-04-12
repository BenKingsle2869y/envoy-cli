# envoy-cli

> A lightweight CLI for managing and syncing `.env` files across local and remote environments with encryption support.

---

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or download a pre-built binary from the [Releases](https://github.com/yourusername/envoy-cli/releases) page.

---

## Usage

```bash
# Initialize envoy in your project
envoy init

# Push your local .env to a remote environment
envoy push --env production

# Pull and decrypt .env from a remote environment
envoy pull --env staging

# Encrypt a .env file before sharing
envoy encrypt --file .env --out .env.enc

# Sync across multiple environments
envoy sync --from staging --to production
```

Run `envoy --help` to see all available commands and flags.

---

## Features

- 🔐 AES-256 encryption for `.env` files
- ☁️ Push/pull env configs to remote storage (S3, GCS, or custom backends)
- 🔄 Sync variables across multiple environments
- 🪶 Zero-dependency single binary

---

## Configuration

Envoy looks for a `envoy.yaml` file in your project root:

```yaml
remote: s3://my-bucket/envs
environments:
  - staging
  - production
```

---

## License

[MIT](LICENSE) © 2024 yourusername