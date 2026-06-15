package scheduler

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/task"
)

func TestDefaultSchedulerConfigUsesProductionDefaults(t *testing.T) {
	config := defaultSchedulerConfig()

	if config.defaultQueue != "default" {
		t.Fatalf("defaultQueue = %q, want default", config.defaultQueue)
	}
	if config.pollInterval != time.Second {
		t.Fatalf("pollInterval = %v, want 1s", config.pollInterval)
	}
	if config.batchSize != 100 {
		t.Fatalf("batchSize = %d, want 100", config.batchSize)
	}
	if config.lockTTL != 30*time.Second {
		t.Fatalf("lockTTL = %v, want 30s", config.lockTTL)
	}
	if config.now == nil {
		t.Fatal("now should be set")
	}
}

func TestSchedulerOptionsValidateValues(t *testing.T) {
	config := defaultSchedulerConfig()
	options := []SchedulerOption{
		WithSchedulerIdentity("scheduler-1"),
		WithSchedulerDefaultQueue("critical"),
		WithSchedulerPollInterval(2 * time.Second),
		WithSchedulerBatchSize(25),
		WithSchedulerLockTTL(time.Minute),
		WithSchedulerCodec(task.JSONPayloadCodec{}),
		WithSchedulerNow(func() time.Time { return time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC) }),
	}

	for _, option := range options {
		if err := option(&config); err != nil {
			t.Fatalf("option returned error: %v", err)
		}
	}

	if config.identity != "scheduler-1" {
		t.Fatalf("identity = %q, want scheduler-1", config.identity)
	}
	if config.defaultQueue != "critical" {
		t.Fatalf("defaultQueue = %q, want critical", config.defaultQueue)
	}
	if config.pollInterval != 2*time.Second {
		t.Fatalf("pollInterval = %v, want 2s", config.pollInterval)
	}
	if config.batchSize != 25 {
		t.Fatalf("batchSize = %d, want 25", config.batchSize)
	}
	if config.lockTTL != time.Minute {
		t.Fatalf("lockTTL = %v, want 1m", config.lockTTL)
	}
}

func TestSchedulerIdentityOptionTrimsWhitespace(t *testing.T) {
	config := defaultSchedulerConfig()

	if err := WithSchedulerIdentity(" scheduler-1 ")(&config); err != nil {
		t.Fatalf("WithSchedulerIdentity returned error: %v", err)
	}

	if config.identity != "scheduler-1" {
		t.Fatalf("identity = %q, want scheduler-1", config.identity)
	}
}

func TestSchedulerOptionsRejectInvalidValues(t *testing.T) {
	tests := []struct {
		name   string
		option SchedulerOption
	}{
		{name: "empty identity", option: WithSchedulerIdentity("")},
		{name: "blank identity", option: WithSchedulerIdentity(" \t ")},
		{name: "invalid queue", option: WithSchedulerDefaultQueue("invalid queue")},
		{name: "non-positive poll", option: WithSchedulerPollInterval(0)},
		{name: "non-positive batch", option: WithSchedulerBatchSize(0)},
		{name: "non-positive lock ttl", option: WithSchedulerLockTTL(0)},
		{name: "nil codec", option: WithSchedulerCodec(nil)},
		{name: "nil now", option: WithSchedulerNow(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := defaultSchedulerConfig()
			if err := tt.option(&config); !errors.Is(err, ErrInvalidSchedulerOption) {
				t.Fatalf("option error = %v, want ErrInvalidSchedulerOption", err)
			}
		})
	}
}
