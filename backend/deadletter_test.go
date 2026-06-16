package backend

import (
	"testing"
	"time"

	"github.com/Newton-School/goqueue/task"
)

func TestDeadLetterRequestValidateRequiresMessage(t *testing.T) {
	err := DeadLetterRequest{}.Validate()

	if err == nil {
		t.Fatal("Validate expected error for missing message")
	}
}

func TestDeadLetterRequestValidateAcceptsValidRecord(t *testing.T) {
	err := DeadLetterRequest{
		Message: task.TaskMessage{
			ID:    "task_01JZ9Z8Z8Z8Z8Z8Z8Z8Z8Z8Z8Z",
			Name:  "email.send",
			Queue: "default",
		},
		Reason:         task.FailureExecution,
		Error:          "handler failed",
		SourceStreamID: "1-0",
		Group:          "workers",
		Consumer:       "pod-1",
		FailedAt:       time.Date(2026, time.June, 14, 10, 0, 0, 0, time.UTC),
	}.Validate()

	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestDeadLetterRequestValidateRejectsUnsafeQueue(t *testing.T) {
	err := DeadLetterRequest{
		Message: task.TaskMessage{
			ID:    "malformed-id",
			Name:  "email.send",
			Queue: "bad queue",
		},
		Reason:         task.FailureExecution,
		SourceStreamID: "1-0",
	}.Validate()

	if err == nil {
		t.Fatal("Validate expected error for unsafe queue")
	}
}

func TestReadDeadLettersRequestValidateDefaultsCount(t *testing.T) {
	request := ReadDeadLettersRequest{Queue: "default"}

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestReadDeadLettersRequestRejectsNegativeCount(t *testing.T) {
	request := ReadDeadLettersRequest{Queue: "default", Count: -1}

	if err := request.Validate(); err == nil {
		t.Fatal("Validate expected error for negative count")
	}
}
