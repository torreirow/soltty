## 1. Widen the time parser

- [x] 1.1 Replace the two-branch body of `parseTime` in `cmd/utils.go` with an ordered slice of layouts (`2006-01-02T15:04:05Z07:00`, `2006-01-02T15:04Z07:00`, `2006-01-02T15:04:05`, `2006-01-02T15:04`, `2006-01-02 15:04:05`, `2006-01-02 15:04`, `15:04:05`, `15:04`), parsing zone-less layouts in `time.Local` and stamping time-only layouts onto today; verify with `go build ./...`
- [x] 1.2 Expose whether the parsed value carried its own date (e.g. a second return value or a small result struct), so `add` can apply end-date inheritance; verify `go build ./...` and that `cmd/start.go` still compiles against the signature
- [x] 1.3 Rewrite the parse error message to name the accepted shapes — date `YYYY-MM-DD`, separator `T` or space, seconds optional, timezone offset only after `T`, bare time means today; verify by running `soltty add "x" --start "16-09-2026 14:00" --end "15:00"` and reading the output

## 2. Test the parser

- [x] 2.1 Add `cmd/utils_test.go` with a table test covering all eight accepted layouts, asserting the resolved instant and that zone-less input lands in `time.Local`; verify `go test ./cmd/ -run TestParseTime` passes
- [x] 2.2 Extend the table with rejections that carry meaning — `16-09-2026 14:00`, `09-16-2026 14:00`, `2026-09-16 14:00+02:00`, `2026/09/16 14:00`, `14` — asserting each returns an error; verify `go test ./cmd/` passes

## 3. End-date inheritance in `add`

- [x] 3.1 In `cmd/add.go`, when `--end` carried no date of its own, rebuild it on the date of the resolved `--start`; verify `soltty add "x" --start "2026-09-16 09:00" --end "17:00"` creates an 8-hour entry on 2026-09-16
- [x] 3.2 Confirm no rollover is applied — an end that still falls on or before the start stays an error; verify `soltty add "x" --start "15:00" --end "14:00"` errors instead of creating a 23-hour entry
- [x] 3.3 Extend the "End time must be after start time" message with the midnight hint naming an explicit dated `--end`; verify `soltty add "x" --start "2026-09-16 23:00" --end "01:00"` prints the hint

## 4. Confirmation output

- [x] 4.1 In `cmd/add.go`, print start and end as `2006-01-02 15:04` when the entry's start date is not today, and keep `15:04` when it is; verify by adding one entry for today and one backdated entry and comparing the two confirmations

## 5. Documentation

- [x] 5.1 Update the `add` and `start` command `Long` help text in `cmd/add.go` and `cmd/start.go` with examples of the new formats; verify `soltty add --help` and `soltty start --help` show them
- [x] 5.2 Update the README time-format documentation (the `add` examples around line 222 and the custom start-time section) to describe the four-part grid and the end-date inheritance rule; verify by reading the rendered section
- [x] 5.3 Add a `## NEXT VERSION` CHANGELOG entry under `### Changed` describing the widened time formats and the end-date inheritance; verify the section exists and follows Keep a Changelog format

## 6. Verification

- [x] 6.1 Run `go build ./... && go test ./...` and confirm both pass
- [x] 6.2 Walk the scenarios in `specs/solidtime-cli/spec.md` against the built binary — full RFC3339, dated space form, bare time, inheritance, midnight error, rejected date notations, both confirmation formats — and confirm each behaves as specified
