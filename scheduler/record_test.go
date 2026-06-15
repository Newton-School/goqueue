package scheduler

import (
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestPeriodicTaskToBackendRecordUsesIntervalSchedule(t *testing.T) {
	now := time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC)
	definition := validPeriodicTask()
	definition.Queue = ""

	record, err := definition.toBackendRecord("critical", now)
	if err != nil {
		t.Fatalf("toBackendRecord returned error: %v", err)
	}

	if record.Name != definition.Name.String() {
		t.Fatalf("Name = %q, want %q", record.Name, definition.Name)
	}
	if record.Queue != "critical" {
		t.Fatalf("Queue = %q, want critical", record.Queue)
	}
	if record.ScheduleKind != backend.PeriodicScheduleInterval {
		t.Fatalf("ScheduleKind = %q, want interval", record.ScheduleKind)
	}
	if record.Interval != definition.Schedule.Interval {
		t.Fatalf("Interval = %v, want %v", record.Interval, definition.Schedule.Interval)
	}
	if !record.NextDueAt.Equal(now.Add(10 * time.Minute)) {
		t.Fatalf("NextDueAt = %v, want interval after now", record.NextDueAt)
	}
	if record.Priority != task.DefaultPriority {
		t.Fatalf("Priority = %d, want default", record.Priority)
	}
}

func TestPeriodicTaskFromBackendRecordRestoresDefinition(t *testing.T) {
	record := backend.PeriodicTaskRecord{
		Name:         "welcome-email",
		TaskName:     "email.send",
		Queue:        "critical",
		Args:         []any{"u_123"},
		Kwargs:       map[string]any{"template": "welcome"},
		Metadata:     map[string]string{"source": "scheduler"},
		ScheduleKind: backend.PeriodicScheduleInterval,
		Interval:     10 * time.Minute,
		NextDueAt:    time.Date(2026, time.June, 15, 10, 10, 0, 0, time.UTC),
		Priority:     task.DefaultPriority,
		RetryPolicy:  task.DefaultRetryPolicy(),
	}

	definition, err := periodicTaskFromBackendRecord(record)
	if err != nil {
		t.Fatalf("periodicTaskFromBackendRecord returned error: %v", err)
	}

	if definition.Name != "welcome-email" {
		t.Fatalf("Name = %q, want welcome-email", definition.Name)
	}
	if definition.Queue != "critical" {
		t.Fatalf("Queue = %q, want critical", definition.Queue)
	}
	if definition.Schedule.Interval != 10*time.Minute {
		t.Fatalf("Interval = %v, want 10m", definition.Schedule.Interval)
	}
	if definition.Metadata["source"] != "scheduler" {
		t.Fatalf("Metadata source = %q, want scheduler", definition.Metadata["source"])
	}
}
