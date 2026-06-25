---
# soltty-f23d
title: soltty toggle command for bar widgets
status: todo
type: feature
priority: normal
created_at: 2026-06-25T20:01:29Z
updated_at: 2026-06-25T20:01:29Z
blocked_by:
    - soltty-ciky
---

Add a 'toggle' command that stops the running timer if one is active, or starts/continues
if none is running. Fully non-interactive, designed for bar widget use.

## Motivation
Bar widgets need a single command to act as play/pause. The toggle command removes the need
for the widget to check state before acting.

## Proposed behaviour
- Timer running  → stop it (equivalent to soltty stop)
- Timer stopped, last-entry available → continue last entry (equivalent to soltty continue <last-id>)
- Timer stopped, no entries yet + description provided → start new timer

## Usage
```bash
# Stop if running, continue last if stopped
soltty toggle

# Stop if running, start with these params if stopped
soltty toggle --description "General. Tasks" --project "TMCS-General"
```

## Depends on
- soltty-ciky (--yes flag, same non-interactive stop logic needed internally)
