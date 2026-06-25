## Context

Soltty is a Go CLI using cobra. All current output is human-readable text printed directly via `fmt.Printf`. Commands like `start` and `continue` use `bufio.NewReader(os.Stdin)` to prompt for confirmation when a timer is already running — this blocks forever in non-interactive environments.

The target use case is bar widgets (waybar, i3bar) and shell scripts that need to (a) read timer state as structured data and (b) trigger start/stop/toggle without stdin interaction.

## Goals / Non-Goals

**Goals:**
- Machine-parseable JSON output for `current`, `list`, `list clients`, `list projects`
- Non-blocking `--yes` flag on `start` and `continue`
- A single `toggle` command for play/pause widget integration
- Zero new external dependencies

**Non-Goals:**
- Streaming or watch mode for JSON output
- JSON output for `stop`, `add`, `delete`, or `info` commands
- Any API server or socket interface
- Changing the default (human-readable) output of any command

## Decisions

### D1: `--json` as persistent flag on `listCmd`

`listCmd` is the parent cobra command for `list clients` and `list projects`. A persistent flag defined on `listCmd` is automatically available to all subcommands. This avoids duplicating flag registration in `list_clients.go` and `list_projects.go`.

The flag value (`listJSON bool`) is read from the parent command's flag set inside each subcommand's run function using `cmd.Root().PersistentFlags()` — or simply declared as a package-level var and bound once on `listCmd`.

**Alternative considered:** Add `--json` independently to each subcommand. Rejected — three identical flag registrations with no benefit.

### D2: JSON structs defined inline per command, not in a shared package

Each command serialises its own purpose-built struct (e.g. `currentJSONOutput`, `listEntryJSON`) rather than reusing the internal API client structs. This keeps the JSON contract stable even if internal types change, and avoids exposing internal client fields (some contain unexported or irrelevant fields).

**Alternative considered:** Serialise `client.TimeEntry` directly. Rejected — the internal struct has unexported fields and includes raw API data not suitable for output (e.g. `WorkspaceID`). Also decouples the public JSON contract from the API client.

### D3: `--yes` replaces the stdin branch, not wraps it

In `runStart` and `runContinue`, the existing stdin prompt block is guarded by `if current != nil`. We add an inner check: `if yesFlag { /* stop immediately */ } else { /* existing prompt */ }`. The stop logic itself is unchanged — it calls `c.StopTimeEntry(current.ID)` either way.

**Alternative considered:** Refactor confirmation into a helper function. Rejected — would add abstraction for a two-branch `if` that is already clear. The task is to add a flag, not refactor.

### D4: `toggle` fetches the last entry via `ListTimeEntries(1)`

The toggle "continue last entry" path needs the most recent entry. Rather than a dedicated API endpoint, it reuses the existing `c.ListTimeEntries(1)` call. If the result is empty, it falls through to the "no entries" branch.

`toggle` internally calls the same stop/start logic as `start --yes` — it stops without prompting (toggle is always non-interactive by definition).

### D5: `elapsed` field in `current --json` is a pre-formatted string, not seconds

Bar widgets typically display the elapsed time directly. A formatted string (`"2h 15m"`) is more useful than raw seconds, and the formatting logic already exists in `formatElapsedTime()` in `cmd/utils.go`. Widgets that need raw seconds can compute from `start`.

## Risks / Trade-offs

- **JSON shape is now a public API contract** — once bar widgets depend on the field names, renaming fields is a breaking change. Mitigation: use clear, stable field names from the start; document them in the README.
- **`--json` on `list` suppresses the "No time entries found" message** — the command outputs `[]` instead. This is correct for machine consumers but could surprise interactive users who accidentally pass `--json`. Mitigation: the flag name makes intent clear; no mitigation needed.
- **toggle "continue last" uses `ListTimeEntries(1)`** — this fetches from the API every call. For a widget polling every few seconds this is an extra API call per toggle. Acceptable: toggle is user-triggered, not polled.

## Open Questions

None — all edge cases (no entries, no description) are resolved in the specs.
