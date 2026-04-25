# dexctl

[![CI](https://github.com/dennismdejong/dexctl/actions/workflows/ci.yml/badge.svg)](https://github.com/dennismdejong/dexctl/actions/workflows/ci.yml)
[![Release](https://github.com/dennismdejong/dexctl/actions/workflows/release.yml/badge.svg)](https://github.com/dennismdejong/dexctl/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dennismdejong/dexctl)](https://github.com/dennismdejong/dexctl)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A command-line interface (CLI) tool for managing Dex clients. dexctl allows you to perform CRUD operations on Dex clients, list all clients, and connect to a Dex server with optional TLS certificate verification.

## Features

- Create, read, update, and delete Dex clients
- List all clients with formatted or JSON output
- Optional TLS certificate support (CA, client cert/key, or insecure)
- Cross-platform builds (macOS, Linux, Windows) via GoReleaser
- JSON output option for easy parsing

## Installation

### From Releases

Download the appropriate binary for your platform from the [GitHub Releases](https://github.com/dennismdejong/dexctl/releases) page.

### Building from Source

#### Using Go

```bash
git clone https://github.com/dennismdejong/dexctl.git
cd dexctl
go build -o dexctl .
```

#### Using GoReleaser (for cross-platform builds)

```bash
# Install GoReleaser: https://goreleaser.com/install/
# Then build snapshot releases
goreleaser release --snapshot --clean
# The binaries will be in the ./dist directory
```

## Usage

```bash
dexctl [command] [flags]
```

### Global Flags

| Flag | Description |
|------|-------------|
| `--server` | Dex gRPC server address (default: "localhost:5556") |
| `--cert` | TLS certificate file (for mutual TLS) |
| `--key` | TLS key file (for mutual TLS) |
| `--ca` | TLS CA certificate file (for server verification) |
| `--insecure` | Skip TLS certificate verification (use plaintext) |
| `--json` | Output as JSON |
| `-h, --help` | Help for dexctl |

### Commands

#### Client Management

Manage Dex clients with the `client` command.

```bash
dexctl client [subcommand] [flags]
```

##### Subcommands

| Subcommand | Description |
|------------|-------------|
| `create` | Create a new client |
| `get` | Get a client by ID |
| `update` | Update a client by ID |
| `delete` | Delete a client by ID |
| `list` | List all clients |

---

#### `dexctl client create`

Create a new Dex client.

**Flags:**

| Flag | Description | Required |
|------|-------------|----------|
| `--name` | Client name | Yes |
| `--public` | Set client as public (default: false) | No |
| `--redirect-uri` | Redirect URI (can be specified multiple times) | No |
| `--secret` | Client secret | No |
| `--trusted-peer` | Trusted peer (can be specified multiple times) | No |
| `--logo-url` | Logo URL | No |

**Example:**

```bash
# Create a confidential client
dexctl client create \
  --name "my-app" \
  --redirect-uri "https://myapp.com/callback" \
  --secret "supersecret" \
  --server "dex.example.com:5556" \
  --ca "/path/to/ca.crt"

# Create a public client
dexctl client create \
  --name "public-app" \
  --public \
  --redirect-uri "https://publicapp.com/auth" \
  --json
```

---

#### `dexctl client get`

Get a Dex client by ID.

**Flags:**

| Flag | Description | Required |
|------|-------------|----------|
| `--id` | Client ID | Yes |

**Example:**

```bash
# Get a client
dexctl client get --id "my-client-id" --server "dex.example.com:5556"

# Get a client with JSON output
dexctl client get --id "my-client-id" --json
```

---

#### `dexctl client update`

Update an existing Dex client by ID.

**Flags:**

| Flag | Description | Required |
|------|-------------|----------|
| `--id` | Client ID | Yes |
| `--name` | New client name | No |
| `--public` | Set public client (if provided) | No |
| `--redirect-uri` | New redirect URIs (can be specified multiple times) | No |
| `--trusted-peer` | New trusted peers (can be specified multiple times) | No |
| `--logo-url` | New logo URL | No |

**Note:** The Dex API does not allow updating the client secret via the UpdateClient endpoint. To change the secret, you must delete and recreate the client.

**Example:**

```bash
# Update client name and redirect URI
dexctl client update \
  --id "my-client-id" \
  --name "updated-app" \
  --redirect-uri "https://updatedapp.com/callback" \
  --server "dex.example.com:5556"

# Update with JSON output
dexctl client update \
  --id "my-client-id" \
  --name "updated-app" \
  --json
```

---

#### `dexctl client delete`

Delete a Dex client by ID.

**Flags:**

| Flag | Description | Required |
|------|-------------|----------|
| `--id` | Client ID | Yes |

**Example:**

```bash
# Delete a client
dexctl client delete --id "my-client-id" --server "dex.example.com:5556"

# Force deletion without confirmation (no prompt by default)
dexctl client delete --id "my-client-id" --insecure
```

---

#### `dexctl client list`

List all Dex clients.

**Flags:** None specific (uses global flags)

**Example:**

```bash
# List all clients (default format)
dexctl client list --server "dex.example.com:5556"

# List all clients as JSON
dexctl client list --json

# List all clients with TLS verification
dexctl client list --ca "/path/to/ca.crt" --server "dex.example.com:5556"
```

**Output Format (default):**

```
ID: client-id-1
Name: my-client
Public: false
RedirectURIs: [https://myapp.com/callback]
TrustedPeers: [trusted-peer-1, trusted-peer-2]
LogoURL: https://example.com/logo.png

---
ID: client-id-2
Name: another-client
Public: true
RedirectURIs: [https://another.com/auth]
TrustedPeers: []
LogoURL: 
```

## TLS Connection Options

dexctl supports multiple ways to secure the connection to the Dex server:

1. **Insecure (Plaintext)**: Use `--insecure` for plaintext gRPC (not recommended for production)
2. **Server Verification Only**: Provide `--ca` to verify the server's certificate
3. **Mutual TLS**: Provide `--cert`, `--key`, and `--ca` for client certificate authentication

**Examples:**

```bash
# Insecure (plaintext)
dexctl client list --insecure --server localhost:5556

# Server verification with CA
dexctl client list --ca /path/to/ca.crt --server dex.example.com:5556

# Mutual TLS
dexctl client list \
  --cert /path/to/client.crt \
  --key /path/to/client.key \
  --ca /path/to/ca.crt \
  --server dex.example.com:5556
```

## GitHub Actions CI/CD

This project uses GitHub Actions for continuous integration and releases. The workflows are defined in the `.github/workflows/` directory.

### Workflows

#### CI Workflow (`ci.yml`)

Runs on every push to the main, master, or develop branches, and on pull requests.

**Jobs:**
- **Test**: Runs tests on multiple Go versions (1.21, 1.22, 1.23, 1.24) and operating systems (Ubuntu, macOS, Windows)
- **Lint**: Runs golangci-lint for code quality checks
- **Security**: Runs Gosec security scanner
- **Build**: Builds binaries for all platforms
- **Build Multi-arch**: Builds multi-platform binaries using GoReleaser
- **Build Docker**: Builds Docker images (on push only)

**Environment Variables:**
- `CODECOV_TOKEN`: For uploading test coverage to Codecov (optional)

#### Release Workflow (`release.yml`)

Runs when a version tag is pushed (e.g., `v1.0.0`).

**Jobs:**
- **Release**: Creates GitHub releases with binaries for all platforms using GoReleaser
- **Release Docker**: Pushes Docker images to Docker Hub
- **Release AUR**: Updates Arch Linux AUR package (for major releases)
- **Verify**: Verifies the release artifacts

**Requirements:**
- `GITHUB_TOKEN`: Automatically provided by GitHub Actions
- `GORELEASER_KEY`: Your GoReleaser Pro key (for Pro features, optional)
- `DOCKERHUB_USERNAME`: Docker Hub username (for Docker releases)
- `DOCKERHUB_TOKEN`: Docker Hub access token (for Docker releases)

**Usage:**
```bash
# Create a tag and push to trigger a release
git tag v1.0.0
git push origin v1.0.0
```

#### Snapshot Workflow (`snapshot.yml`)

Allows creating snapshot releases manually via GitHub Actions UI.

**Jobs:**
- **Snapshot**: Creates a snapshot release for testing
- **Build Standalone**: Builds standalone binaries for testing

**Usage:**
1. Go to the repository's Actions tab in GitHub
2. Select "Snapshot Release" workflow
3. Click "Run workflow"
4. Download artifacts from the workflow run

### Setting Up Secrets

To set up the required secrets for GitHub Actions:

1. Go to your repository Settings > Secrets and variables > Actions
2. Add the following secrets:

| Secret | Description | Required |
|--------|-------------|----------|
| `CODECOV_TOKEN` | Codecov upload token | No |
| `GORELEASER_KEY` | GoReleaser Pro license key | No |
| `DOCKERHUB_USERNAME` | Docker Hub username | No |
| `DOCKERHUB_TOKEN` | Docker Hub access token | No |

## Development

### Prerequisites

- Go 1.20+
- Git
- (Optional) GoReleaser for building cross-platform releases

### Running Tests

```bash
go test ./...
```

### Building for Development

```bash
go build -o dexctl .
```

### Building Releases

```bash
# Install GoReleaser first
# See: https://goreleaser.com/install/

# Create a snapshot release (for testing)
goreleaser release --snapshot --clean

# Create a real release (requires GitHub token with repo scope)
goreleaser release --rm-dist
```

## License

MIT

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing-feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

*Note: Remember to replace `dennismdejong` with your actual GitHub username in the URLs and badges above.*