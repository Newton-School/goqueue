# README Policy

Use this file for any edits to `README.md`.

## Scope

- README must contain only setup instructions.
- Allowed content:
  - prerequisites
  - installation
  - environment/secret setup
  - quick start app construction
  - local setup validation commands
  - minimal CI/repository setup notes
- Disallowed content:
  - feature flow details
  - runtime behavior explanations
  - package internals
  - operational playbooks
  - roadmap/phase history
  - performance/reliability deep dives

## Setup Change Contract

- When setup configuration changes, README must be updated in the same commit.
- Setup configuration changes include:
  - required environment variables
  - default values
  - installation prerequisites
  - CLI/setup command changes
  - bootstrapping examples or flags
- If config is changed in `.env.example`, `Makefile`, or app bootstrap code,
  README setup sections must be kept in sync.

## Enforcement

- Keep long-form usage and behavior notes in `.agents/rules/` files.
- Keep CLI usage docs in `.agents/rules/operations.md` unless it is strictly setup-only.
- If a contributor adds non-setup content to README, move it to the relevant
  rule file under `.agents/rules/`.
