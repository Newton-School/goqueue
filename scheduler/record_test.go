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
