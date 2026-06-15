# Workflow Testing

Workflow changes should include package-level tests for the affected layer.

## Recommended Coverage

- `workflow`: validate signatures, chains, groups, chords, metadata merges, and
  canvas dispatch behavior.
- `backend`: validate backend-neutral record and request contracts.
- `redisbackend`: validate key generation, codec round trips, parser guards,
  Lua script behavior, and Redis integration lifecycle when Redis is available.
- `worker`: validate workflow advancement after result persistence and before
  acknowledgement.

## Commands

Use focused tests while developing:

```sh
go test ./workflow
go test ./redisbackend
go test ./worker
```

Before release, run the full verification target.
