---
title: Scheduler
---

Scheduler turns periodic definitions into dispatched tasks.

## Create scheduler

```go
sched, err := app.NewScheduler(
    goqueue.WithSchedulerPollInterval(time.Second),
    goqueue.WithSchedulerBatchSize(50),
)
```

## Register periodic task

```go
taskDef := goqueue.PeriodicTask{
    Name:     "email-cleanup-hourly",
    TaskName: "cleanup_inbox",
    Queue:    "default",
    Schedule: goqueue.Every(time.Hour),
    Priority: goqueue.DefaultPriority,
    RetryPolicy: goqueue.DefaultRetryPolicy(),
}

if err := sched.RegisterPeriodicTask(context.Background(), taskDef); err != nil {
    // handle
}
```

## Runtime methods

- `PollOnce(ctx)` claims due definitions and dispatches immediate results.
- `Start(ctx)` loops every `PollInterval` until context close.
- `DeletePeriodicTask(ctx, name)` removes a definition.

## Defaults

- interval: `1s`
- batch size: `100`
- lock ttl: `30s`
- scheduler codec: `JSONPayloadCodec`

Periodic metadata keys:

- `goqueue.PeriodicMetadataNameKey`
- `goqueue.PeriodicMetadataDueAtKey`
