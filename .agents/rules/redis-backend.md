# Redis Backend Rules

Use this file when changing `backend/`, `redisbackend/`, Redis key layout,
message codecs, Lua scripts, scheduled promotion, state/result persistence,
periodic coordination, dead-letter storage, pending recovery, or workflow Redis
state.

## Ownership

- `backend/` defines storage contracts. Keep it Redis-client-free.
- `redisbackend/` is the only package that should directly use Redis clients,
  Redis commands, Redis keys, Lua scripts, or Redis-specific encodings.
- Other packages should depend on backend interfaces and request/response
  structs instead of Redis implementation details.

## Key And Codec Rules

- Build Redis keys through `keyBuilder`; do not hand-format keys in operation
  code.
- Keep namespace validation strict and compatible with queue-name validation.
- Store task messages through the backend message codec and preserve raw payload
  bytes exactly.
- Keep state, result, periodic, dead-letter, and workflow codecs backward
  tolerant where operationally possible.
- Never log Redis URLs with credentials. Use redacted values only.

## Queue Storage

- Ready queues are Redis Streams.
- Scheduled queues are sorted sets plus retained task message keys.
- Due scheduled promotion must be atomic and bounded by a caller-provided limit.
- Ready and scheduled enqueue should persist the message with a positive TTL.
- Consumer group creation, stream reads, acknowledgements, and pending claims
  belong behind backend methods.

## Documentation

- Any Redis model/contract changes must be reflected in
  `docs/docs/concepts/redis-backend.md` and, when user-visible, `docs/docs/concepts/task-model.md`.
- If storage/error wording changes for users, sync `docs/docs/reference/errors.md`.
- Apply documentation updates through `docs-sync.md` for backend-facing behavior changes.

## Reliability Storage

- Dead-letter records must keep enough information for replay and debugging:
  queue, stream ID, worker group, consumer, task message, failure category,
  error, and timestamp.
- Pending recovery uses Redis pending-entry ownership transfer. Preserve
  `MinIdle`, batch limits, and start ID semantics.
- State and result TTL behavior must stay explicit and positive.
- Queue purge and cleanup paths must avoid deleting outside the configured
  namespace.

## Lua Scripts

- Use Lua when an operation needs Redis-side atomicity across multiple keys.
- Keep scripts deterministic and avoid unbounded loops.
- Add parser tests for every script result shape.
- Add integration tests when a script changes Redis state in a way mocks cannot
  prove.

## Tests

- Unit-test validation, key layout, codecs, parser behavior, and option errors.
- Redis integration tests should be gated by `GOQUEUE_RUN_INTEGRATION_TESTS`.
- Run package tests with `-count=1` for Redis integration fixes to avoid cached
  false positives.
