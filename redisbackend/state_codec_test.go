package redisbackend

import (
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestStateCodecRoundTripsRecord(t *testing.T) {
	record := backend.TaskStateRecord{
		TaskID:    "4ac0a01f-1b16-4330-b3e7-e99826eacb1a",
		State:     task.TaskStarted,
		UpdatedAt: time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC),
	}

	encoded, err := (stateCodec{}).encode(record)
	if err != nil {
		t.Fatalf("encode returned error: %v", err)
	}

	decoded, err := (stateCodec{}).decode(encoded)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.TaskID != record.TaskID {
		t.Fatalf("TaskID = %q, want %q", decoded.TaskID, record.TaskID)
	}
	if decoded.State != record.State {
		t.Fatalf("State = %s, want %s", decoded.State, record.State)
	}
}
