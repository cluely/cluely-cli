<div align="center">
<br>
<img src="assets/banner.png" alt="Cluely" width="600">
<br><br>
<p><strong>Command-line interface for the <a href="https://cluely.com">Cluely</a> platform.</strong></p>
<p>Manage your meeting sessions, transcripts, and insights from the terminal.</p>
<a href="https://github.com/cluely/cluely-cli/releases/latest"><img src="https://img.shields.io/github/v/release/cluely/cluely-cli" alt="Release"></a>
<a href="https://github.com/cluely/cluely-cli/actions"><img src="https://github.com/cluely/cluely-cli/actions/workflows/release.yml/badge.svg" alt="Build"></a>
<a href="https://github.com/cluely/cluely-cli"><img src="https://img.shields.io/github/license/cluely/cluely-cli" alt="License"></a>
<br><br>
</div>

## Installation

### Homebrew (macOS and Linux)

```bash
brew tap cluely/tap
brew install cluely
```

### Shell script (macOS and Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/cluely/cluely-cli/main/install.sh | sh
```

### Manual download

Download the latest binary for your platform from the [Releases](https://github.com/cluely/cluely-cli/releases/latest) page.

| Platform       | Download |
|----------------|----------|
| macOS (Apple Silicon) | `cluely_*_darwin_arm64.tar.gz` |
| macOS (Intel)  | `cluely_*_darwin_amd64.tar.gz` |
| Linux (x86_64) | `cluely_*_linux_amd64.tar.gz` |
| Linux (ARM64)  | `cluely_*_linux_arm64.tar.gz` |
| Windows (x86_64) | `cluely_*_windows_amd64.zip` |
| Windows (ARM64)  | `cluely_*_windows_arm64.zip` |

### Build from source

Requires [Go](https://go.dev/dl/) 1.26+.

```bash
git clone https://github.com/cluely/cluely-cli.git
cd cluely-cli
make install
```

## Quick start

```bash
# Authenticate with your Cluely account
cluely auth login

# List your recent sessions
cluely sessions list

# View a specific session with transcript
cluely sessions get <session-id>
```

## Commands

### `cluely auth`

Manage authentication. Credentials are stored securely in your OS keyring (macOS Keychain, Windows Credential Manager, or Linux Secret Service).

```
cluely auth login     # Open browser to sign in
cluely auth logout    # Clear stored credentials
cluely auth status    # Check if you're logged in
```

### `cluely sessions`

List and view meeting sessions. Aliased as `cluely s`.

```
cluely sessions list                  # List recent sessions
cluely sessions list -n 5             # Show only 5 sessions
cluely sessions list --state finished # Filter by state
cluely sessions list --since 24h      # Sessions from the last 24 hours
cluely sessions list --since 7d       # Sessions from the last 7 days
cluely sessions get <session-id>      # View session details and transcript
```

**JSON output** -- add `--json` to get raw JSON, useful for scripting and piping:

```bash
cluely sessions list --json
cluely sessions list --json | jq '.items[].title'
cluely sessions get <session-id> --json
```

### `cluely sessions watch`

Watch for sessions to finish in real time. Runs continuously until Ctrl+C.

```bash
cluely sessions watch                          # Print when sessions finish
cluely sessions watch --exec "say done"        # Run a command on completion
cluely sessions watch --exec "./process.sh"    # Run a script
```

The `--exec` command has access to these environment variables:

| Variable | Description |
|----------|-------------|
| `CLUELY_SESSION_ID` | Session ID |
| `CLUELY_SESSION_TITLE` | Session title (if available) |

Example -- automatically export transcripts when sessions finish:

```bash
cluely sessions watch --exec "cluely sessions get \$CLUELY_SESSION_ID --json > ~/transcripts/\$CLUELY_SESSION_ID.json"
```

### `cluely completion`

Generate shell completions:

```bash
# Bash
cluely completion bash > /etc/bash_completion.d/cluely

# Zsh
cluely completion zsh > "${fpath[1]}/_cluely"

# Fish
cluely completion fish > ~/.config/fish/completions/cluely.fish
```

## Exit codes

| Code | Meaning |
|------|---------|
| `0`  | Success |
| `1`  | Runtime error (auth failure, network error, etc.) |
| `2`  | Usage error (invalid flags or arguments) |

`cluely auth status` returns exit code `1` when not authenticated, which is useful for scripting:

```bash
cluely auth status && cluely sessions list
```

## Updating

### Homebrew

```bash
brew upgrade cluely
```

### Shell script

Re-run the install script -- it always fetches the latest version:

```bash
curl -fsSL https://raw.githubusercontent.com/cluely/cluely-cli/main/install.sh | sh
```

## License

Proprietary. Copyright Cluely, Inc.
