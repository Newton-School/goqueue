# goqueue

goqueue is a Redis-backed task queue SDK for Go services. The target is a
Celery-style developer experience for Go: named tasks, queue routing, immediate
execution, scheduled execution, retries, workers, periodic jobs, and workflow
primitives such as groups and chains.

This repository has completed Phase 7 canvas and workflow primitives. The public surface
now includes task identity primitives, producer APIs, Redis backend storage, and
a production-grade worker runtime with acknowledgements, retries, dead-letter
queues, pending recovery, task state/result persistence, and Redis-coordinated
periodic task dispatch, chains, groups, and chords.
Phase 8 added inspection APIs for queue health and task observability. Phase 9 adds
control-plane operations for retries, revocation, dead-letter replay and cleanup,
and queue purge through the root `Admin` API and dedicated CLI commands.

## Installation

Requires Go 1.26.4 or newer.

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

## Scheduler Runtime

```go
package main

import (
	"context"
	"log"
	"time"

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

	scheduler, err := app.NewScheduler(
		goqueue.WithSchedulerIdentity("scheduler-pod-1"),
		goqueue.WithSchedulerPollInterval(time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = scheduler.RegisterPeriodicTask(context.Background(), goqueue.PeriodicTask{
		Name:     "welcome-email",
		TaskName: "email.send",
		Schedule: goqueue.Every(10 * time.Minute),
		Args:     []any{"u_123"},
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := scheduler.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
```

## Canvas Workflows

```go
canvas, err := app.NewCanvas()
if err != nil {
	log.Fatal(err)
}

chainResult, err := canvas.ApplyChain(context.Background(), goqueue.Chain{
	Signatures: []goqueue.Signature{
		{Name: "email.prepare", Args: []any{"u_123"}},
		{Name: "email.send", Args: []any{"u_123"}},
	},
})
if err != nil {
	log.Fatal(err)
}

chordResult, err := canvas.ApplyChord(context.Background(), goqueue.Chord{
	Header: goqueue.Group{
		Signatures: []goqueue.Signature{
			{Name: "email.send", Args: []any{"u_1"}},
			{Name: "email.send", Args: []any{"u_2"}},
		},
	},
	Callback: goqueue.Signature{Name: "email.report"},
})
if err != nil {
	log.Fatal(err)
}

_, _ = chainResult, chordResult
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
- Periodic task definitions backed by Redis hashes and due-time sorted sets.
- Short-lived Redis leases for multi-pod scheduler coordination.
- Chain and group workflow state backed by Redis hashes, sets, and Lua scripts.
- Task state and task result storage with TTL-ready APIs.
- Queue stats for ready, scheduled, and dead-letter counts.

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

## CLI

The CLI is available under `cmd/goqueue` for inspection and operational control:

```bash
go run ./cmd/goqueue inspect task --id <task-id>
go run ./cmd/goqueue inspect stats --queue default
go run ./cmd/goqueue inspect deadletters --queue default --count 20
go run ./cmd/goqueue inspect ping
go run ./cmd/goqueue control retry-task --id <task-id> --queue critical --json
go run ./cmd/goqueue control revoke-task --id <task-id> --reason "operator request"
go run ./cmd/goqueue control replay-dead-letter --queue default --stream-id 1-0
go run ./cmd/goqueue control delete-dead-letter --queue default --stream-id 1-0,1-1
go run ./cmd/goqueue control purge-queue --queue default --yes --delete-messages
```

Set `--json` when you need machine-readable output.

## Development Checks

```bash
make verify
make audit
GOQUEUE_RUN_INTEGRATION_TESTS=true GOQUEUE_REDIS_URL=redis://localhost:6379/0 make integration-test
```

`make audit` runs formatting, `go vet`, unit tests, `staticcheck`,
`govulncheck`, and the race detector. Redis integration tests are separate
because they require a running Redis server.

## Package Layout

```text
.
├── app.go, config.go, errors.go
│   Public SDK facade, app construction, and top-level configuration.
├── task/
│   Redis-independent task domain model: identifiers, payloads, envelopes,
│   messages, handlers, results, retry policy, timing, and registry.
├── backend/
│   Backend interfaces and storage request/response contracts used by
│   producers, schedulers, workers, inspectors, and admin controls.
├── admin/
│   Control APIs for retries, revocations, dead-letter replay/cleanup, and queue
│   purge.
├── producer/
│   Producer API for enqueuing immediate and scheduled tasks.
├── inspect/
│   Read-only APIs for task state, result, dead-letter queues, and queue stats.
├── worker/
│   Worker runtime for consuming and executing task messages.
├── scheduler/
│   Periodic task definitions, scheduler runtime, and dispatch coordination.
├── workflow/
│   Canvas primitives: signatures, chains, groups, chords, and dispatch APIs.
├── redisbackend/
│   Redis Streams, sorted sets, Lua scripts, scheduler leases, workflow state, task state, and result storage.
├── docs/reliability/
│   Operational notes for DLQ, recovery, and failure metadata.
├── docs/scheduler/
│   Operational notes for periodic jobs and Redis scheduler coordination.
├── docs/workflows/
│   Usage and Redis state notes for canvas workflows.
└── .github/workflows/
    CI verification.
```

The root package keeps the convenient `goqueue.X` API. Focused subpackages own
the implementation boundaries used by Redis storage, workers, schedulers,
workflow primitives, inspection, control operations, and the CLI.

## Roadmap

1. ✅ Producer API for immediate and delayed tasks.
2. ✅ Worker runtime with acknowledgements and graceful shutdown.
3. ✅ Retries, dead-letter queues, and task expiration.
4. ✅ Scheduler and periodic jobs.
5. ✅ Canvas primitives: chains, groups, and chords.
6. ✅ Observability, inspection APIs, and CLI commands.
7. ✅ Operational control-plane APIs and command surface.

## Security

Do not place credentials in task payloads, queue names, logs, or test fixtures.
`Config.RedactedRedisURL` is available for safe connection logging. Report
security issues privately through GitHub Security Advisories for this
repository.
