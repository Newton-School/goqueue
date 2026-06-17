---
title: CLI
---

Command line entrypoint is `goqueue`.

## Install

```bash
go install github.com/Newton-School/goqueue/cmd/goqueue@latest
```

## Base commands

- `goqueue inspect ...` for read-style operations
- `goqueue control ...` for mutating operations

## Inspect commands

- `goqueue inspect ping`
- `goqueue inspect stats --queue <queue>`
- `goqueue inspect deadletters --queue <queue> [--count <n>]`
- `goqueue inspect task --id <task_id> [--json]`
- `goqueue inspect state --id <task_id> [--json]`
- `goqueue inspect result --id <task_id> [--json]`
- `goqueue inspect forget-result --id <task_id>`

## Control commands

- `goqueue control retry-task --id <task_id> [--queue <queue>] [--scheduled-at <RFC3339>] [--countdown <duration>] [--preserve-attempt] [--clear-state] [--clear-result]`
- `goqueue control revoke-task --id <task_id> [--reason <text>]`
- `goqueue control replay-dead-letter --queue <queue> --stream-id <id> [--destination-queue <queue>] [--delete-source]`
- `goqueue control delete-dead-letter --queue <queue> --stream-id <id>[,<id>...]`
- `goqueue control purge-queue --queue <queue> --yes [--delete-messages] [--delete-states] [--delete-results]`

## Global flags

All commands support:

- `--redis-url`
- `--namespace`
- `--json`
