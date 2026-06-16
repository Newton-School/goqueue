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

- [ ] 1. Create `admin` package and exported error surface.
- [ ] 2. Add control backend capabilities for task message reads and dead-letter single read/delete.
- [ ] 3. Add admin request/result models for task retry, revoke, and dead-letter replay.
- [ ] 4. Implement `Admin` constructor and runtime nil/backend guards.
- [ ] 5. Implement `TaskMessage` read-through helper in admin.
- [ ] 6. Implement `RetryTask` operation with attempt reset and queue override.
- [ ] 7. Implement optional scheduling for retry when `retry-after` is provided.
- [ ] 8. Implement `RevokeTask` operation with state transition to `TaskRevoked`.
- [ ] 9. Implement `PurgeQueue` operation in admin surface.
- [ ] 10. Implement dead-letter replay operation with optional source deletion.
- [ ] 11. Implement dead-letter bulk delete operation.
- [ ] 12. Add request validation and error wrapping for all control operations.
- [ ] 13. Add `redisbackend` implementations for all control backend methods.
- [ ] 14. Add unit tests for redisbackend control primitives.
- [ ] 15. Add `App.NewAdmin` helper returning configured admin client.
- [ ] 16. Export admin types and errors from root facade.
- [ ] 17. Add `goqueue control ...` CLI group.
- [ ] 18. Add CLI actions for `retry-task`, `revoke-task`, `replay-dead-letter`,
  `delete-dead-letter`, and `purge-queue`.
- [ ] 19. Add CLI tests for argument parsing and JSON/text rendering branches.
- [ ] 20. Add phase 9 docs for operations APIs.
- [ ] 21. Update README with new operations CLI usage and roadmap status.
- [ ] 22. Update root docs for control-plane safety defaults.
- [ ] 23. Run full `go test ./...` and `go vet ./...`.

## Acceptance Criteria

- Production users can read a task payload, requeue it, and force revoke state.
- Operators can replay and delete dead-letter records.
- Operators can purge queue storage with explicit intent and clear/retain message
  state/result artifacts.
- CLI control commands fail-fast on bad IDs/queue names and return structured output.
- Control operations are separated from inspection APIs and do not mutate unless
  explicitly requested.
- `go test ./...` and `go vet ./...` pass.

