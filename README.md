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

## Table of contents

- [Installation](#installation)
- [Quick start](#quick-start)
- [Commands](#commands)
  - [cluely auth](#cluely-auth)
  - [cluely sessions](#cluely-sessions)
  - [cluely sessions watch](#cluely-sessions-watch)
  - [cluely tags](#cluely-tags)
  - [cluely daemon](#cluely-daemon)
  - [cluely completion](#cluely-completion)
- [Exit codes](#exit-codes)
- [Updating](#updating)
- [License](#license)

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

# Watch for sessions to finish and run a script
cluely daemon start --exec "./on-complete.sh"
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
cluely sessions list                   # List recent sessions
cluely sessions list -n 5              # Show only 5 sessions
cluely sessions list --state finished  # Filter by state
cluely sessions list --since 24h       # Sessions from the last 24 hours
cluely sessions list --since 7d        # Sessions from the last 7 days
cluely sessions list --tag <tag-id>    # Filter by tag
cluely sessions get <session-id>       # View session details and transcript
```

Sessions and their tags are displayed with colored badges in the terminal.

**Updating and deleting sessions:**

```bash
cluely sessions update <session-id> --title "Q2 Planning"                    # Update title
cluely sessions update <session-id> --summary "Discussed roadmap priorities"  # Update summary
cluely sessions update <session-id> --title "Standup" --summary "Quick sync"  # Both at once
cluely sessions delete <session-id>                                           # Delete a session
```

**Tagging sessions:**

```bash
cluely sessions tag <session-id> <tag-id>      # Add a tag to a session
cluely sessions untag <session-id> <tag-id>    # Remove a tag from a session
```

**Field filtering** -- control which columns/sections are displayed:

```bash
# List: show only specific columns
cluely sessions list --fields id,title,tags

# List: hide specific columns
cluely sessions list --no-fields tags,state

# Get: show only summary and transcript
cluely sessions get <session-id> --fields summary,transcript

# Get: hide the transcript
cluely sessions get <session-id> --no-fields transcript
```

List columns: `id`, `state`, `title`, `tags`, `created`.
Get sections: `id`, `title`, `state`, `created`, `ended`, `tags`, `attendees`, `summary`, `transcript`.

**JSON output** -- add `--json` to get raw JSON, useful for scripting and piping:

```bash
cluely sessions list --json
cluely sessions list --json | jq '.items[].title'
cluely sessions get <session-id> --json
```

### `cluely sessions watch`

Watch for session starts and completions in real time. Runs in the foreground until Ctrl+C.

```bash
cluely sessions watch                                                    # Print all events
cluely sessions watch --exec "echo \$CLUELY_EVENT: \$CLUELY_SESSION_TITLE"  # Run on every event
cluely sessions watch --on end --exec "./on-complete.sh"                 # Only on session end
cluely sessions watch --on start --exec "notify-send 'Meeting started'"  # Only on session start
```

The `--exec` command has access to these environment variables:

| Variable | Description |
|----------|-------------|
| `CLUELY_EVENT` | Event type: `start` or `end` |
| `CLUELY_SESSION_ID` | Session ID |
| `CLUELY_SESSION_TITLE` | Session title (if available) |

Use `--on` to filter which events trigger the command (`start`, `end`, or both by default).

Example -- automatically export transcripts when sessions finish:

```bash
cluely sessions watch --on end --exec "cluely sessions get \$CLUELY_SESSION_ID --json > ~/transcripts/\$CLUELY_SESSION_ID.json"
```

### `cluely tags`

Create, list, and delete tags for organizing sessions. Aliased as `cluely t`. Tags are displayed as colored badges throughout the CLI.

```bash
cluely tags list                                    # List all tags
cluely tags create "Sales Call" --color "#4f46e5"   # Create a tag with a color
cluely tags create "Interview" --color "#059669"    # Colors are hex values
cluely tags delete <tag-id>                         # Delete a tag
```

Tags support `--json` output:

```bash
cluely tags list --json
```

### `cluely daemon`

Run the session watcher as a persistent background service. Uses launchd on macOS and systemd on Linux. The service auto-restarts on failure and runs on login.

```bash
# Start the service
cluely daemon start --exec "./on-complete.sh"

# Check if it's running
cluely daemon status

# View logs
cluely daemon logs
cluely daemon logs -f    # Follow mode (like tail -f)

# Stop and remove the service
cluely daemon stop
```

Logs are written to `~/.config/cluely/logs/watch.log`.

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
