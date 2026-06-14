# Dead-Letter Queues

Phase 5 adds dead-letter queue support for worker failures that cannot be
processed safely on the ready stream.

## What Is Dead-Lettered

Workers write DLQ records for malformed task payloads, unknown task names,
retry-exhausted failures, retry scheduling failures, and expired tasks.

## Terminal State Rules

Retry-exhausted, unknown, malformed, and retry-schedule failures end as
`DEAD_LETTERED`. Expired tasks keep the `EXPIRED` state while still receiving a
DLQ record for inspection.
