# Operations Rules

Use this file when changing `inspect/`, `admin/`, `cmd/goqueue/`, operational
docs, queue health, task observability, retry/revoke/replay/delete/purge
controls, or CLI output.

## Inspect APIs

- Inspect APIs are read-only except explicit result-forget operations.
- Inspect should use backend contracts and avoid direct Redis access.
- Task snapshots should distinguish missing state/result from backend errors.
- Queue stats should expose ready, scheduled, and dead-letter counts without
  mutating queue contents.
- Health checks should prove backend reachability without leaking credentials.

## Admin APIs

- Admin APIs are mutating operational controls. Validate every task ID, queue
  name, stream ID, countdown, scheduled time, and option combination before
  changing backend state.
- `RetryTask` reloads persisted task messages, optionally changes queue/timing,
  optionally resets attempt, then enqueues ready or scheduled.
- `RevokeTask` writes revoked state with the operator reason.
- `ReplayDeadLetter` reads a DLQ entry, resets attempt, optionally changes
  destination queue, enqueues ready, and optionally deletes the source entry.
- `DeleteDeadLetters` requires at least one non-empty stream ID.
- `PurgeQueue` must require explicit caller intent at CLI level and must never
  delete outside the configured namespace.

## CLI Behavior

- Keep command names stable: `inspect` for read-oriented operations and
  `control` for mutating operations.
- Preserve `--json` for machine-readable output where it already exists.
- Write human-readable errors to stderr and exit non-zero on failure.
- Do not print Redis URLs with credentials.
- Require `--yes` for destructive purge commands.
- Prefer explicit flags over positional arguments for task IDs, queues, stream
  IDs, namespaces, Redis URLs, and output format.

## Documentation

- Update `README.md` and `docs/observability/` when CLI behavior or inspect
  output changes.
- Update `docs/reliability/` when retry, revocation, DLQ, pending recovery, or
  purge semantics change.

## Tests

- Cover success and validation failures for inspect/admin methods.
- Cover CLI parsing, required flags, JSON output, human output, error messages,
  and destructive confirmation behavior.
- Use fake backends for API semantics and Redis integration tests for
  end-to-end recovery flows.
