---
# soltty-prlz
title: soltty current --json flag
status: todo
type: feature
priority: high
created_at: 2026-06-25T20:01:11Z
updated_at: 2026-06-25T20:01:11Z
---

Add a --json flag to the 'current' command so output is machine-parseable.

## Motivation
Bar widgets (e.g. wayle custom module) and shell scripts need to parse the current timer
state without fragile regex on human-readable output.

## Proposed output (when timer running)
```json
{
  "running": true,
  "id": "f1ad5503-...",
  "description": "General. Tasks not further specified",
  "project_id": "...",
  "project": "TMCS-General",
  "elapsed": "8h 6m",
  "start": "2026-06-25T10:32:00Z"
}
```

## Proposed output (when stopped)
```json
{"running": false}
```

## Usage
```bash
soltty current --json
```
