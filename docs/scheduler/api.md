# Scheduler API

## Types

`Scheduler` is the runtime that registers periodic definitions and dispatches
due task instances.

`PeriodicTask` is the public definition type:

```go
type PeriodicTask struct {
	Name        goqueue.PeriodicTaskName
	TaskName    goqueue.TaskName
	Queue       goqueue.QueueName
	Args        []any
	Kwargs      map[string]any
	Metadata    map[string]string
	Schedule    goqueue.IntervalSchedule
	StartAt     time.Time
	Priority    goqueue.Priority
	RetryPolicy goqueue.RetryPolicy
}
```

`Every(duration)` creates an interval schedule. The current scheduler supports
fixed intervals; cron-style schedules are reserved for a future release.

## Methods

`RegisterPeriodicTask(ctx, definition)` stores or updates a periodic definition.

`DeletePeriodicTask(ctx, name)` removes a definition, its due index entry, and
any active lease.

`PollOnce(ctx)` runs one due scan and dispatch pass. It is useful for tests,
controlled jobs, and custom process supervision.

`Start(ctx)` runs `PollOnce` immediately and then repeats on the configured poll
interval until the context is canceled.

## Options

| Option | Purpose |
| --- | --- |
| `WithSchedulerIdentity` | Stable identity stored in due-task lease requests. |
| `WithSchedulerDefaultQueue` | Queue used when a definition omits `Queue`. |
| `WithSchedulerPollInterval` | Delay between `Start` loop polls. |
| `WithSchedulerBatchSize` | Maximum due definitions claimed in one poll. |
| `WithSchedulerLockTTL` | Redis lease TTL for claimed due definitions. |
| `WithSchedulerCodec` | Payload codec used by the dispatch producer. |
| `WithSchedulerNow` | Deterministic clock for tests. |
