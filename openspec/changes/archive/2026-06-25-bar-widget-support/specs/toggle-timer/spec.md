## ADDED Requirements

### Requirement: toggle command acts as non-interactive play/pause
The `soltty toggle` command SHALL stop the running timer if one is active, or start/continue a timer if none is running. It is fully non-interactive and designed for bar widget use.

#### Scenario: toggle stops running timer
- **WHEN** user runs `soltty toggle` and a timer is currently running
- **THEN** the running timer is stopped, a confirmation is printed, and the command exits with code 0

#### Scenario: toggle continues last entry when stopped
- **WHEN** user runs `soltty toggle` and no timer is running but at least one previous entry exists
- **THEN** the most recent time entry is continued (same description and project), a confirmation is printed, and the command exits with code 0

#### Scenario: toggle starts new timer when stopped with description
- **WHEN** user runs `soltty toggle --description "Task name"` and no timer is running and no previous entry exists
- **THEN** a new timer is started with the provided description and optional project, a confirmation is printed, and the command exits with code 0

#### Scenario: toggle exits with error when no entry and no description
- **WHEN** user runs `soltty toggle` and no timer is running and no previous entries exist and no --description is provided
- **THEN** the command prints `"No timer running and no previous entry to continue."` and exits with code 1

### Requirement: toggle command accepts optional --description and --project flags
The `soltty toggle` command SHALL accept `--description / -d` and `--project / -p` flags used only when starting a fresh timer (no previous entries).

#### Scenario: toggle ignores --description when stopping
- **WHEN** user runs `soltty toggle --description "Something"` and a timer is running
- **THEN** the running timer is stopped (--description is ignored in stop mode)

#### Scenario: toggle ignores --description when continuing last entry
- **WHEN** user runs `soltty toggle --description "Something"` and no timer is running but a previous entry exists
- **THEN** the last entry is continued using its original description (--description is ignored in continue mode)
