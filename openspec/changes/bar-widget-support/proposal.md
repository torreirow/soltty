## Why

Bar widgets (waybar, i3bar) and shell scripts cannot handle interactive stdin prompts and need machine-parseable output. Currently soltty hangs waiting for user input when a timer is already running, and all output is human-readable text — making it unusable in automated contexts.

## What Changes

- Add `--json` flag to `soltty current` for machine-parseable timer state
- Add `--json` persistent flag to `soltty list` (covers `list`, `list clients`, `list projects` subcommands)
- Add `--yes / -y` flag to `soltty start` and `soltty continue` to skip the Y/N confirmation prompt
- Add new `soltty toggle` command: stops if running, continues last entry if stopped (play/pause for widgets)

## Capabilities

### New Capabilities

- `json-output`: Machine-readable JSON output for `current`, `list`, `list clients`, and `list projects` commands
- `non-interactive-confirm`: `--yes` flag on `start` and `continue` to bypass stdin confirmation when a timer is already running
- `toggle-timer`: New `toggle` command that acts as a non-interactive play/pause — stops the running timer or continues the last entry

### Modified Capabilities

<!-- No existing spec-level requirements are changing -->

## Impact

**Files to modify:**
- `cmd/current.go` — add `--json` flag
- `cmd/list.go` — add `--json` as persistent flag, JSON output path for time entries
- `cmd/list_clients.go` — JSON output using inherited persistent flag
- `cmd/list_projects.go` — JSON output using inherited persistent flag
- `cmd/start.go` — add `--yes / -y` flag
- `cmd/continue.go` — add `--yes / -y` flag
- `cmd/root.go` — register toggle command

**New files:**
- `cmd/toggle.go` — toggle command implementation

**Dependencies:** No new external dependencies (uses stdlib `encoding/json`)
