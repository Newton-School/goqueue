# goqueue

goqueue is a Redis-backed task queue SDK for Go services. The target is a
Celery-style developer experience for Go: named tasks, queue routing, immediate
execution, scheduled execution, retries, workers, periodic jobs, and workflow
primitives such as groups and chains.

This repository has completed Phase 1. The current public surface establishes
the module, configuration model, task identity primitives, payload codecs, task
envelopes, handler contracts, and task registration. Redis producer and worker
execution APIs will be built on top of this foundation.

## Installation

```bash
go get github.com/Newton-School/goqueue
```

## Current Usage

```go
package main

import (
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

	err = app.RegisterTask("email.send_welcome", goqueue.TaskHandlerFunc(
		func(ctx goqueue.HandlerContext, payload goqueue.TaskPayload) (goqueue.TaskResult, error) {
			return goqueue.SucceededResult("queued-model-ready"), nil
		},
	))
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

## Configuration

The SDK does not read `.env` files from library code. Applications should load
their own configuration and pass it explicitly to `goqueue.New`.

| Variable | Required | Default | Purpose |
| --- | --- | --- | --- |
| `GOQUEUE_REDIS_URL` | Yes for apps and integration tests | Empty | Redis connection URL passed by the application as `WithRedisURL`. |
| `GOQUEUE_RUN_INTEGRATION_TESTS` | No | `false` | Enables Redis-backed integration tests once those tests exist. |

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

1. Redis Streams backend for ready queues.
2. Producer API for immediate and delayed tasks.
3. Worker runtime with acknowledgements and graceful shutdown.
4. Retries, dead-letter queues, and task expiration.
5. Scheduler and periodic jobs.
6. Canvas primitives: chains, groups, and chords.
7. Observability, inspection APIs, and CLI commands.

## Security

Do not place credentials in task payloads, queue names, logs, or test fixtures.
`Config.RedactedRedisURL` is available for safe connection logging. Report
security issues privately through GitHub Security Advisories for this
repository.
