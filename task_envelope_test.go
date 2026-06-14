package goqueue

import (
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
