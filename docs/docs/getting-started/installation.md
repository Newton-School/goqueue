---
title: Installation
sidebar_position: 1
---

## Prerequisites

- Go 1.26.4 or newer
- Redis 5+ accessible with `redis://` or `rediss://`
- `GOQUEUE_REDIS_URL` for local CLI use (recommended)

## Install SDK

```bash
go get github.com/Newton-School/goqueue
```

## Install CLI

```bash
go install github.com/Newton-School/goqueue/cmd/goqueue@latest
```

## Optional integration test setup

```bash
GOQUEUE_RUN_INTEGRATION_TESTS=true \
GOQUEUE_REDIS_URL=redis://localhost:6379/0 \
make integration-test
```

## Next

- Continue to [quick start](./quick-start) to create your first task and worker.
