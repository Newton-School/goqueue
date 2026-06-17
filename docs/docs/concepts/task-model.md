---
title: Task Model
---

goqueue tasks are data + metadata + runtime policy.

## Core task objects

- `TaskName`: stable task identifier.
- `QueueName`: target queue for a single dispatch.
- `TaskID`: random RFC4122 task ID.
- `TaskPayload`: args and kwargs passed to handler.
- `TaskMetadata`: map of workflow and user-defined metadata.
- `TaskEnvelope`: full validated message payload used by backend.

## Timing and scheduling

- `TaskTiming.ETA`: when a task becomes executable.
- `TaskTiming.ExpiresAt`: hard deadline after which task is marked expired.
- Use `goqueue.TaskTimingFromCountdown(now, duration)` for delayed starts.

## Priority and retries

- `Priority` range: `0..9` (`0` is default).
- `RetryPolicy` fields:
  - `MaxAttempts` (`>= 1`, default `1`)
  - `Backoff`
  - `MaxBackoff`
- `DelayForAttempt(attempt)` doubles backoff per attempt with cap.

## Result shape

Tasks return:

- `task.TaskResult` with:
  - `State`: terminal task state
  - `Value`: optional output
  - `Error`: optional error text
  - `Metadata`: optional map

## Task states

- `PENDING`, `SCHEDULED`, `RECEIVED`, `STARTED`, `RETRYING`
- `SUCCEEDED`, `FAILED`, `REVOKED`, `EXPIRED`, `DEAD_LETTERED`
