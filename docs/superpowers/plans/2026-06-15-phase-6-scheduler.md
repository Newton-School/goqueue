# Phase 6 Scheduler And Periodic Jobs Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a production-grade Redis-backed scheduler that can register periodic task definitions, coordinate scheduler pods with leases, and enqueue due task instances through the existing producer path.

**Architecture:** Keep schedule domain types in a new `scheduler` package and storage contracts in `backend`. The Redis backend stores periodic definitions and next-run state, while the scheduler runtime leases due definitions before dispatching tasks through `producer.ApplyAsync`. Workers remain unchanged for execution and continue to consume ready/scheduled task messages.

**Tech Stack:** Go 1.26, Redis sorted sets and hashes, go-redis/v9, standard Go tests.

---

### Task 1: Scheduler Domain

**Files:**
- Create: `scheduler/errors.go`
- Create: `scheduler/schedule.go`
- Create: `scheduler/periodic_task.go`
- Create: `scheduler/schedule_test.go`
- Create: `scheduler/periodic_task_test.go`

- [x] Add interval schedule validation and next-time calculation.
- [x] Add periodic task definition validation with queue, task name, args, kwargs, metadata, priority, retry policy, and schedule.
- [x] Keep domain types Redis-independent.
- [x] Verify with `go test ./scheduler`.
- [x] Commit each narrow behavior.

### Task 2: Backend Schedule Contracts

**Files:**
- Create: `backend/periodic.go`
- Create: `backend/periodic_test.go`
- Modify: `backend/backend.go`
- Modify: `backend/backend_test.go`

- [x] Add `UpsertPeriodicTask`, `DeletePeriodicTask`, `ListDuePeriodicTasks`, and `MarkPeriodicTaskDispatched` backend contracts.
- [x] Add validation for periodic storage, due scans, and dispatch marking.
- [x] Update the backend interface acceptance test.
- [x] Verify with `go test ./backend`.
- [x] Commit each contract slice.

### Task 3: Redis Schedule Persistence

**Files:**
- Modify: `redisbackend/keys.go`
- Create: `redisbackend/periodic_codec.go`
- Create: `redisbackend/periodic_codec_test.go`
- Create: `redisbackend/periodic.go`
- Create: `redisbackend/periodic_test.go`
- Create: `redisbackend/periodic_integration_test.go`

- [x] Store periodic definitions as JSON records in Redis hashes.
- [x] Index next due times in a Redis sorted set.
- [x] Lease due definitions atomically so concurrent scheduler pods do not dispatch the same due occurrence.
- [ ] Mark successful dispatches by advancing the next due time.
- [ ] Verify with `go test ./redisbackend`.
- [ ] Commit each Redis behavior.

### Task 4: Scheduler Runtime

**Files:**
- Create: `scheduler/options.go`
- Create: `scheduler/scheduler.go`
- Create: `scheduler/dispatcher.go`
- Create: `scheduler/scheduler_test.go`
- Create: `scheduler/options_test.go`

- [ ] Add scheduler options for identity, poll interval, batch size, lock TTL, default queue, codec, and clock.
- [ ] Register periodic tasks through the backend with deterministic first-run state.
- [ ] Poll due definitions, dispatch task instances through producer, and mark only successful dispatches.
- [ ] Stop cleanly on context cancellation.
- [ ] Verify with `go test ./scheduler`.
- [ ] Commit each runtime behavior.

### Task 5: Public API And Documentation

**Files:**
- Modify: `app.go`
- Modify: `task_exports.go`
- Modify: `README.md`
- Modify: `doc.go`
- Create: `docs/scheduler/periodic-jobs.md`
- Modify: `docs/superpowers/plans/2026-06-15-phase-6-scheduler.md`

- [ ] Add `App.NewScheduler` with app defaults.
- [ ] Re-export scheduler types and options from the root package.
- [ ] Document periodic job registration, scheduler pod behavior, and Redis coordination.
- [ ] Mark this plan complete.
- [ ] Verify with `go test ./...` and `go vet ./...`.
- [ ] Commit final docs and verification.

### Acceptance Criteria

- [ ] SDK users can define interval-based periodic tasks without Redis-specific types.
- [ ] Scheduler pods can register definitions and dispatch due task instances.
- [ ] Multiple scheduler pods do not enqueue the same due occurrence when Redis leasing succeeds.
- [ ] Scheduler dispatch uses the existing producer path so task state/result behavior stays consistent.
- [ ] Redis-backed periodic definitions survive process restarts.
- [ ] Public root package exports the scheduler API.
- [ ] `go test ./...` and `go vet ./...` pass.
- [ ] Phase 6 creates at least 100 commits from the pre-phase baseline.
