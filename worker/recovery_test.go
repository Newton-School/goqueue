package worker

import (
	"testing"
	"time"
)

func TestWorkerShouldClaimPendingHonorsInterval(t *testing.T) {
	worker := &Worker{
		pendingRecoveryEnabled: true,
		pendingClaimInterval:   time.Minute,
	}
	now := time.Date(2026, time.June, 14, 10, 0, 0, 0, time.UTC)

	if !worker.shouldClaimPending(now, time.Time{}) {
		t.Fatal("first claim should run")
	}
	if worker.shouldClaimPending(now.Add(30*time.Second), now) {
		t.Fatal("claim should not run before interval")
	}
	if !worker.shouldClaimPending(now.Add(time.Minute), now) {
		t.Fatal("claim should run at interval")
	}
}

func TestWorkerShouldClaimPendingRespectsDisabledRecovery(t *testing.T) {
	worker := &Worker{pendingRecoveryEnabled: false}
	now := time.Date(2026, time.June, 14, 10, 0, 0, 0, time.UTC)

	if worker.shouldClaimPending(now, time.Time{}) {
		t.Fatal("claim should not run when recovery is disabled")
	}
}
