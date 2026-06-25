## ADDED Requirements

### Requirement: current command supports JSON output
The `soltty current` command SHALL accept a `--json` flag that outputs the timer state as a single JSON object to stdout instead of human-readable text.

#### Scenario: JSON output when timer is running
- **WHEN** user runs `soltty current --json` and a timer is active
- **THEN** the command outputs a JSON object with `"running": true`, `"id"` (full UUID), `"description"`, `"project_id"` (null if none), `"project"` (resolved name or null), `"elapsed"` (formatted string e.g. `"2h 15m"`), and `"start"` (RFC3339 UTC timestamp)

#### Scenario: JSON output when no timer is running
- **WHEN** user runs `soltty current --json` and no timer is active
- **THEN** the command outputs `{"running": false}` and exits with code 0

### Requirement: list command supports JSON output
The `soltty list` command SHALL accept a `--json` flag as a persistent flag inherited by all subcommands (`list clients`, `list projects`), outputting a JSON array instead of a table.

#### Scenario: list --json outputs time entries array
- **WHEN** user runs `soltty list --json`
- **THEN** the command outputs a JSON array of time entry objects, each containing `"id"` (full UUID), `"description"`, `"project_id"` (null if none), `"project"` (resolved name or null), `"date"` (YYYY-MM-DD), `"start"` (RFC3339), `"end"` (RFC3339 or null if running), `"duration"` (integer seconds), `"running"` (boolean)

#### Scenario: list --json respects --limit flag
- **WHEN** user runs `soltty list --json --limit 5`
- **THEN** the JSON array contains at most 5 entries

#### Scenario: list --json uses full UUIDs
- **WHEN** user runs `soltty list --json`
- **THEN** all `"id"` fields contain full 36-character UUIDs, not 8-character short IDs

#### Scenario: list --json with no entries
- **WHEN** user runs `soltty list --json` and there are no time entries
- **THEN** the command outputs `[]` and exits with code 0

### Requirement: list clients supports JSON output
The `soltty list clients` command SHALL output a JSON array when the `--json` flag is provided (inherited from parent `list` command).

#### Scenario: list clients --json outputs client array
- **WHEN** user runs `soltty list clients --json`
- **THEN** the command outputs a JSON array of active client objects, each containing `"id"` (full UUID), `"name"`, and `"project_count"` (integer count of active projects)

#### Scenario: list clients --json with no clients
- **WHEN** user runs `soltty list clients --json` and no active clients exist
- **THEN** the command outputs `[]` and exits with code 0

### Requirement: list projects supports JSON output
The `soltty list projects` command SHALL output a JSON array when the `--json` flag is provided (inherited from parent `list` command).

#### Scenario: list projects --json outputs project array
- **WHEN** user runs `soltty list projects --json`
- **THEN** the command outputs a JSON array of active project objects, each containing `"id"` (full UUID), `"name"`, `"client_id"` (null if no client), and `"client"` (resolved client name or null)

#### Scenario: list projects --json respects --client filter
- **WHEN** user runs `soltty list projects --json --client "TMCS"`
- **THEN** the JSON array contains only projects belonging to clients matching the filter

#### Scenario: list projects --json with no matching projects
- **WHEN** user runs `soltty list projects --json` and no active projects exist
- **THEN** the command outputs `[]` and exits with code 0
