---
title: goqueue
sidebar_position: 1
slug: /
---

goqueue is a Redis-backed Go SDK for queueing and executing background work.

It has one app entry point and clear runtime roles:

1. **Producers** create tasks and push them to Redis.
2. **Workers** consume registered task handlers and persist results.
3. **Schedulers** dispatch periodic tasks.
4. **Canvas** dispatches chain/group/chord workflows.
5. **Inspect/Admin** provide visibility and operational control.

## Core setup pattern

```go
app, err := goqueue.New(
    goqueue.WithRedisURL("redis://localhost:6379/0"),
    goqueue.WithDefaultQueue("default"),
    goqueue.WithNamespace("goqueue"),
)
if err != nil {
    log.Fatal(err)
}
```

From this app you can create producer, worker, scheduler, canvas, inspector, and admin clients.

## What this documentation covers

- Application setup and validation
- Producer and worker behavior
- Retries, deadlines, and dead-letter behavior
- Scheduler and periodic task registration
- Workflow primitives (chain/group/chord)
- Inspect/admin operations and CLI usage
