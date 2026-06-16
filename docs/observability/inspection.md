# Task Inspection API

The `inspect` package offers read-only visibility into task lifecycle state stored by the
Redis backend:

- task state (`TaskState`)
- task result (`TaskResult`)
- combined state/result snapshot (`TaskSnapshot`)
- dead-letter stream for a queue
- queue counts (`ready`, `scheduled`, `dead`)
- connectivity health check (`Ping`)

Example:

```go
app, err := goqueue.New(goqueue.WithRedisURL("redis://localhost:6379/0"))
if err != nil {
  panic(err)
}

inspector, err := app.NewInspector()
if err != nil {
  panic(err)
}

snapshot, err := inspector.TaskSnapshot(context.Background(), "123e4567-e89b-42d3-a456-556642440111")
if err != nil {
  panic(err)
}

fmt.Printf("task=%s state=%s result=%s\\n", snapshot.TaskID, snapshot.State.State, snapshot.Result.State)
```

Notes:

- Inspection APIs are intentionally read-only.
- State and result retrieval are separate backend reads.
- Snapshot requests aggregate best-effort: if one side is unavailable it still returns the
  available side and surfaces request error only when both miss.
