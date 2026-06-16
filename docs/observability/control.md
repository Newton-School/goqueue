# Control API

The `admin` package provides explicit, opt-in production control operations for
task recovery and queue hygiene.

## Available Operations

- Retry tasks from stored payloads with explicit scheduling controls.
- Revoke tasks by state.
- Replay dead-letter entries into a queue.
- Delete dead-letter entries in bulk.
- Purge queue containers and optional message/state/result payloads.

```go
app, err := goqueue.New(
	goqueue.WithRedisURL("redis://localhost:6379/0"),
)
if err != nil {
	panic(err)
}

adminClient, err := app.NewAdmin()
if err != nil {
	panic(err)
}

result, err := adminClient.RetryTask(context.Background(), taskID, goqueue.RetryTaskOptions{
	Queue:    "critical",
	ClearState: true,
})
if err != nil {
	panic(err)
}

fmt.Printf("retried=%s scheduled=%v\n", result.TaskID, result.EnqueueResult.Scheduled)
```

## Safety Notes

- All operations mutate queue state and require explicit command/API calls.
- `PurgeQueue` has opt-in data deletion flags. Omit `delete-messages`/`delete-states`/
  `delete-results` to keep durable payloads while clearing stream entries.
- Replayed dead-letter tasks always reset attempt counters to `0`.

## CLI

Use `goqueue control ...` for operational workflows from shell environments:

- `goqueue control retry-task --id <task-id> [--queue <queue>]`
- `goqueue control revoke-task --id <task-id> [--reason <text>]`
- `goqueue control replay-dead-letter --queue <queue> --stream-id <id>`
- `goqueue control delete-dead-letter --queue <queue> --stream-id <id>[,<id>...]`
- `goqueue control purge-queue --queue <queue> [--delete-messages] [--delete-states] [--delete-results]`
