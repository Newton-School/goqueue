# goqueue

goqueue is a Redis-backed task queue SDK for Go services. The target is a
Celery-style developer experience for Go: named tasks, queue routing, immediate
execution, scheduled execution, retries, workers, periodic jobs, and workflow
primitives such as groups and chains.

This repository is in Phase 0. The current public surface establishes the
module, configuration model, validation rules, CI, and documentation baseline.
Producer and worker execution APIs will be built on top of this foundation.

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

	_ = app
}
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

## Roadmap

1. Core task envelope and registry.
2. Redis Streams backend for ready queues.
3. Producer API for immediate and delayed tasks.
4. Worker runtime with acknowledgements and graceful shutdown.
5. Retries, dead-letter queues, and task expiration.
6. Scheduler and periodic jobs.
7. Canvas primitives: chains, groups, and chords.
8. Observability, inspection APIs, and CLI commands.

## Security

Do not place credentials in task payloads, queue names, logs, or test fixtures.
`Config.RedactedRedisURL` is available for safe connection logging. Report
security issues privately through GitHub Security Advisories for this
repository.
