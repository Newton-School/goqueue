# goqueue Agent Guide

goqueue is a public Go SDK for Redis-backed task queues. Treat every change as
library code that other services may import directly.

## Start Here

1. Stay on the current branch. Do not create, switch, or use another branch or
   worktree unless the user explicitly asks for it.
2. Read the rule file for the flow you are changing under `.agents/rules/`.
   Start with `.agents/rules/README.md` when you are unsure which rule applies.
3. Use `rg` or `rg --files` before broader search tools.
4. Keep commits small and intentional when making product changes.
5. Preserve public API compatibility unless the user explicitly asks for a
   breaking change.
6. Do not add credentials, local paths, machine-specific values, or private
   planning notes to the repository.

## Product Flow

The root package is the public facade. It re-exports stable SDK types from
focused subpackages and wires the default Redis-backed producer, worker,
scheduler, workflow canvas, inspector, and admin clients.

The normal execution path is:

1. Applications create `goqueue.App` with explicit configuration.
2. Producers convert task envelopes into encoded task messages.
3. `redisbackend` stores ready, scheduled, periodic, state, result, dead-letter,
   and workflow data in Redis.
4. Workers read queue streams, execute registered task handlers, persist state
   and results, schedule retries, dead-letter unrecoverable failures, and
   advance workflows.
5. Schedulers dispatch due periodic tasks with Redis-backed coordination.
6. Inspect and admin surfaces provide read-only observability and operational
   recovery without bypassing backend contracts.

## Rule Map

- `.agents/rules/public-api.md` for root facade, package boundaries, and export
  compatibility.
- `.agents/rules/task-model.md` for task names, queues, payload codecs, retry
  policies, timing, metadata, and handlers.
- `.agents/rules/producer-flow.md` for immediate and scheduled task publishing.
- `.agents/rules/redis-backend.md` for Redis keys, streams, sorted sets, Lua
  scripts, leases, codecs, and TTL behavior.
- `.agents/rules/worker-runtime.md` for handler execution, ack ordering,
  retries, DLQ, pending recovery, shutdown, and concurrency.
- `.agents/rules/scheduler-flow.md` for periodic jobs and scheduler
  coordination.
- `.agents/rules/workflow-canvas.md` for signatures, chains, groups, chords,
  result aggregation, and worker advancement.
- `.agents/rules/operations.md` for inspect, admin, CLI, observability, and
  control-plane behavior.
- `.agents/rules/testing-release.md` for verification, CI, integration tests,
  public packaging, and release readiness.
- `.agents/rules/security-config.md` for secrets, environment variables,
  validation, Redis URLs, and safe logging.

## Verification

For code changes, run the narrow package tests that cover the behavior first,
then run `make verify`. Run `make audit` before claiming release readiness,
security hardening, or full-product stability. Redis integration tests require
a running Redis instance and `GOQUEUE_RUN_INTEGRATION_TESTS=true`.
