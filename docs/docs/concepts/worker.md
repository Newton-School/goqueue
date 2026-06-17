---
title: Worker
---

Workers execute tasks from a Redis consumer group.

## Register handlers

```go
app.RegisterTask("send_email", task.TaskHandlerFunc(func(ctx task.HandlerContext, p task.TaskPayload) (task.TaskResult, error) {
    // do work
    return goqueue.SucceededResult("ok"), nil
}))
```

## Create and run worker

```go
worker, err := app.NewWorker(
    goqueue.WithWorkerQueue("emails"),
    goqueue.WithWorkerGroup("workers"),
    goqueue.WithWorkerConsumer("worker-a"),
)
if err != nil { ... }

err = worker.Start(context.Background())
```

## Worker defaults

- group: `goqueue`
- consumer: `worker`
- codec: `JSONPayloadCodec`
- concurrency: `1`
- read batch: `1`
- block: `250ms`
- move due enabled: `true`
- move due limit: `100`
- idle delay: `50ms`
- dead-letter enabled: `true`
- pending recovery enabled: `false`

## Worker flow

1. read message -> `RECEIVED`
2. decode + validate
3. if expired -> `EXPIRED` and dead-letter optionally
4. execute handler
5. if failed and retry allowed -> schedule retry with backoff
6. write final state and result
7. for workflow tasks, advance chain/group/chord progress
8. acknowledge consumed stream entry

## Worker options

- `WithWorkerQueue`, `WithWorkerGroup`, `WithWorkerConsumer`
- `WithWorkerCodec`
- `WithWorkerConcurrency`
- `WithWorkerReadBatch`
- `WithWorkerBlock`
- `WithWorkerMoveDueEnabled`, `WithWorkerMoveDueLimit`
- `WithWorkerIdleDelay`
- `WithWorkerDeadLetterEnabled`
- `WithWorkerPendingRecoveryEnabled`
- `WithWorkerPendingMinIdle`, `WithWorkerPendingClaimBatch`, `WithWorkerPendingClaimInterval`
