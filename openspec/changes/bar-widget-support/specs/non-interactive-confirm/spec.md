## ADDED Requirements

### Requirement: start command accepts --yes flag to skip confirmation
The `soltty start` command SHALL accept a `--yes / -y` flag that automatically confirms stopping the currently running timer without prompting the user via stdin.

#### Scenario: start --yes stops running timer without prompt
- **WHEN** user runs `soltty start --yes "New task"` and a timer is currently running
- **THEN** the running timer is stopped immediately, the new timer is started, and no Y/N prompt is displayed

#### Scenario: start --yes with no running timer behaves normally
- **WHEN** user runs `soltty start --yes "New task"` and no timer is running
- **THEN** the new timer starts normally (--yes has no effect)

#### Scenario: start without --yes still prompts
- **WHEN** user runs `soltty start "New task"` (without --yes) and a timer is running
- **THEN** the command displays the Y/N confirmation prompt as before (existing behaviour unchanged)

### Requirement: continue command accepts --yes flag to skip confirmation
The `soltty continue` command SHALL accept a `--yes / -y` flag that automatically confirms stopping the currently running timer without prompting the user via stdin.

#### Scenario: continue --yes stops running timer without prompt
- **WHEN** user runs `soltty continue --yes <entry-id>` and a timer is currently running
- **THEN** the running timer is stopped immediately, the continued timer is started, and no Y/N prompt is displayed

#### Scenario: continue --yes with no running timer behaves normally
- **WHEN** user runs `soltty continue --yes <entry-id>` and no timer is running
- **THEN** the entry is continued normally (--yes has no effect)

#### Scenario: continue without --yes still prompts
- **WHEN** user runs `soltty continue <entry-id>` (without --yes) and a timer is running
- **THEN** the command displays the Y/N confirmation prompt as before (existing behaviour unchanged)
