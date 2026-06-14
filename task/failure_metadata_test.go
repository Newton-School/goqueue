package task

import (
	"testing"
	"time"
)

func TestFailureMetadataToMapIncludesStableFields(t *testing.T) {
	nextRetryAt := time.Date(2026, time.June, 14, 10, 30, 0, 0, time.UTC)
	deadLetteredAt := time.Date(2026, time.June, 14, 10, 35, 0, 0, time.UTC)

	metadata := FailureMetadata{
		Category:       FailureExecution,
		Attempt:        2,
		MaxAttempts:    3,
		Retryable:      true,
		NextRetryAt:    nextRetryAt,
		DeadLettered:   true,
		DeadLetteredAt: deadLetteredAt,
		LastError:      "handler failed",
	}

	values := metadata.ToMap()

	if values["goqueue.failure.category"] != string(FailureExecution) {
		t.Fatalf("category = %q, want %q", values["goqueue.failure.category"], FailureExecution)
	}
	if values["goqueue.failure.attempt"] != "2" {
		t.Fatalf("attempt = %q, want 2", values["goqueue.failure.attempt"])
	}
	if values["goqueue.failure.max_attempts"] != "3" {
		t.Fatalf("max attempts = %q, want 3", values["goqueue.failure.max_attempts"])
	}
	if values["goqueue.failure.retryable"] != "true" {
		t.Fatalf("retryable = %q, want true", values["goqueue.failure.retryable"])
	}
	if values["goqueue.failure.next_retry_at"] != nextRetryAt.Format(time.RFC3339Nano) {
		t.Fatalf("next retry = %q, want RFC3339 timestamp", values["goqueue.failure.next_retry_at"])
	}
	if values["goqueue.failure.dead_lettered"] != "true" {
		t.Fatalf("dead lettered = %q, want true", values["goqueue.failure.dead_lettered"])
	}
	if values["goqueue.failure.dead_lettered_at"] != deadLetteredAt.Format(time.RFC3339Nano) {
		t.Fatalf("dead lettered at = %q, want RFC3339 timestamp", values["goqueue.failure.dead_lettered_at"])
	}
	if values["goqueue.failure.last_error"] != "handler failed" {
		t.Fatalf("last error = %q, want handler failed", values["goqueue.failure.last_error"])
	}
}
