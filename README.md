# goqueue

goqueue is a Redis-backed task queue SDK for Go services. The target is a
Celery-style developer experience for Go: named tasks, queue routing, immediate
execution, scheduled execution, retries, workers, periodic jobs, and workflow
primitives such as groups and chains.

This repository has completed Phase 5 reliability hardening. The public surface
now includes task identity primitives, producer APIs, Redis backend storage, and
a production-grade worker runtime with acknowledgements, retries, dead-letter
queues, pending recovery, and task state/result persistence.

## Installation

```bash
go get github.com/Newton-School/goqueue
```

## Producer Usage

```go
package main

import (
	"context"
	"log"

	"github.com/Newton-School/goqueue"
)

func main() {
	app, err := goqueue.New(
		goqueue.WithRedisURL("redis://localhost:6379/0"),
		goqueue.WithDefaultQueue("default"),
		goqueue.WithNamespace("goqueue"),
	)
	if err != nil {
		log.Fatal(err)
	}

	producer, err := app.NewProducer()
	if err != nil {
		log.Fatal(err)
	}

	result, err := producer.ApplyAsync(context.Background(), "email.send_welcome", []any{"u_123"}, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = app.RegisterTask("email.send_welcome", goqueue.TaskHandlerFunc(
		func(ctx goqueue.HandlerContext, payload goqueue.TaskPayload) (goqueue.TaskResult, error) {
			return goqueue.SucceededResult("queued-model-ready"), nil
		},
	))
	_ = result
	if err != nil {
		log.Fatal(err)
	}
}
```

## Task Envelopes

Task envelopes are Redis-independent. They are the shared model future
producers, schedulers, and workers will use.

```go
envelope, err := goqueue.NewTaskEnvelope(goqueue.TaskEnvelopeInput{
	Name:   "email.send_welcome",
	Queue:  "default",
	Args:   []any{"u_123"},
	Kwargs: map[string]any{"template": "welcome"},
})
if err != nil {
	log.Fatal(err)
}

message, err := goqueue.TaskEnvelopeToMessage(envelope, goqueue.JSONPayloadCodec{})
if err != nil {
	log.Fatal(err)
}

_ = message
```

## Worker Runtime

```go
package main

import (
	"context"
	"log"

	"github.com/Newton-School/goqueue"
)

func main() {
	app, err := goqueue.New(
		goqueue.WithRedisURL("redis://localhost:6379/0"),
		goqueue.WithDefaultQueue("default"),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = app.RegisterTask("email.send_welcome", goqueue.TaskHandlerFunc(
		func(ctx goqueue.HandlerContext, payload goqueue.TaskPayload) (goqueue.TaskResult, error) {
			return goqueue.SucceededResult("email sent"), nil
		},
	))
	if err != nil {
		log.Fatal(err)
	}

	worker, err := app.NewWorker(
		goqueue.WithWorkerGroup("workers"),
		goqueue.WithWorkerConsumer("pod-1"),
		goqueue.WithWorkerConcurrency(4),
		goqueue.WithWorkerPendingRecoveryEnabled(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := worker.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
```

## Reliability

Phase 5 workers use strict ack ordering: messages are acknowledged only after
state, result, retry scheduling, or DLQ persistence succeeds.

Unrecoverable tasks are written to Redis-backed dead-letter streams with a
failure reason, source stream ID, worker group, consumer, error, and timestamp.

Pending recovery is opt-in and uses Redis `XAUTOCLAIM` to reclaim messages that
were read by a worker but never acknowledged.

Failure metadata is stored on `TaskResult.Metadata` using public
`FailureMetadata*` keys and `Failure*` category constants.

## Redis Backend

The Redis backend persists task messages and queue state.

```go
backend, err := app.NewRedisBackend()
if err != nil {
	log.Fatal(err)
}
defer backend.Close()
```

Phase 5 and earlier backend capabilities:

- Ready queues backed by Redis Streams.
- Scheduled queues backed by Redis sorted sets.
- Atomic Lua scripts for ready enqueue, scheduled enqueue, and due-task moves.
- Consumer group creation, stream reads, and acknowledgements.
- Dead-letter streams for unrecoverable worker failures.
- Stale pending message recovery through Redis `XAUTOCLAIM`.
- Task state and task result storage with TTL-ready APIs.
- Queue stats for ready and scheduled counts.

## Configuration

The SDK does not read `.env` files from library code. Applications should load
their own configuration and pass it explicitly to `goqueue.New`.

| Variable | Required | Default | Purpose |
| --- | --- | --- | --- |
| `GOQUEUE_REDIS_URL` | Yes for apps and Redis integration tests | Empty | Redis connection URL passed by the application as `WithRedisURL`. |
| `GOQUEUE_RUN_INTEGRATION_TESTS` | No | `false` | Enables Redis-backed integration tests. |

Redis URLs must use `redis://` or `rediss://`. Queue names and namespaces must
use 1-128 characters from `A-Z`, `a-z`, `0-9`, `.`, `_`, `:`, and `-`.

## Development

```bash
make verify
```

`make verify` checks formatting, runs `go vet`, and executes the full test
suite.

## Package Layout

```text
.
├── app.go, config.go, errors.go
│   Public SDK facade, app construction, and top-level configuration.
├── task/
│   Redis-independent task domain model: identifiers, payloads, envelopes,
│   messages, handlers, results, retry policy, timing, and registry.
├── backend/
│   Backend interfaces and storage request/response contracts used by future
│   producers, schedulers, and workers.
├── producer/
│   Producer API for enqueuing immediate and scheduled tasks.
├── worker/
│   Worker runtime for consuming and executing task messages.
├── redisbackend/
│   Redis Streams, sorted sets, Lua scripts, task state, and result storage.
├── docs/superpowers/plans/
│   Phase implementation plans and acceptance checklists.
└── .github/workflows/
    CI verification.
```

The root package keeps the convenient `goqueue.X` API. The `task` package owns
the core task model implementation so future Redis backend, worker, scheduler,
and CLI packages can depend on focused domain packages instead of a crowded
module root.

## Roadmap

1. ✅ Producer API for immediate and delayed tasks.
2. ✅ Worker runtime with acknowledgements and graceful shutdown.
3. Retries, dead-letter queues, and task expiration.
4. Scheduler and periodic jobs.
5. Canvas primitives: chains, groups, and chords.
6. Observability, inspection APIs, and CLI commands.

## Security

Do not place credentials in task payloads, queue names, logs, or test fixtures.
`Config.RedactedRedisURL` is available for safe connection logging. Report
security issues privately through GitHub Security Advisories for this
repository.
