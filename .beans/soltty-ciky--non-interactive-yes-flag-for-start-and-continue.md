---
# soltty-ciky
title: Non-interactive --yes flag for start and continue
status: completed
type: feature
priority: high
created_at: 2026-06-25T20:01:20Z
updated_at: 2026-06-26T06:48:29Z
---

Add a --yes / -y flag to 'start' and 'continue' commands to skip the Y/N confirmation
prompt when a timer is already running.

## Motivation
Bar widgets and shell scripts cannot handle interactive prompts. When a timer is running
and a new start/continue is triggered, the command currently hangs waiting for stdin.

## Proposed behaviour
With --yes: if a timer is running, stop it immediately and start the new one without prompting.

## Usage
```bash
soltty start --yes "General. Tasks not further specified" -p "TMCS-General"
soltty continue --yes f1ad5503
```

## Affects
- cmd/start.go
- cmd/continue.go
