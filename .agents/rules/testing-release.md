# Testing And Release Rules

Use this file when changing tests, CI, Makefile targets, public docs, package
layout, module metadata, examples, release readiness, or any cross-package
behavior.

## Local Verification

- Run focused package tests first for the behavior you changed.
- Run `make verify` before claiming normal code changes are ready.
- Run `make audit` before claiming full-product stability, security readiness,
  release readiness, or public packaging readiness.
- Run Redis integration tests when Redis behavior, backend contracts,
  scheduling, pending recovery, DLQ, workflow state, or CLI Redis flows change:
  `GOQUEUE_RUN_INTEGRATION_TESTS=true GOQUEUE_REDIS_URL=redis://localhost:6379/0 make integration-test`.
- Use `go test -count=1` for integration reruns where cached results would hide
  stateful failures.

## CI Expectations

- CI must run the audit target.
- CI must run Redis integration tests against a Redis service.
- Keep Go version alignment between `go.mod`, CI, and README.
- Do not add CI steps that require private credentials for public repository
  validation.

## Test Placement

- Put unit tests next to the package they cover.
- Use fake backends for ordering, validation, and error injection.
- Use Redis integration tests for streams, sorted sets, Lua scripts, pending
  entries, leases, TTL behavior, and key deletion semantics.
- Add regression tests for every bug fix that changes runtime behavior.

## Public Packaging

- Keep public docs free of internal phase plans, local paths, private
  operational context, and generated scratch artifacts.
- Keep package docs present for public packages.
- Keep examples compileable when possible.
- Use root facade examples for SDK onboarding and focused package examples when
  documenting internal extension points.
- `CLAUDE.md` should remain a symlink to `AGENTS.md` so agent instructions have
  one source of truth.

## Git Hygiene

- Keep commits small when implementing larger tasks.
- Do not rewrite unrelated user changes.
- Do not run destructive git commands unless explicitly requested.
- Commit messages must be imperative, present tense, and start with a capital
  letter without `feat`, `docs`, or AI prefixes.
