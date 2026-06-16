# Phase 8 Observability, Inspection, and CLI Commands

## Goal

Add a first-class inspection API for production operations (task state/result lookup,
dead-letter inspection, queue metrics, and connectivity checks) and a lightweight
command-line interface for read-only operations.

## Phase 8 Scope

- Add a dedicated inspection package with explicit error handling and ergonomic
  query methods.
- Export an inspection entrypoint from the root `goqueue` API.
- Add a CLI under `cmd/goqueue` for common operations.
- Add structured task and queue inspection output formats (text + JSON).
- Add tests for all new inspection and CLI behavior surfaces.
- Update documentation and roadmap status.

## Work Items

- [ ] 1. Add inspection package skeleton and interfaces.
- [ ] 2. Implement task state lookup API.
- [ ] 3. Implement task result lookup API.
- [ ] 4. Implement task result cleanup API.
- [ ] 5. Implement dead-letter read API.
- [ ] 6. Implement queue stats API.
- [ ] 7. Implement backend ping/health check wrapper.
- [ ] 8. Add typed inspection output models.
- [ ] 9. Add root facade export for inspector creation.
- [ ] 10. Add app-level `NewInspector` helper.
- [ ] 11. Add CLI command bootstrap and config parsing.
- [ ] 12. Add `inspect task state` command.
- [ ] 13. Add `inspect task result` command.
- [ ] 14. Add `inspect task forget-result` command.
- [ ] 15. Add `inspect dead-letters` command.
- [ ] 16. Add `inspect stats` command.
- [ ] 17. Add `inspect ping`/health command.
- [ ] 18. Add CLI support for env and flag fallback configuration.
- [ ] 19. Add JSON output mode and stable text mode.
- [ ] 20. Add unit tests for inspection package request/response validation.
- [ ] 21. Add command parsing and output tests for all commands.
- [ ] 22. Add docs for CLI usage.
- [ ] 23. Add roadmap update and production hardening notes.
- [ ] 24. Run formatting, vet, and tests.

### Commit Plan Tracking

- [ ] Split implementation into production-safe incremental commits.
- [ ] Keep exports minimal and dependency boundaries explicit.
- [ ] No command should mutate task payloads or execution state except optional
  result cleanup when explicitly requested.
