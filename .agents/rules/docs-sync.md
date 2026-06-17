# Documentation Sync Rules

Use this file for any change to SDK behavior, APIs, or operator-facing behavior.

## Change Coverage

- Any modification to implementation behavior in `app.go`, `task/`, `producer/`,
  `worker/`, `scheduler/`, `workflow/`, `backend/`, `redisbackend/`,
  `inspect/`, `admin/`, `cmd/goqueue/`, or shared option/config surfaces must be
  reflected in `docs/docs`.
- API/contract updates should also update any affected reference pages and
  concept pages for discoverability.
- CLI behavior updates in `cmd/goqueue/` must update `docs/reference/cli.md`.
- Setup/configuration changes to environment variables, bootstrap commands, or
  installation prerequisites require synchronized updates in `README.md` per
  `readme-policy.md` and, when helpful, `docs/getting-started/*`.

## Required Artifacts

- Keep docs pages focused and short.
- Keep file names aligned with the owning flow:
  - `task-model`, `producer`, `worker`, `scheduler`, `workflow`,
    `configuration`, `redis-backend`.
- For operational behavior changes, update `concepts/inspect-and-admin.md`,
  `reference/cli.md`, and/or `reference/errors.md` depending on scope.
- Do not move behavior detail into `README.md`; keep setup-only constraints in
  `README.md`.

## Update Workflow

When you change SDK files:

1. Identify the owning rule file for the behavior change (for example,
   `producer-flow.md` for producer timing).
2. Identify affected docs pages in `docs/docs`.
3. Update rules in `.agents/rules/*` for behavior semantics where needed.
4. Update docs in the same commit when behavior visibility changed for users.
5. Run docs build before finalizing: `cd docs && npm run docs-build`.

## Verification

- For documentation-only changes, run docs lint/build and confirm no broken links
  in the modified pages.
- For behavior changes, combine docs updates with focused code tests and
  `make verify` (plus `GOQUEUE_RUN_INTEGRATION_TESTS=true make integration-test`
  where Redis contracts change).
