package worker

import (
	"context"
	"time"

	"github.com/Newton-School/goqueue/backend"
)

func (w *Worker) shouldClaimPending(now time.Time, lastClaim time.Time) bool {
	if !w.pendingRecoveryEnabled {
		return false
	}
	if lastClaim.IsZero() {
		return true
	}
	return !now.Before(lastClaim.Add(w.pendingClaimInterval))
}

func (w *Worker) claimStalePending(ctx context.Context) ([]backend.ReadyMessage, error) {
	return w.backend.ClaimStaleReady(ctx, backend.ClaimStaleReadyRequest{
		Queue:    w.queue,
		Group:    w.group,
		Consumer: w.consumer,
		MinIdle:  w.pendingMinIdle,
		Count:    w.pendingClaimBatch,
		StartID:  "0-0",
	})
}
