# Security And Configuration Rules

Use this file when changing configuration, environment variables, Redis URLs,
logging, payload handling, CLI output, `.gitignore`, `.env`, `.env.example`, or
operator-facing errors.

## Secrets

- Never commit real credentials, tokens, private Redis URLs, local paths,
  customer data, or machine-specific values.
- Do not store secrets in task payloads, task metadata, task results, queue
  names, state errors, CLI fixtures, or docs.
- Do not log raw task payloads by default. Payloads may contain user data.
- Use `Config.RedactedRedisURL` or equivalent redaction for diagnostics.

## Environment Variables

- Library code must not load `.env` files directly.
- Applications and tests pass Redis URLs into `goqueue.New` or backend options.
- Keep `.env` and `.env.example` synchronized whenever environment variables
  change.
- `.env.example` must include every project variable. If a value has a safe
  default, include it; if required with no safe default, leave it empty.
- Required runtime configuration must fail validation when missing.

## Current Variables

- `GOQUEUE_REDIS_URL` is required for apps and Redis integration tests when they
  need a Redis connection.
- `GOQUEUE_NAMESPACE` is optional and controls the SDK/CLI Redis namespace.
- `GOQUEUE_RUN_INTEGRATION_TESTS` enables Redis-backed integration tests.

## Validation

- Redis URLs must use `redis://` or `rediss://` and include a host.
- Queue names and namespaces must remain bounded and restricted to safe
  characters.
- TTLs, poll intervals, lock TTLs, batch sizes, retry attempts, countdowns, and
  concurrency values must be validated before use.
- Avoid unbounded loops, unbounded Redis reads, and unbounded task dispatch.

## Ignore Files

- Keep `.gitignore` grouped by use case with comment headers.
- Ignore local env files, build artifacts, caches, generated coverage files,
  editor metadata, and temporary outputs.
- Keep `.env.example` explicitly unignored.

## Security Review Triggers

- Unsafe deserialization, custom codecs, admin controls, purge/delete paths,
  credential handling, Redis key construction, Lua scripts, and CLI output all
  require explicit security review and tests.
- Report vulnerabilities privately as described in `SECURITY.md`; do not add
  public issue-style vulnerability details to docs.
