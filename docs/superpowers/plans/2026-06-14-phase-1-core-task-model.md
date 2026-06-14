# Phase 1 Core Task Model Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the production-grade, Redis-independent task model that future producer, backend, scheduler, and worker phases will share.

**Architecture:** Keep Phase 1 in the root `goqueue` package so the public SDK API is direct and stable. Split the model into small files by responsibility: identifiers, timing, retry policy, payloads, envelopes, handler contracts, registry, and app registration. Every behavior-changing slice gets tests before implementation and a small commit.

**Tech Stack:** Go 1.26 standard library only in this phase; no Redis client dependency until Phase 2.

---

### Commit Plan

- [ ] Add the Phase 1 plan document.
- [ ] Add task name validation.
- [ ] Add queue name validation.
- [ ] Add task ID validation and generation.
- [ ] Add priority validation.
- [ ] Add task state constants and terminal checks.
- [ ] Add retry policy defaults and validation.
- [ ] Add task timing options for ETA, countdown, and expiration.
- [ ] Add task payload copying.
- [ ] Add JSON payload codec.
- [ ] Add task envelope defaults.
- [ ] Add task envelope validation.
- [ ] Add task envelope copy behavior.
- [ ] Add task message conversion.
- [ ] Add handler context.
- [ ] Add task result contract.
- [ ] Add task handler adapter.
- [ ] Add registry registration.
- [ ] Add registry lookup.
- [ ] Add app-level task registration.
- [ ] Update public docs for Phase 1.

### Acceptance Criteria

- `make verify` passes.
- Phase 1 has at least 20 commits after the Phase 0 commits.
- The SDK exposes task model primitives without requiring Redis.
- Public errors support `errors.Is`.
- Payload and metadata APIs copy mutable data before storing or returning it.
- README clearly distinguishes available Phase 1 APIs from planned Redis execution features.
