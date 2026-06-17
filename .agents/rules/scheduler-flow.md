# Scheduler Flow Rules

Use this file when changing `scheduler/`, periodic task definitions, fixed
interval schedules, Redis scheduler coordination, or scheduler documentation.

## Responsibilities

- The scheduler stores periodic task definitions through backend contracts.
- The scheduler dispatches due periodic tasks through `producer.Producer`.
- Redis coordination prevents multiple scheduler pods from dispatching the same
  due definition at the same time.
- The current public scheduler supports fixed intervals. Cron-style schedules
  are reserved for a future release unless the user explicitly asks to add them.

## Registration Flow

1. Validate periodic task name, task name, queue, schedule, priority, retry
   policy, metadata, and timing fields.
2. Use the scheduler default queue when a definition omits one.
3. Convert the definition into a backend record.
4. Store the definition and next due timestamp through
   `UpsertPeriodicTask`.

## Dispatch Flow

1. `PollOnce` lists due periodic tasks with scheduler identity, batch size, and
   lock TTL.
2. Each due record is validated and converted back into a periodic definition.
3. Dispatch uses the scheduler-owned producer so payload codec, queue defaults,
   retry policy, priority, and created time stay consistent.
4. Metadata must include periodic task name and due timestamp while preserving
   user metadata.
5. Mark dispatch with the backend lock token, dispatched task ID, dispatch time,
   and next due timestamp.

## Coordination

- Scheduler identity may be caller-provided or generated from cryptographic
  randomness.
- Poll interval, batch size, and lock TTL must be positive.
- Lock ownership must be checked before advancing a periodic definition.
- Do not bypass backend lease methods with direct Redis calls from `scheduler/`.

## Documentation

- Keep scheduler and periodic task behavior docs current in
  `docs/docs/concepts/scheduler.md`.
- If default interval behavior or retry lock semantics change, update workflow and
  operation docs as needed.
- Apply `docs-sync.md` when scheduler-facing behavior changes.

## Tests

- Use deterministic clocks through `WithSchedulerNow`.
- Cover registration, deletion, due polling, metadata merge, default queue,
  custom queue, batch size, lock TTL, generated identity, and dispatch failure.
- Add Redis integration coverage for lease contention or script behavior.
