## Context

`parseTime` in `cmd/utils.go` is a two-branch function: try `time.RFC3339`, else try `"15:04"` and stamp it onto today in `time.Local`. It is shared by `add` (`--start`, `--end`) and `start` (`--time`), so any widening lands in both commands at once.

Two properties of the current code shape this design:

1. `add` parses `--start` and `--end` independently. That is harmless today because a bare `HH:MM` always lands on today, so both flags land on the same day by construction. Once a date can appear on one flag, independent parsing breaks.
2. The `add` confirmation formats with `Format("15:04")`, discarding the date. With date input becoming easy, that turns a mistyped day into silent data.

See proposal.md - Why for motivation.

## Goals / Non-Goals

**Goals:**

- One parsing routine whose accepted set can be stated in a sentence, not enumerated as trivia.
- Every date the tool fills in is either one the user typed or today. Never a date the tool inferred by shifting.
- Preserve the exact meaning of both currently accepted shapes.

**Non-Goals:**

- Relative expressions (`yesterday`, `-1d`, `mon`). Possible follow-up; deliberately out of scope here.
- A separate `--date` flag. The widened value formats make it redundant.
- Per-command parsing rules. `add` and `start` share one parser.
- Correct handling of DST-ambiguous wall-clock times. See Risks.

## Decisions

### Grid of layouts, not an ad-hoc list

The accepted set is the product of four parts rather than a hand-picked list:

| part      | values                        | absent means      |
|-----------|-------------------------------|-------------------|
| date      | `YYYY-MM-DD`                  | today             |
| separator | `T` or space                  | (time-only input) |
| time      | `HH:MM` or `HH:MM:SS`         | -                 |
| zone      | `Z` or `±HH:MM`, after `T` only | local timezone  |

Implemented as an ordered slice of Go layouts tried in sequence, longest/most specific first. Adding a shape later means adding a layout, not restructuring a branch.

**Date is `YYYY-MM-DD` only.** Rejected: also accepting `DD-MM-YYYY`. It reads naturally to a Dutch user, but `05-09` is unreadable without knowing which convention is in force, and accepting both makes `09-16-2026` the only thing standing between a typo and a wrong month. Restricting to ISO makes every non-ISO date notation a loud parse error.

**Zone only after `T`.** Rejected: allowing `2026-09-16 14:00+02:00`. The `T` form is what gets pasted from an API response, a log line, or a script - machine provenance, so it carries a zone. The space form is what a person types by hand, and nobody hand-types an offset. Dropping the combination removes two layouts and two scenarios that would never be exercised.

**An explicit zone is honoured, not overridden.** Rejected: forcing every input to local. `Z` and `±HH:MM` are unambiguous instants and someone who types one means it; overriding would also silently change the meaning of the documented `2026-03-31T14:00:00Z` examples. "Always local" therefore applies to input that carries *no* zone information - which is the actual gap.

### `--end` inherits the date of `--start`

Defaulting chains rather than branching:

```
--start with no date  ->  today
--end   with no date  ->  the date of --start
```

Because a date-less `--start` already resolves to today, `--start "09:00" --end "17:00"` keeps behaving exactly as it does now. The rule needs no special case for the all-bare-times path.

Alternatives considered:

- **Leave both flags independent.** Rejected: `--start "2026-09-16 23:00" --end "01:00"` would silently produce an end on today's date - a wrong entry with no error.
- **Reject the mixed form** (one flag dated, the other not, is an error). Consistent with the strictness elsewhere in this design, but it refuses the most natural shorthand (`--start "2026-09-16 09:00" --end "17:00"`) whose meaning is unambiguous. Inheritance gives the same safety without the friction.

### No midnight rollover

When the resolved end still falls on or before the start, that is an error - the end date is never shifted forward a day, even when `--end` carried no date of its own.

Rejected: rolling over. It is convenient for overnight work, but it converts the far more common typo into silent bad data:

```
intended:  --start "15:00" --end "16:00"
typed:     --start "15:00" --end "14:00"

  with rollover  ->  a 23-hour entry, accepted
  without        ->  "End time must be after start time"
```

A 23-hour entry in a timesheet is a worse outcome than typing an end date once. Overnight entries are rare in time tracking, and the error message carries the fix (`--end "2026-09-17 01:00"`), so the case is guided rather than blocked.

### Confirmation shows the date only when it is not today

Unconditional dates would add noise to the common same-day case; never showing one leaves date input unverifiable. Showing `YYYY-MM-DD HH:MM` exactly when the entry is not on today's date makes the date visible precisely when it is the thing that could be wrong - and matches the input notation, so what is echoed is what could have been typed.

## Risks / Trade-offs

- **DST-ambiguous or non-existent wall-clock times** (the repeated or skipped hour at a transition) -> `time.Date` with `time.Local` resolves these without signalling. Accepted: prompting or erroring would hurt every ordinary call to protect two hours a year. Documented as a known limitation.
- **The grid accepts more than the three shapes discussed as the minimum** (it adds optional seconds and the `T`-without-zone form) -> Each is a natural read of a value a user or a script might hold, and each is one layout string plus one test row. The cost of a surprise rejection is higher than the cost of a layout.
- **`start --time` gains the new layouts without being asked for** -> Not gated: divergent parsing between two flags that accept "a time" would itself be a defect. Called out in the proposal's Impact so it is not a surprise at review.
- **`parseTime` has no existing test coverage** -> The table test lands in the same change as the widening, covering every accepted layout and the rejections that carry meaning (non-ISO dates, zone after a space separator).
