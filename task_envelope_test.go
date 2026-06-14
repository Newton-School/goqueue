package goqueue

import (
	"errors"
	"testing"
	"time"
)

func TestNewTaskEnvelopeAppliesDefaults(t *testing.T) {
	createdAt := time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC)

	envelope, err := NewTaskEnvelope(TaskEnvelopeInput{
		Name:      "email.send_welcome",
		Queue:     "default",
		CreatedAt: createdAt,
	})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}

	if envelope.ID == "" {
		t.Fatal("ID should be generated")
	}
	if envelope.Priority != DefaultPriority {
		t.Fatalf("Priority = %d, want %d", envelope.Priority, DefaultPriority)
	}
	if envelope.RetryPolicy != DefaultRetryPolicy() {
		t.Fatalf("RetryPolicy = %+v, want default", envelope.RetryPolicy)
	}
	if !envelope.CreatedAt.Equal(createdAt) {
		t.Fatalf("CreatedAt = %s, want %s", envelope.CreatedAt, createdAt)
	}
}

func TestNewTaskEnvelopePreservesMinimumPriority(t *testing.T) {
	envelope, err := NewTaskEnvelope(TaskEnvelopeInput{
		Name:     "email.send_welcome",
		Queue:    "default",
		Priority: MinPriority,
	})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}

	if envelope.Priority != MinPriority {
		t.Fatalf("Priority = %d, want %d", envelope.Priority, MinPriority)
	}
}

func TestNewTaskEnvelopeRejectsInvalidName(t *testing.T) {
	_, err := NewTaskEnvelope(TaskEnvelopeInput{
		Name:  "email/send",
		Queue: "default",
	})
	if !errors.Is(err, ErrInvalidTaskName) {
		t.Fatalf("NewTaskEnvelope error = %v, want ErrInvalidTaskName", err)
	}
}

func TestNewTaskEnvelopeRejectsInvalidQueue(t *testing.T) {
	_, err := NewTaskEnvelope(TaskEnvelopeInput{
		Name:  "email.send",
		Queue: "default queue",
	})
	if !errors.Is(err, ErrInvalidQueueName) {
		t.Fatalf("NewTaskEnvelope error = %v, want ErrInvalidQueueName", err)
	}
}

func TestNewTaskEnvelopeRejectsInvalidRetryPolicy(t *testing.T) {
	_, err := NewTaskEnvelope(TaskEnvelopeInput{
		Name:        "email.send",
		Queue:       "default",
		RetryPolicy: RetryPolicy{MaxAttempts: -1},
	})
	if !errors.Is(err, ErrInvalidRetryPolicy) {
		t.Fatalf("NewTaskEnvelope error = %v, want ErrInvalidRetryPolicy", err)
	}
}

func TestTaskEnvelopeCloneCopiesPayloadAndMetadata(t *testing.T) {
	envelope, err := NewTaskEnvelope(TaskEnvelopeInput{
		Name:     "email.send",
		Queue:    "default",
		Args:     []any{"welcome"},
		Metadata: map[string]string{"trace_id": "trace-1"},
	})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}

	cloned := envelope.Clone()
	clonedArgs := cloned.Payload.Args()
	clonedArgs[0] = "mutated"
	clonedMetadata := cloned.Metadata.Values()
	clonedMetadata["trace_id"] = "trace-2"

	if got := envelope.Payload.Args()[0]; got != "welcome" {
		t.Fatalf("original payload arg = %v, want welcome", got)
	}
	if got := envelope.Metadata.Values()["trace_id"]; got != "trace-1" {
		t.Fatalf("original metadata trace_id = %q, want trace-1", got)
	}
}
