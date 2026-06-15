package workflow

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/task"
)

func TestSignatureValidateRequiresTaskName(t *testing.T) {
	signature := validSignature()
	signature.Name = ""

	if err := signature.Validate(); !errors.Is(err, task.ErrInvalidTaskName) {
		t.Fatalf("Validate error = %v, want ErrInvalidTaskName", err)
	}
}

func TestSignatureValidateRejectsInvalidQueueWhenSet(t *testing.T) {
	signature := validSignature()
	signature.Queue = "invalid queue"

	if err := signature.Validate(); !errors.Is(err, task.ErrInvalidQueueName) {
		t.Fatalf("Validate error = %v, want ErrInvalidQueueName", err)
	}
}

func TestSignatureNormalizeAppliesDefaults(t *testing.T) {
	signature := validSignature()
	signature.Queue = ""
	signature.Priority = 0
	signature.RetryPolicy = task.RetryPolicy{}

	normalized, err := signature.Normalize("critical")
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}

	if normalized.Queue != "critical" {
		t.Fatalf("Queue = %q, want critical", normalized.Queue)
	}
	if normalized.Priority != task.DefaultPriority {
		t.Fatalf("Priority = %d, want default", normalized.Priority)
	}
	if normalized.RetryPolicy != task.DefaultRetryPolicy() {
		t.Fatalf("RetryPolicy = %+v, want default", normalized.RetryPolicy)
	}
}

func TestSignatureNormalizeCopiesMutableFields(t *testing.T) {
	signature := validSignature()

	normalized, err := signature.Normalize("default")
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}

	signature.Args[0] = "u_999"
	signature.Kwargs["template"] = "changed"
	signature.Metadata["source"] = "changed"

	if normalized.Args[0] != "u_123" {
		t.Fatalf("Args[0] = %v, want copied value", normalized.Args[0])
	}
	if normalized.Kwargs["template"] != "welcome" {
		t.Fatalf("Kwargs template = %v, want copied value", normalized.Kwargs["template"])
	}
	if normalized.Metadata["source"] != "workflow" {
		t.Fatalf("Metadata source = %v, want copied value", normalized.Metadata["source"])
	}
}

func validSignature() Signature {
	return Signature{
		Name:        "email.send",
		Queue:       "default",
		Args:        []any{"u_123"},
		Kwargs:      map[string]any{"template": "welcome"},
		Metadata:    map[string]string{"source": "workflow"},
		Timing:      task.TaskTiming{ExpiresAt: time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC)},
		Priority:    task.DefaultPriority,
		RetryPolicy: task.DefaultRetryPolicy(),
	}
}
