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

func TestWorkflowChainRecordValidateAcceptsCompleteRecord(t *testing.T) {
	record := validWorkflowChainRecord()

	if err := record.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestWorkflowChainRecordValidateRequiresID(t *testing.T) {
	record := validWorkflowChainRecord()
	record.ID = ""

	if err := record.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestWorkflowChainRecordValidateRequiresSignatures(t *testing.T) {
	record := validWorkflowChainRecord()
	record.Signatures = nil

	if err := record.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestAdvanceWorkflowChainRequestValidateAcceptsCompleteRequest(t *testing.T) {
	request := AdvanceWorkflowChainRequest{
		WorkflowID:      "chain-1",
		CompletedTaskID: "11111111-1111-4111-8111-111111111111",
		CompletedIndex:  0,
		CompletedAt:     time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
	}

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestAdvanceWorkflowChainRequestValidateRejectsNegativeIndex(t *testing.T) {
	request := AdvanceWorkflowChainRequest{
		WorkflowID:      "chain-1",
		CompletedTaskID: "11111111-1111-4111-8111-111111111111",
		CompletedIndex:  -1,
		CompletedAt:     time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
	}

	if err := request.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestWorkflowGroupRecordValidateAcceptsCompleteRecord(t *testing.T) {
	record := validWorkflowGroupRecord()

	if err := record.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestWorkflowGroupRecordValidateRequiresTaskIDs(t *testing.T) {
	record := validWorkflowGroupRecord()
	record.TaskIDs = nil

	if err := record.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestRecordWorkflowTaskCompletedRequestValidateAcceptsTerminalState(t *testing.T) {
	request := RecordWorkflowTaskCompletedRequest{
		GroupID:     "group-1",
		TaskID:      "11111111-1111-4111-8111-111111111111",
		State:       task.TaskSucceeded,
		CompletedAt: time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
	}

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestRecordWorkflowTaskCompletedRequestValidateRejectsNonTerminalState(t *testing.T) {
	request := RecordWorkflowTaskCompletedRequest{
		GroupID:     "group-1",
		TaskID:      "11111111-1111-4111-8111-111111111111",
		State:       task.TaskStarted,
		CompletedAt: time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
	}

	if err := request.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func validWorkflowGroupRecord() WorkflowGroupRecord {
	return WorkflowGroupRecord{
		ID: "group-1",
		TaskIDs: []task.TaskID{
			"11111111-1111-4111-8111-111111111111",
			"22222222-2222-4222-8222-222222222222",
		},
		CreatedAt: time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC),
	}
}

func validWorkflowChainRecord() WorkflowChainRecord {
	return WorkflowChainRecord{
		ID:         "chain-1",
		Signatures: []WorkflowSignatureRecord{validWorkflowSignatureRecord(), validWorkflowSignatureRecord()},
		CreatedAt:  time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC),
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
