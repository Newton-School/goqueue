package redisbackend

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestPeriodicTaskCodecRoundTripsRecord(t *testing.T) {
	record := testPeriodicTaskRecord()

	encoded, err := (periodicTaskCodec{}).encode(record)
	if err != nil {
		t.Fatalf("encode returned error: %v", err)
	}

	decoded, err := (periodicTaskCodec{}).decode(encoded)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.Name != record.Name {
		t.Fatalf("Name = %q, want %q", decoded.Name, record.Name)
	}
	if decoded.TaskName != record.TaskName {
		t.Fatalf("TaskName = %q, want %q", decoded.TaskName, record.TaskName)
	}
	if decoded.Queue != record.Queue {
		t.Fatalf("Queue = %q, want %q", decoded.Queue, record.Queue)
	}
	if decoded.Interval != record.Interval {
		t.Fatalf("Interval = %v, want %v", decoded.Interval, record.Interval)
	}
	if !decoded.NextDueAt.Equal(record.NextDueAt) {
		t.Fatalf("NextDueAt = %v, want %v", decoded.NextDueAt, record.NextDueAt)
	}
	if decoded.Metadata["source"] != "scheduler" {
		t.Fatalf("Metadata source = %q, want scheduler", decoded.Metadata["source"])
	}
}

func TestPeriodicTaskCodecRejectsInvalidJSON(t *testing.T) {
	_, err := (periodicTaskCodec{}).decode([]byte("{"))
	if !errors.Is(err, ErrInvalidRedisMessage) {
		t.Fatalf("decode error = %v, want ErrInvalidRedisMessage", err)
	}
}

func TestPeriodicTaskCodecRejectsInvalidRecord(t *testing.T) {
	record := testPeriodicTaskRecord()
	record.Name = ""

	_, err := (periodicTaskCodec{}).encode(record)
	if !errors.Is(err, backend.ErrInvalidBackendRequest) {
		t.Fatalf("encode error = %v, want ErrInvalidBackendRequest", err)
	}
}

func testPeriodicTaskRecord() backend.PeriodicTaskRecord {
	return backend.PeriodicTaskRecord{
		Name:         "welcome-email",
		TaskName:     "email.send",
		Queue:        "default",
		Args:         []any{"u_123"},
		Kwargs:       map[string]any{"template": "welcome"},
		Metadata:     map[string]string{"source": "scheduler"},
		ScheduleKind: backend.PeriodicScheduleInterval,
		Interval:     10 * time.Minute,
		NextDueAt:    time.Date(2026, time.June, 15, 10, 10, 0, 0, time.UTC),
		Priority:     task.DefaultPriority,
		RetryPolicy:  task.DefaultRetryPolicy(),
		UpdatedAt:    time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC),
	}
}
