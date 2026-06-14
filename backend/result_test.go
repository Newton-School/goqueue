package backend

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/task"
)

func TestTaskResultRecordValidateAcceptsResult(t *testing.T) {
	record := TaskResultRecord{
		TaskID:    "4ac0a01f-1b16-4330-b3e7-e99826eacb1a",
		Result:    task.SucceededResult("ok"),
		UpdatedAt: time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC),
		TTL:       time.Hour,
	}

	if err := record.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestTaskResultRecordValidateRejectsNegativeTTL(t *testing.T) {
	record := TaskResultRecord{
		TaskID: "4ac0a01f-1b16-4330-b3e7-e99826eacb1a",
		Result: task.SucceededResult("ok"),
		TTL:    -time.Second,
	}

	err := record.Validate()
	if !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}
