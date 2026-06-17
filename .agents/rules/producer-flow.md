# Producer Flow Rules

Use this file when changing `producer/`, root producer helpers, scheduled task
publish behavior, enqueue options, or producer-facing examples.

## Responsibilities

- Producers convert user input into validated task envelopes and backend
  messages.
- Producers choose ready vs scheduled enqueue based on task timing.
- Producers persist an initial task state before storing the queue message.
- Producers should remain backend-agnostic and talk only through
  `backend.QueueBackend`.

## ApplyAsync Flow

1. Validate the task name and producer configuration.
2. Resolve apply options: queue override, task ID, metadata, priority, retry
   policy, countdown, ETA, expiration, attempt, and created time.
3. Resolve countdowns using the producer clock.
4. Build a validated task envelope.
5. Serialize the envelope through the configured payload codec.
6. Persist initial task state as `TaskPending` or `TaskScheduled`.
7. Enqueue through `EnqueueReady` or `EnqueueScheduled`.
8. Return `AsyncResult` for state/result lookup.

## Error Handling

- If enqueue fails after state persistence, mark the task failed when possible
  and include both enqueue and state-write errors when state update also fails.
- Keep context handling tolerant of nil contexts by replacing them with
  `context.Background()`.
- Validate option values at option application time where possible.

## Scheduling

- Countdown values must not be negative.
- ETA and expiration validation belongs in `task.TaskTiming`.
- Scheduled tasks should not enter ready queues until the backend due-task move
  path promotes them.

## Documentation

- Update concept docs for producer timing/enqueue semantics and any queueing
  behavior changes.
- If retry policy defaults or metadata behavior changes, refresh the relevant
  section under `docs/docs/concepts/task-model.md` and/or concept flow pages.
- If producer options affect setup, route config updates through `readme-policy.md`.
- When producer behavior changes, apply `docs-sync.md` before finalizing.

## Tests

- Cover ready enqueue, scheduled enqueue, initial state persistence, custom task
  IDs, custom queues, metadata, retry policy, priority, countdown/ETA behavior,
  codec failures, and enqueue failure state transitions.
- Use deterministic clocks through `WithProducerNow` for timing tests.
