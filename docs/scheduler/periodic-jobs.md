# Periodic Jobs

Phase 6 adds a beat-style scheduler for recurring task dispatch.

The scheduler does not execute handlers. It only registers periodic definitions,
leases due definitions from Redis, and enqueues task instances through the same
producer path used by `ApplyAsync`. Workers continue to consume and execute task
messages from ready queues.

## Basic Usage

```go
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
	Metadata: map[string]string{"source": "billing"},
})
if err != nil {
	log.Fatal(err)
}

if err := scheduler.Start(context.Background()); err != nil {
	log.Fatal(err)
}
```

## Runtime Flow

1. `RegisterPeriodicTask` normalizes the definition with scheduler defaults.
2. The Redis backend stores the definition in a hash and indexes `NextDueAt` in
   a sorted set.
3. `Start` calls `PollOnce` immediately, then repeats on the configured poll
   interval.
4. `PollOnce` asks Redis for due definitions using the scheduler identity,
   batch size, and lock TTL.
5. Redis leases each due definition with a short `SET NX` lock.
6. The scheduler dispatches each leased definition with `producer.ApplyAsync`.
7. The scheduler marks the definition dispatched only after enqueue succeeds.
8. The backend verifies the lease token before advancing `NextDueAt`.

## Multi-Pod Coordination

Multiple scheduler pods may run for the same namespace. Redis lease keys ensure
only one pod can claim a due definition at a time. If a scheduler pod dies after
claiming a definition, the lease expires and a later poll can claim it again.

The scheduler identity should be stable for the pod process and safe for logs.
If one is not configured, goqueue generates a random `scheduler-<hex>` identity.

## Dispatch Metadata

Dispatched task instances include trace metadata:

| Key | Value |
| --- | --- |
| `goqueue.periodic.name` | Periodic definition name. |
| `goqueue.periodic.due_at` | Due occurrence timestamp in RFC3339Nano format. |

User-provided metadata is preserved unless it uses the same reserved keys.
