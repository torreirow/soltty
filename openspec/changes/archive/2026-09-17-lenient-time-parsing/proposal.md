## Why

Backdating a time entry is unnecessarily painful. `parseTime` accepts only two shapes: a full RFC3339 timestamp (seconds and timezone suffix mandatory) or bare `HH:MM`, which always resolves to today. Everything in between — `2026-09-16 14:00`, `2026-09-16T14:00` — is rejected. The most common real-world case, "I forgot to log yesterday", forces the user to type `2026-09-16T14:00:00+02:00` twice. On top of that, the `add` confirmation prints only `15:04`, so an entry created on the wrong day is indistinguishable from one created today.

## What Changes

- Widen `parseTime` to accept a grid of layouts instead of two fixed ones:
  - date is always `YYYY-MM-DD`; no other date notation is accepted
  - the date/time separator may be `T` or a space
  - seconds are optional
  - a timezone offset is accepted only after a `T` separator
  - input without timezone information is interpreted in the **local** timezone
  - input without a date resolves to **today**, local timezone
- `add` lets a date-less `--end` inherit the date of `--start` instead of defaulting to today. No midnight rollover: an end that still falls before the start is an error, never a silently shifted day.
- The "End time must be after start time" error explains how to record an entry crossing midnight (give `--end` with an explicit date).
- The `add` confirmation prints the date alongside the time whenever the entry is not on today's date.
- `start --time` inherits the widened layouts automatically, since it shares `parseTime`.
- Add `cmd/utils_test.go` with a table test covering every accepted layout and representative rejections; `parseTime` currently has no test coverage.

No breaking changes: both currently accepted shapes keep their exact meaning, including the timezone carried by an explicit `Z` or offset.

## Capabilities

### New Capabilities
<!-- No new capabilities being introduced -->

### Modified Capabilities
- `solidtime-cli`: broaden the accepted time input formats for `add` and `start`, define local-timezone defaulting for zone-less input, define `--end` date inheritance in `add`, and require the date in the `add` confirmation when the entry is not today

## Impact

- **Files Modified**: `cmd/utils.go` (`parseTime`), `cmd/add.go` (end inheritance, error message, confirmation output), `README.md` (documented time formats)
- **Files Added**: `cmd/utils_test.go`
- **Indirectly Affected**: `cmd/start.go` uses `parseTime` and gains the new layouts without code changes
- **Backward Compatibility**: No breaking changes. `2026-03-31T14:00:00Z` and `14:00` behave exactly as before
- **Rejected Input**: date notations other than `YYYY-MM-DD` (e.g. `09-16-2026`, `16-09-2026`) now produce a clear parse error rather than being silently misread
- **Known Limitation**: during a DST transition a local wall-clock time can be ambiguous or non-existent; Go's `time.Date` resolves it without prompting
