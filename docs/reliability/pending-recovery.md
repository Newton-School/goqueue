# Pending Recovery

Phase 5 adds optional worker recovery for Redis Stream entries that are pending
because a worker claimed them but did not acknowledge them.

## Claim Flow

When `WithWorkerPendingRecoveryEnabled(true)` is set, the worker periodically
calls `ClaimStaleReady` before normal ready reads. Redis backends implement this
with `XAUTOCLAIM`.
