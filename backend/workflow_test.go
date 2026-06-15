package backend

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/task"
)

func TestWorkflowSignatureRecordValidateAcceptsCompleteSignature(t *testing.T) {
	record := validWorkflowSignatureRecord()

	if err := record.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestWorkflowSignatureRecordValidateRequiresTaskName(t *testing.T) {
	record := validWorkflowSignatureRecord()
	record.Name = ""

	if err := record.Validate(); !errors.Is(err, task.ErrInvalidTaskName) {
		t.Fatalf("Validate error = %v, want ErrInvalidTaskName", err)
	}
}

func TestWorkflowSignatureRecordValidateRequiresQueue(t *testing.T) {
	record := validWorkflowSignatureRecord()
	record.Queue = ""

	if err := record.Validate(); !errors.Is(err, task.ErrInvalidQueueName) {
		t.Fatalf("Validate error = %v, want ErrInvalidQueueName", err)
	}
}

func validWorkflowSignatureRecord() WorkflowSignatureRecord {
	return WorkflowSignatureRecord{
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
