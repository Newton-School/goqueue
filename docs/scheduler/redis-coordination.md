# Redis Scheduler Coordination

The scheduler is designed for multi-pod deployments. Every scheduler instance
uses the same Redis namespace and competes for due periodic definitions through
short-lived leases.

## Redis Keys

| Purpose | Shape |
| --- | --- |
| Periodic definitions | `<namespace>:scheduler:periodic:definitions` |
| Due index | `<namespace>:scheduler:periodic:due` |
| Lease key | `<namespace>:scheduler:periodic:<name>:lease` |

Definitions are JSON records stored in a Redis hash. Due times are indexed in a
sorted set using Unix milliseconds as scores. Lease keys store random tokens and
expire after the configured scheduler lock TTL.

## Dispatch Safety

`ListDuePeriodicTasks` claims due definitions with `SET NX`. A second scheduler
that sees the same due definition cannot claim it while the lease key exists.

`MarkPeriodicTaskDispatched` verifies the lease token in Lua before advancing
the next due time and deleting the lease. If the token is missing or different,
the backend returns `backend.ErrPeriodicTaskLeaseLost`.

## Failure Behavior

If a scheduler crashes after claiming a definition but before dispatching it,
the lease expires and another scheduler can claim it later.

If dispatch succeeds but marking fails, the scheduler returns the error. The
definition may be retried after the lease expires, so task handlers should remain
idempotent when they are used as periodic jobs.

## Deployment Notes

Run at least one scheduler process for each Redis namespace that owns periodic
definitions. Running more than one scheduler process is supported for
availability, but all scheduler processes for the namespace should use the same
application configuration and task registration code.
