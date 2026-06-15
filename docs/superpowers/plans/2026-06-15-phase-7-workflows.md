# Phase 7 Canvas And Workflow Primitives Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add Celery-style canvas primitives for reusable signatures, chains, groups, and chords on top of the existing producer, worker, Redis backend, and scheduler foundations.

**Architecture:** Keep public workflow composition types in a new `workflow` package. Store backend-neutral workflow records in `backend`, persist chain/group/chord state in `redisbackend`, and advance workflows from `worker` after successful final result persistence but before acknowledgement. Public root exports keep the convenient `goqueue.X` API.

**Tech Stack:** Go 1.26, Redis hashes/sets/Lua scripts, go-redis/v9, standard Go tests.

---

### Task 1: Workflow Domain

**Files:**
- Create: `workflow/errors.go`
- Create: `workflow/signature.go`
- Create: `workflow/chain.go`
- Create: `workflow/group.go`
- Create: `workflow/chord.go`
- Create tests for each domain type.

- [x] Add signature validation, normalization, and defensive copying.
- [x] Add chain validation for ordered task signatures.
- [x] Add group validation for fan-out task signatures.
- [x] Add chord validation for group header plus callback signature.
- [x] Verify with `go test ./workflow`.
- [x] Commit each narrow behavior.

### Task 2: Backend Workflow Contracts

**Files:**
- Create: `backend/workflow.go`
- Create: `backend/workflow_test.go`
- Modify: `backend/backend.go`
- Modify: package test fakes as needed.

- [x] Add workflow signature storage records.
- [x] Add chain save/advance contracts.
- [x] Add group save/progress contracts.
- [x] Add validation for workflow IDs, indexes, terminal states, and callbacks.
- [x] Verify with `go test ./backend`.
- [x] Commit each contract slice.

### Task 3: Redis Workflow Persistence

**Files:**
- Modify: `redisbackend/keys.go`
- Modify: `redisbackend/scripts.go`
- Create: `redisbackend/workflow_codec.go`
- Create: `redisbackend/workflow.go`
- Create workflow Redis tests and integration tests.

- [x] Store chain metadata and signatures in Redis hashes.
- [x] Advance chains atomically and return the next signature at most once.
- [x] Store group metadata, task membership, and optional chord callback.
- [x] Record group task completion atomically and return a callback at most once.
- [x] Verify with `go test ./redisbackend`.
- [x] Commit each Redis behavior.

### Task 4: Canvas Producer API

**Files:**
- Create: `workflow/canvas.go`
- Create: `workflow/options.go`
- Create: `workflow/result.go`
- Create canvas tests.

- [x] Add `Canvas` with default queue, codec, and clock options.
- [x] Add `ApplySignature`, `ApplyChain`, `ApplyGroup`, and `ApplyChord`.
- [x] Persist workflow state before dispatching workflow task instances.
- [x] Add workflow metadata to dispatched task instances.
- [x] Verify with `go test ./workflow`.
- [x] Commit each canvas behavior.

### Task 5: Worker Workflow Advancement

**Files:**
- Modify: `worker/worker.go`
- Create: `worker/workflow.go`
- Modify worker tests and fakes.

- [x] Advance chain workflows after successful task completion.
- [x] Record group child terminal states after task result persistence.
- [x] Dispatch chord callbacks once when a group succeeds.
- [ ] Preserve strict ack ordering by advancing workflows before acknowledgement.
- [ ] Verify with `go test ./worker`.
- [ ] Commit each worker behavior.

### Task 6: Public API And Documentation

**Files:**
- Modify: `app.go`
- Modify: `task_exports.go`
- Modify: `README.md`
- Modify: `doc.go`
- Create: `docs/workflows/canvas.md`
- Create: `docs/workflows/redis-state.md`
- Modify this plan.

- [ ] Add `App.NewCanvas` with app defaults.
- [ ] Re-export workflow types, results, options, and metadata constants.
- [ ] Document signature, chain, group, and chord usage.
- [ ] Document Redis workflow state and idempotency behavior.
- [ ] Verify with `go test ./...` and `go vet ./...`.
- [ ] Phase 7 creates exactly 100 commits from the pre-phase baseline.

### Acceptance Criteria

- [ ] SDK users can create reusable signatures without Redis-specific types.
- [ ] SDK users can dispatch chains, groups, and chords.
- [ ] Workers advance chains only after persisted successful completion.
- [ ] Workers record group progress for terminal child states.
- [ ] Chord callbacks dispatch once when all group members succeed.
- [ ] Redis state prevents duplicate chain advancement and duplicate chord callbacks.
- [ ] Public root package exports workflow APIs.
- [ ] `go test ./...` and `go vet ./...` pass.
- [ ] Phase 7 creates exactly 100 commits from the Phase 6 baseline.
