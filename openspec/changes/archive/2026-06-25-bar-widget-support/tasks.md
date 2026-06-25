## 1. current --json

- [x] 1.1 Add `currentJSON bool` package-level var and `--json` flag registration in `cmd/current.go` `init()`
- [x] 1.2 Define `currentJSONOutput` struct with fields: `Running`, `ID`, `Description`, `ProjectID`, `Project`, `Elapsed`, `Start`
- [x] 1.3 In `runCurrent`: when `--json` is set and no timer is running, output `{"running":false}` and return
- [x] 1.4 In `runCurrent`: when `--json` is set and timer is running, resolve project name via `c.GetProjects()`, populate struct, marshal to JSON and print
- [x] 1.5 Verify `soltty current --json` outputs valid JSON in both running and stopped states

## 2. list --json (time entries)

- [x] 2.1 Add `listJSON bool` package-level var and register `--json` as a persistent flag on `listCmd` in `cmd/list.go` `init()`
- [x] 2.2 Define `listEntryJSON` struct with fields: `ID`, `Description`, `ProjectID`, `Project`, `Date`, `Start`, `End`, `Duration`, `Running`
- [x] 2.3 In `runList`: when `--json` is set, build JSON array using full UUIDs and resolved project names, marshal and print; skip table output
- [x] 2.4 Verify `soltty list --json` outputs JSON array, full UUIDs used, `[]` when no entries
- [x] 2.5 Verify `soltty list --json --limit 3` returns at most 3 entries

## 3. list clients --json

- [x] 3.1 Define `listClientJSON` struct with fields: `ID`, `Name`, `ProjectCount`
- [x] 3.2 In `runListClients`: check `listJSON` flag; when set, build JSON array from `activeClients` with project counts, marshal and print; skip text output
- [x] 3.3 Verify `soltty list clients --json` outputs JSON array, `[]` when no active clients

## 4. list projects --json

- [x] 4.1 Define `listProjectJSON` struct with fields: `ID`, `Name`, `ClientID`, `Client`
- [x] 4.2 In `runListProjects`: check `listJSON` flag; when set, build JSON array from `activeProjects` with resolved client names, marshal and print; skip table output
- [x] 4.3 Verify `soltty list projects --json` outputs JSON array with client info
- [x] 4.4 Verify `soltty list projects --json --client "TMCS"` filters correctly

## 5. --yes flag on start and continue

- [x] 5.1 Add `startYes bool` var and `--yes / -y` flag in `cmd/start.go` `init()`
- [x] 5.2 In `runStart`: when timer is running and `startYes` is true, stop the timer immediately without prompting; print the same "✓ Stopped" confirmation
- [x] 5.3 Add `continueYes bool` var and `--yes / -y` flag in `cmd/continue.go` `init()`
- [x] 5.4 In `runContinue`: when timer is running and `continueYes` is true, stop the timer immediately without prompting; print the same "✓ Stopped" confirmation
- [x] 5.5 Verify `soltty start --yes "Task"` stops running timer without prompt
- [x] 5.6 Verify `soltty continue --yes <id>` stops running timer without prompt
- [x] 5.7 Verify `soltty start "Task"` (without --yes) still shows Y/N prompt

## 6. toggle command

- [x] 6.1 Create `cmd/toggle.go` with `toggleCmd` cobra command, `--description / -d` and `--project / -p` flags
- [x] 6.2 In `runToggle`: call `c.GetCurrentTimeEntry()`; if timer running, stop it, print confirmation, exit 0
- [x] 6.3 In `runToggle` (stopped path): call `c.ListTimeEntries(1)`; if entry exists, start new timer with its description and project_id, print confirmation, exit 0
- [x] 6.4 In `runToggle` (no entries + description provided): resolve project by name if `--project` given, start new timer with provided description, print confirmation, exit 0
- [x] 6.5 In `runToggle` (no entries + no description): print `"No timer running and no previous entry to continue."` and `os.Exit(1)`
- [x] 6.6 Register `toggleCmd` in `cmd/root.go`
- [x] 6.7 Verify `soltty toggle` stops a running timer
- [x] 6.8 Verify `soltty toggle` continues the last entry when stopped
- [x] 6.9 Verify `soltty toggle --description "Task" --project "TMCS-General"` starts new timer when no entries exist
- [x] 6.10 Verify `soltty toggle` exits 1 with error message when no entries and no description

## 7. Documentation

- [x] 7.1 Update README: add `--json` flag to `current`, `list`, `list clients`, `list projects` examples
- [x] 7.2 Update README: add `--yes` flag to `start` and `continue` examples
- [x] 7.3 Update README: add `toggle` command to Features list and Usage section
- [x] 7.4 Update CHANGELOG.md under `## NEXT VERSION` with all new features
