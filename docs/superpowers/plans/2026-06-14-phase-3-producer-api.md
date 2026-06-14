# Phase 3 Producer API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

## Goal

Build a public producer API that publishes immediate and scheduled tasks using the
existing task model and Redis backend contracts.

## Architecture

Keep the root `goqueue` package as a facade and add a focused `producer`
subpackage for enqueue logic. The producer should:

- build a validated `task.TaskEnvelope`,
- serialize it via a payload codec,
- write initial task state,
- dispatch to ready or scheduled queue paths.

## Commit Plan

- [x] Add producer errors.
- [x] Add producer configuration/options types.
- [x] Add producer constructor from `backend.QueueBackend`.
- [x] Add `ApplyAsync` with task name/args/kwargs API.
- [x] Add default queue handling and fallback to app config.
- [x] Add countdown and ETA scheduling path.
- [x] Add immediate enqueue path.
- [x] Add async result handle for task state/result inspection.
- [x] Add failure-state writeback for enqueue failures.
- [x] Add producer-specific option constructors for queue, id, timing, priority,
  retry, and metadata.
- [x] Re-export producer options and types in `task_exports.go`.
- [x] Add `App.NewProducer` integration.
- [x] Add producer tests with fake backend.
- [x] Add app-level producer creation test.
- [x] Update README and docs with Phase 3 status and usage example.

## Acceptance Criteria

- `goqueue.App.NewProducer` returns a configured `producer.Producer`.
- `ApplyAsync` publishes to ready queue when ETA is absent.
- `ApplyAsync` publishes to scheduled queue when countdown or ETA is present.
- Invalid task names, apply options, and backend errors return clear errors.
- Enqueue failures record task state as `FAILED` when state persistence works.
- Async result can fetch task state and result and can forget result records.
- New tests cover the producer paths in unit tests.
