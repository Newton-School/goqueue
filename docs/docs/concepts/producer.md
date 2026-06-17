---
title: Producer
---

The producer publishes validated tasks to Redis through ready/scheduled queues.

## Create producer

```go
producer, err := app.NewProducer()
```

Defaults:

- queue: app default queue
- codec: `JSONPayloadCodec`
- clock: `time.Now().UTC`

## ApplyAsync

```go
res, err := producer.ApplyAsync(
    context.Background(),
    "task_name",
    []any{"arg1"},
    map[string]any{"x": "value"},
    goqueue.WithApplyQueue("emails"),
)
```

`ApplyAsync`:

- validates task name and payload
- sets initial state (`PENDING` or `SCHEDULED`)
- stores task message and returns `AsyncResult`

## Apply options

- `WithApplyQueue`
- `WithApplyTaskID`
- `WithApplyMetadata`
- `WithApplyPriority`
- `WithApplyRetryPolicy`
- `WithApplyCountDown`
- `WithApplyETA`
- `WithApplyExpiresAt`
- `WithApplyAttempt`
- `WithApplyCreatedAt`

## Async result

Returned by `ApplyAsync` and `NewAsyncResult`:

- `ID()` task id
- `TaskState(ctx)` fetch latest lifecycle state
- `TaskResult(ctx)` fetch latest execution result
- `ForgetTaskResult(ctx)` clears stored result
