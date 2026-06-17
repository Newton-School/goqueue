# goqueue

goqueue is a Redis-backed Go SDK for background task execution.

## Prerequisites

- Go `1.26.4` or newer
- Redis reachable via `redis://` or `rediss://`
- Optional local Redis for integration test commands

## Installation

```bash
go get github.com/Newton-School/goqueue
```

## Setup

Applications must create an app and pass explicit options.

```go
package main

import "github.com/Newton-School/goqueue"

func main() {
	app, err := goqueue.New(
		goqueue.WithRedisURL("redis://localhost:6379/0"),
		goqueue.WithDefaultQueue("default"),
		goqueue.WithNamespace("goqueue"),
	)
	if err != nil {
		return
	}
	_ = app
}
```

Use your own error handling and constructor lifecycle around app initialization.

## Environment Variables

`goqueue` does not load `.env` files internally.

| Variable | Required | Default | Purpose |
| --- | --- | --- | --- |
| `GOQUEUE_REDIS_URL` | Yes for app startup and integration tests | _empty_ | Redis connection URL passed to `WithRedisURL`. |
| `GOQUEUE_NAMESPACE` | No | `goqueue` | Redis namespace for all SDK keys. |
| `GOQUEUE_RUN_INTEGRATION_TESTS` | No | `false` | Enables Redis-backed integration tests. |

Redis URLs must use `redis://` or `rediss://`.

## Local integration test setup

```bash
GOQUEUE_RUN_INTEGRATION_TESTS=true GOQUEUE_REDIS_URL=redis://localhost:6379/0 make integration-test
```

## Setup validation

Run this locally after setup changes:

```bash
make verify
```

This repository uses `make audit` in CI.

## Documentation

Documentation source is in `docs/`. Install and run it locally with:

```bash
make docs-install
make docs-start
```

Build static documentation site:

```bash
make docs-build
```

Alternatively, run npm directly from the docs folder:

```bash
cd docs
npm install
npm run docs-start
npm run docs-build
```
