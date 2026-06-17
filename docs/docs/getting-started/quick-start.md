---
title: Quick Start
sidebar_position: 2
---

This gives you a runnable baseline flow: task registration, produce, run worker.

```go
package main

import (
    "context"
    "log"

    "github.com/Newton-School/goqueue"
    "github.com/Newton-School/goqueue/task"
)

func main() {
    app, err := goqueue.New(
        goqueue.WithRedisURL("redis://localhost:6379/0"),
        goqueue.WithDefaultQueue("default"),
    )
    if err != nil {
        log.Fatal(err)
    }

    if err := app.RegisterTask("send_email", task.TaskHandlerFunc(
        func(ctx task.HandlerContext, p task.TaskPayload) (task.TaskResult, error) {
            _ = ctx
            return goqueue.SucceededResult("email sent"), nil
        },
    )); err != nil {
        log.Fatal(err)
    }

    producer, err := app.NewProducer()
    if err != nil {
        log.Fatal(err)
    }

    go func() {
        worker, err := app.NewWorker(
            goqueue.WithWorkerGroup("goqueue"),
            goqueue.WithWorkerConsumer("worker-1"),
        )
        if err != nil {
            log.Fatal(err)
        }

        if err := worker.Start(context.Background()); err != nil {
            log.Fatal(err)
        }
    }()

    if _, err := producer.ApplyAsync(
        context.Background(),
        "send_email",
        []any{"user@example.com"},
        map[string]any{"subject": "Welcome"},
    ); err != nil {
        log.Fatal(err)
    }
}
```

## What happens next

- Producer publishes a task message with generated task ID.
- Worker reads from `default` queue and executes handler.
- Result and state are stored in Redis.
- Inspector/admin clients can read state, results, and control execution.
