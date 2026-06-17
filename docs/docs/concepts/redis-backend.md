---
title: Redis Backend
---

goqueue ships with Redis as storage backend and queue broker.

## Defaults

- namespace: `goqueue`
- message TTL: `7d`
- state TTL: `24h`
- result TTL: `24h`

## Build options

`app.NewRedisBackend(...)` uses:

- `goqueue.WithRedisURL(...)` from app config
- `goqueue.WithNamespace(...)` from app config
- optional `redisbackend.BackendOption` overrides:
  - `redisbackend.WithMessageTTL`
  - `redisbackend.WithStateTTL`
  - `redisbackend.WithResultTTL`
  - `redisbackend.WithNamespace`

## Extending backend

`app.NewRedisBackend` accepts:

```go
backend, err := app.NewRedisBackend(
    redisbackend.WithMessageTTL(12*time.Hour),
)
```

You can also pass a custom redis client using `redisbackend.WithClient`.

## Notes

- Backend validates Redis URL (`redis://` / `rediss://`).
- Keep credentials out of logs and docs output.
