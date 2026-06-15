# Workflow Failure Semantics

Phase 7 workflow primitives keep failure handling explicit and queue-safe.

## Chains

Chains advance only when the current task reaches `SUCCEEDED`. If a chain task
fails, retries, expires, or lands in a dead-letter path, the worker does not
dispatch the next chain signature.

## Groups

Groups record every terminal child state. Successful children increment the
completed count. Failed terminal children increment the failed count. This lets
the group finish without waiting forever for a child that already reached a
terminal state.

## Chords

Chords are groups with a callback. The callback dispatches only when all header
tasks succeed. If any header task fails, the chord completes as failed and the
callback is not dispatched.
