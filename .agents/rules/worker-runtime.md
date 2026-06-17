# Worker Runtime Rules

Use this file when changing `worker/`, handler execution, retries,
dead-lettering, ack behavior, pending recovery, workflow advancement from
workers, or worker-facing options.

## Responsibilities

- Workers read ready queue messages through `backend.QueueBackend`.
- Workers decode messages into validated task envelopes with the configured
  payload codec.
- Workers execute registered task handlers from `task.TaskRegistry`.
- Workers persist task state and results, schedule retries, dead-letter
  unrecoverable tasks, advance workflows, and acknowledge messages only after
  durable side effects succeed.

## Processing Flow

1. Ensure the consumer group exists.
2. Optionally move due scheduled tasks into the ready stream.
3. Optionally claim stale pending messages.
4. Read ready messages in bounded batches.
5. Decode the task message; malformed messages go to DLQ when DLQ is enabled.
6. Write `TaskReceived`.
7. Check expiration before handler execution.
8. Look up the handler. Unknown tasks are dead-lettered.
9. Write `TaskStarted`.
10. Execute the handler and normalize the result.
11. Schedule retry, mark retry exhaustion, or persist final state/result.
12. Advance workflow state for successful or terminal workflow tasks.
13. Ack the stream entry after the durable update path succeeds.

## Ack Ordering

- Do not ack before required state, result, retry scheduling, DLQ persistence, or
  workflow advancement completes.
- If retry scheduling fails, attempt to dead-letter the task and return the
  scheduling error.
- If a task expires before retry, persist expired state/result and ack.
- Preserve message ownership semantics for pending recovery; do not ack claimed
  messages until they are processed.

## Concurrency And Shutdown

- Keep concurrency bounded with the worker semaphore.
- Respect context cancellation and wait for in-flight goroutines before
  returning.
- Keep read batch, claim batch, move-due limit, block duration, idle delay,
  pending min idle, and claim interval validated.
- Use deterministic clocks in tests through worker options.

## Results And Failures

- Normalize invalid handler results to failed results.
- Only terminal states or retrying are accepted from handlers.
- Failure metadata should identify category, attempt, max attempts, retryable
  status, next retry time, DLQ status, and last error where relevant.
- Do not leak secrets through result errors or metadata.

## Documentation

- Update docs for execution semantics, retries, DLQ transitions, and worker
  state progression when behavior changes.
- If worker lifecycle output changes affect users, update `docs/docs/reference/errors.md`
  and relevant concept pages under `docs/docs/concepts/worker.md`.
- Use `docs-sync.md` for coordinating all user-facing doc updates in the same
  commit.

## Tests

- Cover ack ordering for success, failure, retry, retry exhaustion, expiration,
  unknown task, malformed message, DLQ disabled/enabled, pending recovery, and
  workflow advancement.
- Include cancellation and concurrency tests when changing the run loop.
- Use fake backends for ordering assertions and Redis integration tests for
  stream/pending behavior.
