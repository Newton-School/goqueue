package task

import (
	"strconv"
	"time"
)

const (
	FailureExecution           FailureCategory = "execution"
	FailureMalformedMessage    FailureCategory = "malformed_message"
	FailureUnknownTask         FailureCategory = "unknown_task"
	FailureExpired             FailureCategory = "expired"
	FailureRetryExhausted      FailureCategory = "retry_exhausted"
	FailureRetryScheduleFailed FailureCategory = "retry_schedule_failed"
)

const (
	FailureMetadataCategoryKey       = "goqueue.failure.category"
	FailureMetadataAttemptKey        = "goqueue.failure.attempt"
	FailureMetadataMaxAttemptsKey    = "goqueue.failure.max_attempts"
	FailureMetadataRetryableKey      = "goqueue.failure.retryable"
	FailureMetadataNextRetryAtKey    = "goqueue.failure.next_retry_at"
	FailureMetadataDeadLetteredKey   = "goqueue.failure.dead_lettered"
	FailureMetadataDeadLetteredAtKey = "goqueue.failure.dead_lettered_at"
	FailureMetadataLastErrorKey      = "goqueue.failure.last_error"
)

// FailureCategory classifies why a task execution failed.
type FailureCategory string

// FailureMetadata is stored in TaskResult.Metadata for operational inspection.
type FailureMetadata struct {
	Category       FailureCategory
	Attempt        int
	MaxAttempts    int
	Retryable      bool
	NextRetryAt    time.Time
	DeadLettered   bool
	DeadLetteredAt time.Time
	LastError      string
}

// ToMap returns stable string metadata suitable for TaskResult.Metadata.
func (m FailureMetadata) ToMap() map[string]string {
	values := make(map[string]string)
	if m.Category != "" {
		values[FailureMetadataCategoryKey] = string(m.Category)
	}
	values[FailureMetadataAttemptKey] = strconv.Itoa(m.Attempt)
	values[FailureMetadataMaxAttemptsKey] = strconv.Itoa(m.MaxAttempts)
	values[FailureMetadataRetryableKey] = strconv.FormatBool(m.Retryable)
	if !m.NextRetryAt.IsZero() {
		values[FailureMetadataNextRetryAtKey] = m.NextRetryAt.UTC().Format(time.RFC3339Nano)
	}
	values[FailureMetadataDeadLetteredKey] = strconv.FormatBool(m.DeadLettered)
	if !m.DeadLetteredAt.IsZero() {
		values[FailureMetadataDeadLetteredAtKey] = m.DeadLetteredAt.UTC().Format(time.RFC3339Nano)
	}
	if m.LastError != "" {
		values[FailureMetadataLastErrorKey] = m.LastError
	}
	return values
}
