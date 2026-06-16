# Phase 9 Operations Control Plane

**Goal:** Add production-safe control APIs and command surface for queue/task
operations so operators can perform safe recovery and remediation actions.

**Architecture:** Add an `admin` package for task and queue control operations
with a strict `controlBackend` interface. Keep operations distinct from observation
APIs in `inspect` so readonly inspection and state mutations are cleanly separated.

## Phase 9 Scope

- Add task-level control APIs for retry/revoke operations.
- Add dead-letter replay and cleanup operations.
- Add queue purge/reset operations for recovery and incident handling.
- Add CLI `control` command group that reuses the same `admin` package.
- Add backend-level primitives needed by control APIs (message lookup, dead-letter
  deletes, queue purge helpers).
- Add operational docs and update public README with CLI command matrix.

## Work Items

- [x] 1. Create `admin` package and exported error surface.
- [x] 2. Add control backend capabilities for task message reads and dead-letter single read/delete.
- [x] 3. Add admin request/result models for task retry, revoke, and dead-letter replay.
- [x] 4. Implement `Admin` constructor and runtime nil/backend guards.
- [x] 5. Implement `TaskMessage` read-through helper in admin.
- [x] 6. Implement `RetryTask` operation with attempt reset and queue override.
- [x] 7. Implement optional scheduling for retry when `retry-after` is provided.
- [x] 8. Implement `RevokeTask` operation with state transition to `TaskRevoked`.
- [x] 9. Implement `PurgeQueue` operation in admin surface.
- [x] 10. Implement dead-letter replay operation with optional source deletion.
- [x] 11. Implement dead-letter bulk delete operation.
- [x] 12. Add request validation and error wrapping for all control operations.
- [x] 13. Add `redisbackend` implementations for all control backend methods.
- [x] 14. Add unit tests for redisbackend control primitives.
- [x] 15. Add `App.NewAdmin` helper returning configured admin client.
- [x] 16. Export admin types and errors from root facade.
- [x] 17. Add `goqueue control ...` CLI group.
- [x] 18. Add CLI actions for `retry-task`, `revoke-task`, `replay-dead-letter`,
  `delete-dead-letter`, and `purge-queue`.
- [x] 19. Add CLI tests for argument parsing and JSON/text rendering branches.
- [x] 20. Add phase 9 docs for operations APIs.
- [x] 21. Update README with new operations CLI usage and roadmap status.
- [x] 22. Update root docs for control-plane safety defaults.
- [x] 23. Run full `go test ./...` and `go vet ./...`.

## Acceptance Criteria

- Production users can read a task payload, requeue it, and force revoke state.
- Operators can replay and delete dead-letter records.
- Operators can purge queue storage with explicit intent and clear/retain message
  state/result artifacts.
- CLI control commands fail-fast on bad IDs/queue names and return structured output.
- Control operations are separated from inspection APIs and do not mutate unless
  explicitly requested.
- `go test ./...` and `go vet ./...` pass.

