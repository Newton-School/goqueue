---
title: Configuration
---

The public entry point is `goqueue.New`.

## Required and optional settings

```go
app, err := goqueue.New(
    goqueue.WithRedisURL("redis://localhost:6379/0"), // required
    goqueue.WithDefaultQueue("default"),               // optional, defaults to "default"
    goqueue.WithNamespace("goqueue"),                   // optional, defaults to "goqueue"
)
```

## Config behavior

- `WithRedisURL` must be `redis://` or `rediss://` and include host.
- Queue names and namespaces must be safe token strings.
- Empty config is validated during `New` and returns an error.

## App-level helpers

Once created, `App` exposes:

- `RegisterTask(name, handler)`
- `LookupTask(name)`, `TaskNames()`
- `NewProducer`, `NewWorker`, `NewScheduler`, `NewCanvas`
- `NewInspector`, `NewAdmin`
- `NewRedisBackend`

## Redis URL safety

Use `Config.RedactedRedisURL()` when logging to avoid exposing credentials.
