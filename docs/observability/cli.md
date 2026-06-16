# goqueue CLI

`goqueue` exposes a CLI for inspection and control operations.

## Install

From module root:

```bash
go install ./cmd/goqueue
```

## Commands

- `goqueue inspect ping`
- `goqueue inspect task --id <task-id> [--json]`
- `goqueue inspect state --id <task-id> [--json]`
- `goqueue inspect result --id <task-id> [--json]`
- `goqueue inspect forget-result --id <task-id>`
- `goqueue inspect deadletters --queue <queue> [--count <n>] [--json]`
- `goqueue inspect stats --queue <queue> [--json]`
- `goqueue control retry-task --id <task-id> [--queue <queue>] [--scheduled-at <RFC3339>] [--countdown <duration>] [--preserve-attempt] [--clear-state] [--clear-result] [--json]`
- `goqueue control revoke-task --id <task-id> [--reason <text>]`
- `goqueue control replay-dead-letter --queue <queue> --stream-id <id> [--destination-queue <queue>] [--delete-source] [--json]`
- `goqueue control delete-dead-letter --queue <queue> --stream-id <id>[,<id>...] [--json]`
- `goqueue control purge-queue --queue <queue> --yes [--delete-messages] [--delete-states] [--delete-results] [--json]`

## Common options

- `--redis-url`: Redis URL (defaults to `GOQUEUE_REDIS_URL` env var)
- `--namespace`: namespace (defaults to `GOQUEUE_NAMESPACE` or `goqueue`)
- `--json`: emit JSON output

## Response semantics

- Read methods return backend errors from the queue data layer.
- `forget-result` does not delete task payloads or state, only removes terminal result blobs.
