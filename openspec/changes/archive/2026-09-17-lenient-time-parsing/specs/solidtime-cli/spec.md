## ADDED Requirements

### Requirement: Time Input Formats

The CLI tool SHALL accept time values for `--start`, `--end` and `--time` in a fixed set of formats built from four parts: an optional `YYYY-MM-DD` date, a `T` or space separator, a `HH:MM` or `HH:MM:SS` time, and an optional timezone offset that is only valid after a `T` separator. Input carrying no timezone information SHALL be interpreted in the user's local timezone. Input carrying no date SHALL resolve to today in the user's local timezone. Any other notation SHALL be rejected.

#### Scenario: Full RFC3339 timestamp
- **WHEN** user provides `2026-09-16T14:00:00Z` or `2026-09-16T14:00+02:00`
- **THEN** the CLI SHALL parse the value including its timezone offset
- **AND** preserve that offset when converting to UTC for the API

#### Scenario: Date and time with T separator, no timezone
- **WHEN** user provides `2026-09-16T14:00` or `2026-09-16T14:00:00`
- **THEN** the CLI SHALL interpret the value in the user's local timezone

#### Scenario: Date and time with space separator
- **WHEN** user provides `2026-09-16 14:00` or `2026-09-16 14:00:00`
- **THEN** the CLI SHALL interpret the value in the user's local timezone

#### Scenario: Time only
- **WHEN** user provides `14:00` or `14:00:00`
- **THEN** the CLI SHALL interpret the value as today in the user's local timezone

#### Scenario: Timezone offset after a space separator is rejected
- **WHEN** user provides `2026-09-16 14:00+02:00`
- **THEN** the CLI SHALL display a parse error listing the accepted formats
- **AND** exit with code 1

#### Scenario: Non-ISO date notation is rejected
- **WHEN** user provides `16-09-2026 14:00` or `09-16-2026 14:00`
- **THEN** the CLI SHALL display a parse error listing the accepted formats
- **AND** exit with code 1
- **AND** SHALL NOT interpret the value as any date

#### Scenario: Parse error names the accepted formats
- **WHEN** a time value cannot be parsed
- **THEN** the error message SHALL state that the date must be `YYYY-MM-DD`
- **AND** show that the separator may be `T` or a space, that seconds are optional, that a timezone offset is only allowed after `T`, and that a bare time resolves to today

## MODIFIED Requirements

### Requirement: Add Completed Entry

The CLI tool SHALL provide an `add` command to create time entries with specific start and end times. When `--end` carries no date, it SHALL inherit the date of `--start`.

#### Scenario: Add entry with explicit times
- **WHEN** user runs `solidtime-cli add "Meeting" --start "2026-03-31T14:00:00Z" --end "2026-03-31T15:30:00Z"`
- **THEN** the CLI SHALL create a completed time entry
- **AND** calculate duration as 90 minutes (5400 seconds)
- **AND** display confirmation with calculated duration

#### Scenario: Add entry with time-only format
- **WHEN** user runs `solidtime-cli add "Code review" --start "14:00" --end "15:30"`
- **THEN** the CLI SHALL interpret times as today in local timezone
- **AND** create entry with full ISO8601 timestamps
- **AND** display confirmation

#### Scenario: Add backdated entry with space separator
- **WHEN** user runs `solidtime-cli add "Client call" --start "2026-09-16 14:00" --end "2026-09-16 15:30"`
- **THEN** the CLI SHALL interpret both times as local timezone on 2026-09-16
- **AND** create a completed time entry with duration 90 minutes
- **AND** display confirmation

#### Scenario: End time inherits the start date
- **WHEN** user runs `solidtime-cli add "Client call" --start "2026-09-16 09:00" --end "17:00"`
- **THEN** the CLI SHALL resolve `--end` to 17:00 on 2026-09-16, not on today's date
- **AND** create a completed time entry with duration 8 hours

#### Scenario: Entry crossing midnight requires an explicit end date
- **WHEN** user runs `solidtime-cli add "Deploy" --start "2026-09-16 23:00" --end "01:00"`
- **THEN** the CLI SHALL resolve `--end` to 01:00 on 2026-09-16
- **AND** SHALL NOT shift the end date forward by a day
- **AND** display an error that the end time must be after the start time
- **AND** the error SHALL suggest giving `--end` with an explicit date, for example `--end "2026-09-17 01:00"`
- **AND** exit with code 1

#### Scenario: Add entry spanning midnight with explicit dates
- **WHEN** user runs `solidtime-cli add "Deploy" --start "2026-09-16 23:00" --end "2026-09-17 01:00"`
- **THEN** the CLI SHALL create a completed time entry with duration 2 hours
- **AND** display confirmation

#### Scenario: Add entry with project
- **WHEN** user runs `solidtime-cli add "Sprint planning" --start "10:00" --end "12:00" --project "Acme-Meetings"`
- **THEN** the CLI SHALL resolve project name to ID
- **AND** create entry with project attached
- **AND** display confirmation with all metadata

#### Scenario: Add entry with end before start
- **WHEN** user runs `solidtime-cli add "Task" --start "15:00" --end "14:00"`
- **THEN** the CLI SHALL display an error that the end time must be after the start time
- **AND** SHALL NOT create an entry by shifting the end time to the next day
- **AND** exit with code 1

#### Scenario: Add entry with missing time flags
- **WHEN** user runs `solidtime-cli add "Task" --start "14:00"`
- **AND** no --end flag is provided
- **THEN** the CLI SHALL display an error "Both --start and --end are required for add command"
- **AND** exit with code 1

#### Scenario: Confirmation shows the date for an entry not on today
- **WHEN** the CLI confirms an added entry whose start date is not today
- **THEN** the confirmation SHALL show start and end as `YYYY-MM-DD HH:MM` in local timezone

#### Scenario: Confirmation omits the date for an entry on today
- **WHEN** the CLI confirms an added entry whose start date is today
- **THEN** the confirmation SHALL show start and end as `HH:MM` in local timezone

### Requirement: Start Time Entry

The CLI tool SHALL provide a `start` command to begin tracking time.

#### Scenario: Start with description only
- **WHEN** user runs `solidtime-cli start "Working on feature X"`
- **THEN** the CLI SHALL create a new time entry with start time set to now
- **AND** display confirmation showing entry ID and start time
- **AND** the entry SHALL have end time set to null (running state)

#### Scenario: Start with project
- **WHEN** user runs `solidtime-cli start "Bug fix" --project "Example-Project"`
- **THEN** the CLI SHALL lookup the project ID from the project name
- **AND** create a time entry with the resolved project_id
- **AND** display confirmation with project name

#### Scenario: Start with custom start time
- **WHEN** user runs `solidtime-cli start "Forgot to start" --time "09:00"`
- **THEN** the CLI SHALL parse the time as today at 09:00 in local timezone
- **AND** create a time entry with start time set to 09:00
- **AND** display confirmation showing custom start time

#### Scenario: Start with ISO8601 start time
- **WHEN** user runs `solidtime-cli start "Task" --time "2026-03-31T08:00:00Z"`
- **THEN** the CLI SHALL parse the full ISO8601 timestamp
- **AND** create a time entry with the specified start time
- **AND** display confirmation

#### Scenario: Start with dated local time
- **WHEN** user runs `solidtime-cli start "Forgot to start" --time "2026-09-16 09:00"`
- **THEN** the CLI SHALL parse the value as 09:00 local time on 2026-09-16
- **AND** create a time entry with that start time
- **AND** display confirmation showing custom start time

#### Scenario: Start with unknown project
- **WHEN** user runs `solidtime-cli start "Task" --project "NonExistent"`
- **AND** the project name does not match any project in the workspace
- **THEN** the CLI SHALL display an error message
- **AND** list available project names as suggestions
- **AND** exit with code 1

#### Scenario: Start when timer already running - user confirms stop
- **WHEN** user runs `solidtime-cli start "New task"`
- **AND** a timer is already running
- **AND** user confirms stopping the running timer
- **THEN** the CLI SHALL stop the running timer
- **AND** start the new timer
- **AND** display confirmation for both actions

#### Scenario: Start when timer already running - user declines stop
- **WHEN** user runs `solidtime-cli start "New task"`
- **AND** a timer is already running
- **AND** user declines stopping the running timer
- **THEN** the CLI SHALL leave the running timer untouched
- **AND** SHALL NOT start a new timer
- **AND** display a message that no new timer was started

#### Scenario: Start when timer already running - stop fails
- **WHEN** user runs `solidtime-cli start "New task"`
- **AND** a timer is already running
- **AND** user confirms stopping the running timer
- **AND** stopping the running timer fails
- **THEN** the CLI SHALL display the stop error
- **AND** SHALL NOT start a new timer
- **AND** exit with code 1

### Requirement: Timezone Handling

The CLI tool SHALL handle timezone conversions between user's local time and API UTC requirements.

#### Scenario: Send times to API in UTC
- **WHEN** the CLI creates or updates a time entry
- **THEN** it SHALL convert all timestamps to UTC before sending to API
- **AND** format times as ISO8601/RFC3339
- **AND** include billable field (default: false)

#### Scenario: Display times in local timezone
- **WHEN** the CLI displays time information to the user
- **THEN** it SHALL convert UTC times from API to user's local timezone
- **AND** format times in user-friendly format (HH:MM or YYYY-MM-DD HH:MM)

#### Scenario: User inputs local time
- **WHEN** user provides time in HH:MM format (e.g., "14:00")
- **THEN** the CLI SHALL interpret it as local timezone
- **AND** convert to UTC for API submission

#### Scenario: User inputs ISO8601 time
- **WHEN** user provides time in ISO8601 format with timezone
- **THEN** the CLI SHALL preserve the timezone information
- **AND** convert to UTC for API submission

#### Scenario: User inputs date and time without timezone
- **WHEN** user provides a value with a date but no timezone offset (e.g., "2026-09-16 14:00" or "2026-09-16T14:00")
- **THEN** the CLI SHALL interpret it as local timezone
- **AND** convert to UTC for API submission

#### Scenario: Local time falls in a DST transition
- **WHEN** a zone-less value names a wall-clock time that is ambiguous or non-existent due to a daylight saving transition
- **THEN** the CLI SHALL resolve it without prompting the user
- **AND** SHALL NOT fail the command
