# goqueue CLI

`goqueue` exposes a read-only CLI for inspection and operations.

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

## Common options

- `--redis-url`: Redis URL (defaults to `GOQUEUE_REDIS_URL` env var)
- `--namespace`: namespace (defaults to `GOQUEUE_NAMESPACE` or `goqueue`)
- `--json`: emit JSON output

## Response semantics

- Read methods return backend errors from the queue data layer.
- `forget-result` does not delete task payloads or state, only removes terminal result blobs.
