---
title: Errors
---

Use exported errors to decide operational behavior and user feedback.

## Setup/validation errors

- invalid redis URL or missing Redis URL
- invalid task/queue/name formats
- invalid task state, retry policy, timing, priority

## Producer/consumer runtime errors

- `goqueue.ErrNilBackend`
- `goqueue.ErrNilWorker`
- `goqueue.ErrNilTaskRegistry`
- `goqueue.ErrMissingTaskName`

## Scheduler and workflow

- `goqueue.ErrInvalidSchedule`
- `goqueue.ErrInvalidPeriodicTask`
- `goqueue.ErrInvalidSchedulerOption`
- `goqueue.ErrInvalidSignature`
- `goqueue.ErrInvalidWorkflow`

## Inspect and admin

- `goqueue.ErrNilInspector`
- `goqueue.ErrNilAdmin`
- `goqueue.ErrAdminBackend`
- `goqueue.ErrTaskMessageNotFound`
- `goqueue.ErrDeadLetterNotFound`
- `goqueue.ErrTaskStateNotFound`
- `goqueue.ErrTaskResultNotFound`

## Behavior guidance

- treat IDs and queue names as user input and validate before calling any API
- prefer inspect output checks over raw Redis access for health and status
- keep credentials out of error payloads
