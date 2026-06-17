# Public API Rules

Use this file when changing the root package, exported SDK names, package docs,
or anything that an importing application can compile against.

## Package Boundary

- The root `goqueue` package is the public facade. Keep the convenient
  `goqueue.X` API there, but place implementation logic in the focused package
  that owns the behavior.
- `task_exports.go` is the re-export layer for public task, producer, worker,
  scheduler, workflow, inspect, and admin names. Update it when adding a public
  type that should be available from the root package.
- `app.go` wires public constructors to the Redis-backed implementations. Do
  not add business logic there when it belongs in `producer`, `worker`,
  `scheduler`, `workflow`, `inspect`, `admin`, or `redisbackend`.
- Every public package should have a short `doc.go` package comment that
  explains the package responsibility.

## Compatibility

- Preserve existing exported names, method signatures, option names, error
  variables, and constants unless the user explicitly requests a breaking
  change.
- Prefer adding option functions over changing constructor signatures.
- Preserve root re-exports for types that are already public through
  `goqueue.X`.
- When adding a public API, add focused tests in the owning subpackage and a
  root-level test when the root facade or re-export behavior changes.

## Configuration

- Library code must not load `.env` files. Applications pass configuration into
  `goqueue.New`.
- `Config.Validate` must reject missing Redis URLs, invalid Redis URL schemes,
  invalid queue names, and invalid namespaces before any backend is created.
- `Config.RedactedRedisURL` is the only safe Redis URL value for logs or
  operator-facing output.

## Documentation

- Update `README.md` only for setup or configuration setup changes.
- Update relevant rule files in `.agents/rules/` when changing runtime
- semantics, operational behavior, or failure handling.
- Run `docs-sync.md` when exported API behavior is user-visible to keep
  `docs/` pages aligned.
- Do not add internal planning notes, local paths, or private rollout history to
  public docs.
