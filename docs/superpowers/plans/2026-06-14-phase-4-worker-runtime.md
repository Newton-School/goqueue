# Phase 4 Worker Runtime Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

## Goal

Build a production-capable worker runtime that consumes from queue groups, executes
registered handlers, updates lifecycle state, persists results, and handles
retries/scheduling for failed tasks.

## Architecture

Keep the root `goqueue` package as facade and add a focused `worker` package for
execution. A worker:

- binds to shared queue registry and storage backend,
- creates/uses a Redis consumer group,
- polls ready tasks via stream reads,
- moves due scheduled messages into ready queues,
- decodes payloads, executes registered handlers, and persists terminal state/result,
- re-enqueues failed tasks when retry policy allows.

## Commit Plan

- [x] Add worker errors.
- [x] Add worker configuration and option types.
- [x] Add worker constructor with queue/backend/registry validation.
- [x] Implement startup path and consumer-group bootstrap.
- [x] Implement schedule-promotion (`MoveDueScheduled`) in worker loop.
- [x] Implement poll/read loop with bounded concurrency.
- [x] Implement task state transitions (RECEIVED, STARTED, RETRYING, terminal).
- [x] Persist task results and handle unknown handler cases.
- [x] Implement retry behavior with exponential backoff and attempt tracking.
- [x] Add graceful stop semantics and context handling.
- [x] Add unit tests for constructor, success path, retry/no-retry, and expired tasks.
- [x] Add app-level worker factory method.
- [x] Re-export worker options and types from root.
- [x] Update README with worker usage and roadmap status.

## Acceptance Criteria

- `App.NewWorker` returns a configured worker bound to app registry/config.
- Workers process ready tasks and mark state transitions.
- Expired tasks are skipped and marked as `EXPIRED`.
- Failed tasks with retry allowance are rescheduled with computed delay.
- Failed tasks at max attempts are marked terminal and not requeued.
- `go test ./...` passes.
