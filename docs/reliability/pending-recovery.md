# Pending Recovery

Phase 5 adds optional worker recovery for Redis Stream entries that are pending
because a worker claimed them but did not acknowledge them.

## Claim Flow

When `WithWorkerPendingRecoveryEnabled(true)` is set, the worker periodically
calls `ClaimStaleReady` before normal ready reads. Redis backends implement this
with `XAUTOCLAIM`.

## Tuning

Use `WithWorkerPendingMinIdle`, `WithWorkerPendingClaimBatch`, and
`WithWorkerPendingClaimInterval` to control when a pending message is considered
stale, how many entries are claimed per pass, and how often recovery runs.

## Failure Behavior

If pending recovery fails, `Worker.Start` returns the claim error and does not
fall through to fresh reads. This keeps Redis recovery failures visible to the
owning process.
