package redisbackend

import (
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestResultCodecRoundTripsRecord(t *testing.T) {
	record := backend.TaskResultRecord{
		TaskID:    "4ac0a01f-1b16-4330-b3e7-e99826eacb1a",
		Result:    task.SucceededResult("ok"),
		UpdatedAt: time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC),
	}

	encoded, err := (resultCodec{}).encode(record)
	if err != nil {
		t.Fatalf("encode returned error: %v", err)
	}

	decoded, err := (resultCodec{}).decode(encoded)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.TaskID != record.TaskID {
		t.Fatalf("TaskID = %q, want %q", decoded.TaskID, record.TaskID)
	}
	if decoded.Result.State != task.TaskSucceeded {
		t.Fatalf("State = %s, want %s", decoded.Result.State, task.TaskSucceeded)
	}
}
