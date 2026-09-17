# Project Context

## Purpose

Soltty is a command-line interface for [Solidtime](https://github.com/solidtime-io/solidtime)
time tracking. It lets you start, stop, toggle, add, list and delete time entries
from the terminal instead of the web app, and is designed to be driven both
interactively and from scripts and status-bar widgets.

Soltty is an independent project. It talks to Solidtime only over its public HTTP
REST API and contains no Solidtime code.

## Tech Stack

- **Go 1.21+** - primary language
- **github.com/spf13/cobra** - command structure, flags and help output
- **github.com/pkg/browser** - opening the Solidtime web UI from `soltty web`
- **net/http, encoding/json** - API client (stdlib, no HTTP framework)
- **Nix flake** - reproducible dev shell and package
- **GoReleaser + GitHub Actions** - multi-platform release builds

No runtime dependencies beyond the two modules above; everything else is stdlib.

## Project Conventions

### Code Style

- One cobra command per file in `cmd/` (`start.go`, `stop.go`, `add.go`, ...)
- Command wiring in `init()`, behaviour in a `runX` function
- English identifiers, comments and user-facing strings
- Exported API types live in `internal/client`, config handling in `internal/config`
- Errors are returned up and printed once by the command, via `formatError`
- `gofmt` is enforced in CI; run `gofmt -w .` before committing

### Architecture Patterns

- Thin CLI layer over a small hand-written API client - no code generation
- `internal/client` owns HTTP, auth and JSON decoding; `cmd/` owns presentation
- Shared helpers (time parsing, duration and date formatting) live in `cmd/utils.go`
- Machine-readable output is opt-in per command via `--json`
- Confirmation prompts are skippable with `--yes` so scripts never block

### Error Handling

- User-facing errors explain the fix, not just the failure
- Config errors name every location searched and show a template
- Parse errors list the accepted input formats
- Commands print the error and return; they do not panic

### Data Flow

1. **Input**: flags and arguments, plus `config.json` (token, workspace, base URL)
2. **Client**: authenticated REST calls to the Solidtime API
3. **Output**: human-readable text by default, JSON with `--json`

### File Structure

```
.
├── main.go                    # Entry point, calls cmd.Execute()
├── cmd/                       # One file per command + shared helpers
│   ├── root.go                # Root command, client construction, error format
│   ├── start.go stop.go       # Timer lifecycle
│   ├── toggle.go continue.go  # Widget-oriented and resume commands
│   ├── add.go delete.go       # Completed entries
│   ├── current.go list*.go    # Queries
│   ├── web.go info.go         # Browser and account info
│   └── utils.go               # Time parsing, duration/date formatting
├── internal/client/           # Solidtime API client
├── internal/config/           # Config discovery and validation
├── scripts/                   # Maintenance scripts
├── THIRD_PARTY_LICENSES/      # License texts of linked dependencies
├── openspec/                  # Specs and change proposals
└── .beans/                    # Task tracking
```

## Domain Context

### Solidtime

- Open-source time tracking application, AGPL-3.0, self-hostable
- Soltty targets any instance; the base URL is configured per user
- Authentication: `Authorization: Bearer {api_token}`
- Relevant endpoints, all scoped to an organization (workspace):
  - `/organizations/{workspace_id}/time-entries`
  - `/organizations/{workspace_id}/projects`
  - `/organizations/{workspace_id}/clients`
  - `/users/me/time-entries/active`
- Timestamps are exchanged in UTC (RFC3339) and displayed in local time

### Time entries

- A running entry has `end: null`; a completed entry has both `start` and `end`
- Entries carry an optional `project_id`; projects carry an optional `client_id`
- Entry IDs are UUIDs, addressable by a 6-36 character prefix (see the
  `short-id-matching` capability)

## Important Constraints

- **`base_url` is required** in `config.json` - there is no default instance
- All timestamps are converted to UTC before sending, and to local time for display
- Input without timezone information is interpreted as local time; input without a
  date resolves to today (see the `solidtime-cli` capability, "Time Input Formats")
- `config.json` holds an API token and is gitignored - it must never be committed
- The binary is statically linked, so every release redistributes its dependencies;
  their license texts must ship with it (`THIRD_PARTY_LICENSES/`)
- Soltty must not bundle or link Solidtime code - API access only

## External Dependencies

### Solidtime API

- Base URL: configured per user, e.g. `https://app.example.com/api/v1`
- Authentication: `Bearer {api_token}`
- Accept: `application/vnd.api+json`

## Setup

### Prerequisites

- Go 1.21 or higher, or Nix with flakes enabled
- A Solidtime instance and an API token

### Initial setup

1. Clone the repository
2. Create the config:
   ```bash
   mkdir -p ~/.config/soltty
   cp config.json.example ~/.config/soltty/config.json
   ```
3. Fill in `api_token`, `workspace_id`, `base_url` and `username`
4. Build and verify:
   ```bash
   go build ./... && go test ./...
   soltty info
   ```

### Development

```bash
nix develop          # dev shell with go, gopls, gotools
go test ./...        # tests
gofmt -l .           # must print nothing; CI enforces this
scripts/update-third-party-licenses.sh --check
```

## Development Notes

- Changes follow the OpenSpec workflow: propose, apply, archive
- `release.sh` bumps `VERSION`, updates the changelog, refreshes the Nix
  `vendorHash` and creates the tag; pushing the tag triggers GoReleaser
- The `## NEXT VERSION` changelog heading is replaced by `release.sh`
