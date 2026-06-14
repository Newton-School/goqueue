# Phase 2 Redis Backend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the Redis-backed queue storage layer that later producer and worker phases will use.

**Architecture:** Keep the public root package as the SDK facade. Add a focused `backend` package for storage contracts and a public `redisbackend` package for Redis Streams, sorted sets, message storage, task state, and task result persistence. Redis integration tests are gated by `GOQUEUE_RUN_INTEGRATION_TESTS=true` and `GOQUEUE_REDIS_URL`.

**Tech Stack:** Go 1.26, `github.com/redis/go-redis/v9`, Redis Streams, Redis sorted sets, Redis strings/hashes, Lua scripts for atomic enqueue and due-scheduled movement.

---

### Commit Plan

- [x] Add the Phase 2 plan document.
- [x] Add go-redis dependency.
- [x] Add backend package errors.
- [x] Add backend queue contract.
- [x] Add backend enqueue/read/ack types.
- [x] Add backend state/result types.
- [x] Add backend stats types.
- [x] Add Redis backend options.
- [x] Add Redis backend constructor.
- [x] Add Redis namespace validation.
- [x] Add Redis key builder.
- [x] Add Redis key tests for ready streams.
- [x] Add Redis key tests for scheduled sets.
- [x] Add Redis key tests for task storage.
- [x] Add Redis message codec.
- [x] Add Redis message codec validation.
- [x] Add Redis message storage abstraction.
- [x] Add ready enqueue Lua script.
- [x] Add ready enqueue command.
- [x] Add scheduled enqueue Lua script.
- [x] Add scheduled enqueue command.
- [x] Add due scheduled Lua script.
- [x] Add due scheduled parser.
- [x] Add due scheduled command.
- [x] Add consumer group creation.
- [x] Add ready stream read options.
- [x] Add ready stream read command.
- [x] Add ready stream ack command.
- [x] Add state storage model.
- [x] Add state set command.
- [x] Add state get command.
- [x] Add result storage model.
- [x] Add result save command.
- [x] Add result get command.
- [x] Add result forget command.
- [x] Add queue stats command.
- [x] Add backend health ping.
- [x] Add close method.
- [x] Add integration test gating.
- [x] Add integration test Redis cleanup namespace.
- [x] Add integration test ready enqueue/read/ack.
- [x] Add integration test scheduled enqueue/move.
- [x] Add integration test state storage.
- [x] Add integration test result storage.
- [x] Add root config Redis backend helper.
- [x] Add root facade docs for Redis backend.
- [x] Update README package layout.
- [x] Update `.env.example` comments.
- [ ] Run full verification.
- [ ] Run race tests.

### Acceptance Criteria

- `make verify` passes.
- `go test -race ./...` passes.
- At least 50 commits are created for Phase 2.
- Redis integration tests skip by default and run only when explicitly enabled.
- Ready queues use Redis Streams.
- Scheduled queues use Redis sorted sets.
- Task state and results are stored with TTL-ready APIs.
- Lua scripts perform atomic ready enqueue, scheduled enqueue, and due scheduled movement.
- Root package remains small and stable; Redis implementation lives outside the root package.
