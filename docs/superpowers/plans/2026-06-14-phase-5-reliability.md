# Phase 5 Reliability And Failure Semantics Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add production-grade worker reliability with dead-letter queues, retry exhaustion semantics, failure metadata, pending-entry recovery, and operational inspection APIs.

**Architecture:** Extend `backend.QueueBackend` with dead-letter and pending-claim contracts, implement those contracts in `redisbackend`, and keep worker policy in the `worker` package. Failure classification lives in `task` so producer, worker, and observability code can share stable metadata keys without Redis-specific coupling.

**Tech Stack:** Go 1.26, Redis Streams consumer groups, Redis sorted sets, go-redis/v9, standard Go tests.

---

### Task 1: Failure Metadata Domain

**Files:**
- Create: `task/failure_metadata.go`
- Create: `task/failure_metadata_test.go`
- Modify: `task_exports.go`

- [x] Add failure category constants for execution errors, malformed messages, unknown tasks, expired tasks, retry exhaustion, and retry scheduling failures.
- [x] Add `FailureMetadata` with stable `ToMap` output for task result metadata.
- [x] Re-export the type and constants from the root package.
- [x] Verify with `go test ./task`.
- [x] Commit each narrow addition.

### Task 2: Backend Reliability Contracts

**Files:**
- Create: `backend/deadletter.go`
- Create: `backend/deadletter_test.go`
- Create: `backend/pending.go`
- Create: `backend/pending_test.go`
- Modify: `backend/backend.go`
- Modify: `backend/stats.go`
- Modify: `backend/stats_test.go`

- [x] Add `EnqueueDeadLetter`, `ReadDeadLetters`, and `ClaimStaleReady` to `QueueBackend`.
- [x] Add validation for dead-letter enqueue/read requests.
- [x] Add validation for stale pending claim requests.
- [x] Add `DeadLetterCount` to queue stats.
- [x] Verify with `go test ./backend`.
- [x] Commit each contract and test slice.

### Task 3: Redis Dead-Letter Backend

**Files:**
- Modify: `redisbackend/keys.go`
- Create: `redisbackend/deadletter_codec.go`
- Create: `redisbackend/deadletter_codec_test.go`
- Create: `redisbackend/deadletter.go`
- Create: `redisbackend/deadletter_test.go`
- Modify: `redisbackend/stats.go`
- Modify: `redisbackend/stats_test.go`

- [ ] Add dead-letter stream key building.
- [ ] Encode/decode dead-letter records as JSON stream fields.
- [ ] Implement `EnqueueDeadLetter` with validation and Redis `XADD`.
- [ ] Implement `ReadDeadLetters` with bounded `XREVRANGE`.
- [ ] Include dead-letter stream length in queue stats.
- [ ] Verify with `go test ./redisbackend`.
- [ ] Commit each Redis behavior.

### Task 4: Redis Pending Recovery Backend

**Files:**
- Create: `redisbackend/pending.go`
- Create: `redisbackend/pending_test.go`
- Modify: `redisbackend/stream_parser.go`
- Modify: `redisbackend/stream_parser_test.go`

- [ ] Implement `ClaimStaleReady` using Redis `XAUTOCLAIM`.
- [ ] Parse claimed messages through existing ready-message parsing.
- [ ] Validate min-idle, count, and start-id behavior.
- [ ] Verify with `go test ./redisbackend`.
- [ ] Commit each recovery behavior.

### Task 5: Worker Reliability Policy

**Files:**
- Modify: `worker/options.go`
- Modify: `worker/worker.go`
- Create: `worker/deadletter.go`
- Create: `worker/recovery.go`
- Modify: `worker/worker_test.go`

- [ ] Add worker options for DLQ enablement, pending recovery enablement, pending claim interval, pending min idle, and pending claim batch.
- [ ] Add startup recovery polling before ready reads.
- [ ] Route malformed messages, unknown tasks, expired tasks, retry-exhausted failures, and retry-schedule failures to DLQ when enabled.
- [ ] Preserve strict ack semantics: only ack after state/result/retry/DLQ persistence succeeds.
- [ ] Persist failure metadata in task results.
- [ ] Verify with `go test ./worker`.
- [ ] Commit each worker behavior.

### Task 6: Public API And Docs

**Files:**
- Modify: `task_exports.go`
- Modify: `README.md`
- Modify: `doc.go`
- Modify: `docs/superpowers/plans/2026-06-14-phase-5-reliability.md`

- [ ] Re-export public reliability options and metadata types.
- [ ] Document Phase 5 completion and operational semantics.
- [ ] Mark this plan complete.
- [ ] Verify with `go test ./...` and `go vet ./...`.
- [ ] Commit final docs and verification.

### Acceptance Criteria

- Worker sends poisoned/unknown/exhausted/expired unrecoverable tasks to a Redis-backed DLQ.
- Worker never acknowledges a message before final state/result/retry/DLQ persistence succeeds.
- Worker can claim stale pending stream messages for pod-crash recovery.
- Redis queue stats include dead-letter counts.
- Public APIs expose reliability options and failure metadata constants.
- `go test ./...` and `go vet ./...` pass.
- Phase 5 creates at least 100 commits from the pre-phase baseline.
