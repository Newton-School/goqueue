# Redis Workflow State

Workflow state is stored separately from queue messages. Producers persist
workflow records before dispatching the first workflow task.

## Chain State

| Key | Purpose |
| --- | --- |
| `<namespace>:workflow:chain:<id>:meta` | Total, completed index, dispatched index, timestamps. |
| `<namespace>:workflow:chain:<id>:signatures` | JSON signatures by chain index. |

`AdvanceWorkflowChain` runs a Lua script that:

1. Ignores duplicate completions for a step already marked complete.
2. Marks the completed step.
3. Returns the next signature only when it has not already been dispatched.

This prevents duplicate next-step dispatch when a worker retries the same stream
message after a crash.

## Group And Chord State

| Key | Purpose |
| --- | --- |
| `<namespace>:workflow:group:<id>:meta` | Total, completed count, failed count, callback flag. |
| `<namespace>:workflow:group:<id>:completed` | Set of child task IDs already counted. |
| `<namespace>:workflow:group:<id>:callback` | Optional JSON callback signature for chords. |

`RecordWorkflowTaskCompleted` runs a Lua script that:

1. Adds the child task ID to the completed set.
2. Ignores duplicate child completions.
3. Increments completed or failed counters.
4. Returns the chord callback only once, and only when all children succeeded.

## Idempotency

Workflow advancement is idempotent at the Redis state layer. Task handlers should
still be idempotent because a task may finish, enqueue a downstream workflow
task, and then crash before stream acknowledgement.
