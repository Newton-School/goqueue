package worker

import (
	"testing"
	"time"
)

func TestDefaultWorkerClockAdvances(t *testing.T) {
	config := defaultWorkerConfig()
	first := config.now()
	time.Sleep(2 * time.Millisecond)
	second := config.now()

	if !second.After(first) {
		t.Fatalf("default worker clock did not advance: first=%s second=%s", first.Format(time.RFC3339Nano), second.Format(time.RFC3339Nano))
	}
}

func TestWorkerReliabilityOptions(t *testing.T) {
	config := defaultWorkerConfig()

	options := []WorkerOption{
		WithWorkerDeadLetterEnabled(false),
		WithWorkerPendingRecoveryEnabled(true),
		WithWorkerPendingMinIdle(5 * time.Minute),
		WithWorkerPendingClaimBatch(20),
		WithWorkerPendingClaimInterval(30 * time.Second),
	}

	for _, opt := range options {
		if err := opt(&config); err != nil {
			t.Fatalf("option returned error: %v", err)
		}
	}

	if config.deadLetterEnabled {
		t.Fatal("dead letter should be disabled")
	}
	if !config.pendingRecoveryEnabled {
		t.Fatal("pending recovery should be enabled")
	}
	if config.pendingMinIdle != 5*time.Minute {
		t.Fatalf("pending min idle = %v, want 5m", config.pendingMinIdle)
	}
	if config.pendingClaimBatch != 20 {
		t.Fatalf("pending claim batch = %d, want 20", config.pendingClaimBatch)
	}
	if config.pendingClaimInterval != 30*time.Second {
		t.Fatalf("pending claim interval = %v, want 30s", config.pendingClaimInterval)
	}
}

func TestWorkerReliabilityOptionsRejectInvalidValues(t *testing.T) {
	tests := []struct {
		name string
		opt  WorkerOption
	}{
		{name: "pending min idle", opt: WithWorkerPendingMinIdle(-time.Second)},
		{name: "pending claim batch", opt: WithWorkerPendingClaimBatch(0)},
		{name: "pending claim interval", opt: WithWorkerPendingClaimInterval(-time.Second)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config := defaultWorkerConfig()
			if err := tc.opt(&config); err == nil {
				t.Fatal("option expected error")
			}
		})
	}
}
