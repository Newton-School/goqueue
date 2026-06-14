package backend

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/task"
)

func TestTaskStateRecordValidateAcceptsState(t *testing.T) {
	record := TaskStateRecord{
		TaskID:    "4ac0a01f-1b16-4330-b3e7-e99826eacb1a",
		State:     task.TaskStarted,
		UpdatedAt: time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC),
	}

	if err := record.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestTaskStateRecordValidateRejectsMissingTaskID(t *testing.T) {
	err := (TaskStateRecord{State: task.TaskStarted}).Validate()
	if !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}
