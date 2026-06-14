# Phase 2 Redis Backend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the Redis-backed queue storage layer that later producer and worker phases will use.

**Architecture:** Keep the public root package as the SDK facade. Add a focused `backend` package for storage contracts and a public `redisbackend` package for Redis Streams, sorted sets, message storage, task state, and task result persistence. Redis integration tests are gated by `GOQUEUE_RUN_INTEGRATION_TESTS=true` and `GOQUEUE_REDIS_URL`.

**Tech Stack:** Go 1.26, `github.com/redis/go-redis/v9`, Redis Streams, Redis sorted sets, Redis strings/hashes, Lua scripts for atomic enqueue and due-scheduled movement.

---

### Commit Plan

- [ ] Add the Phase 2 plan document.
- [ ] Add go-redis dependency.
- [ ] Add backend package errors.
- [ ] Add backend queue contract.
- [ ] Add backend enqueue/read/ack types.
- [ ] Add backend state/result types.
- [ ] Add backend stats types.
- [ ] Add Redis backend options.
- [ ] Add Redis backend constructor.
- [ ] Add Redis namespace validation.
- [ ] Add Redis key builder.
- [ ] Add Redis key tests for ready streams.
- [ ] Add Redis key tests for scheduled sets.
- [ ] Add Redis key tests for task storage.
- [ ] Add Redis message codec.
- [ ] Add Redis message codec validation.
- [ ] Add Redis message storage abstraction.
- [ ] Add ready enqueue Lua script.
- [ ] Add ready enqueue command.
- [ ] Add scheduled enqueue Lua script.
- [ ] Add scheduled enqueue command.
- [ ] Add due scheduled Lua script.
- [ ] Add due scheduled parser.
- [ ] Add due scheduled command.
- [ ] Add consumer group creation.
- [ ] Add ready stream read options.
- [ ] Add ready stream read command.
- [ ] Add ready stream ack command.
- [ ] Add pending message claim placeholder contract.
- [ ] Add state storage model.
- [ ] Add state set command.
- [ ] Add state get command.
- [ ] Add result storage model.
- [ ] Add result save command.
- [ ] Add result get command.
- [ ] Add result forget command.
- [ ] Add queue stats command.
- [ ] Add backend health ping.
- [ ] Add close method.
- [ ] Add integration test gating.
- [ ] Add integration test Redis cleanup namespace.
- [ ] Add integration test ready enqueue/read/ack.
- [ ] Add integration test scheduled enqueue/move.
- [ ] Add integration test state storage.
- [ ] Add integration test result storage.
- [ ] Add root config Redis backend helper.
- [ ] Add root facade docs for Redis backend.
- [ ] Update README package layout.
- [ ] Update `.env.example` comments.
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
