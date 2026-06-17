# Agent Rules Index

Use these rule files as the product map for goqueue. Load only the files needed
for the package or flow you are changing, plus `security-config.md` whenever the
change touches configuration, Redis connection handling, logging, payloads, or
operator-facing output.

## Path Routing

| Path | Rule files |
| --- | --- |
| `app.go`, `config.go`, `errors.go`, `doc.go`, `task_exports.go` | `public-api.md`, `security-config.md` |
| `task/` | `task-model.md`, `security-config.md` |
| `producer/` | `producer-flow.md`, `task-model.md`, `redis-backend.md` |
| `backend/` | `redis-backend.md`, `public-api.md` |
| `redisbackend/` | `redis-backend.md`, `security-config.md`, `testing-release.md` |
| `worker/` | `worker-runtime.md`, `task-model.md`, `workflow-canvas.md` |
| `scheduler/` | `scheduler-flow.md`, `producer-flow.md`, `redis-backend.md` |
| `workflow/` | `workflow-canvas.md`, `producer-flow.md`, `worker-runtime.md` |
| `inspect/`, `admin/`, `cmd/goqueue/` | `operations.md`, `redis-backend.md`, `security-config.md` |
| `docs/`, `README.md`, `.github/`, `Makefile`, `.env.example`, `.gitignore` | `testing-release.md`, `security-config.md` |

## Core Expectations

- Keep package ownership clear. Root re-exports belong in `task_exports.go`;
  implementation logic belongs in the focused package that owns the behavior.
- Prefer existing backend contracts over direct Redis access outside
  `redisbackend`.
- Add or update focused tests next to the changed behavior.
- Keep `.env` and `.env.example` synchronized when environment variables
  change. `.env.example` documents required variables with empty values when no
  safe default exists.
- Use grouped comments in `.gitignore` and environment files.
- Keep public docs free of internal plans, local machine paths, credentials,
  and private operational context.
