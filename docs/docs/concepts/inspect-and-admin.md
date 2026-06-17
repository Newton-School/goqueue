---
title: Inspect and Admin
---

goqueue provides read-only inspection and explicit operator controls.

## Inspect client

```go
inspector, err := app.NewInspector()
```

Read-only methods:

- `Ping(ctx)`
- `QueueStats(ctx, queue)`
- `ReadDeadLetters(ctx, queue, count)`
- `TaskState(ctx, taskID)`
- `TaskResult(ctx, taskID)`
- `TaskSnapshot(ctx, taskID)`
- `ForgetTaskResult(ctx, taskID)`

All inspect methods validate IDs and queue names before read.

## Admin client

```go
admin, err := app.NewAdmin()
```

Admin controls:

- `RetryTask(ctx, taskID, RetryTaskOptions{})`
- `RevokeTask(ctx, taskID, reason)`
- `ReplayDeadLetter(ctx, queue, streamID, options)`
- `DeleteDeadLetters(ctx, queue, ids...)`
- `PurgeQueue(ctx, PurgeQueueOptions{})`

Use admin operations carefully; they mutate queue state.
