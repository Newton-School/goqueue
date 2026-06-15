package backend

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/task"
)

func TestPeriodicTaskRecordValidateAcceptsCompleteIntervalRecord(t *testing.T) {
	record := validPeriodicTaskRecord()

	if err := record.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestPeriodicTaskRecordValidateRequiresName(t *testing.T) {
	record := validPeriodicTaskRecord()
	record.Name = ""

	if err := record.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestPeriodicTaskRecordValidateRequiresTaskName(t *testing.T) {
	record := validPeriodicTaskRecord()
	record.TaskName = ""

	if err := record.Validate(); !errors.Is(err, task.ErrInvalidTaskName) {
		t.Fatalf("Validate error = %v, want ErrInvalidTaskName", err)
	}
}

func TestPeriodicTaskRecordValidateRequiresQueue(t *testing.T) {
	record := validPeriodicTaskRecord()
	record.Queue = ""

	if err := record.Validate(); !errors.Is(err, task.ErrInvalidQueueName) {
		t.Fatalf("Validate error = %v, want ErrInvalidQueueName", err)
	}
}

func TestPeriodicTaskRecordValidateRequiresIntervalSchedule(t *testing.T) {
	record := validPeriodicTaskRecord()
	record.ScheduleKind = "cron"

	if err := record.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestPeriodicTaskRecordValidateRequiresPositiveInterval(t *testing.T) {
	record := validPeriodicTaskRecord()
	record.Interval = 0

	if err := record.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestPeriodicTaskRecordValidateRequiresNextDueAt(t *testing.T) {
	record := validPeriodicTaskRecord()
	record.NextDueAt = time.Time{}

	if err := record.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func validPeriodicTaskRecord() PeriodicTaskRecord {
	return PeriodicTaskRecord{
		Name:         "welcome-email",
		TaskName:     "email.send",
		Queue:        "default",
		Args:         []any{"u_123"},
		Kwargs:       map[string]any{"template": "welcome"},
		Metadata:     map[string]string{"source": "scheduler"},
		ScheduleKind: PeriodicScheduleInterval,
		Interval:     10 * time.Minute,
		NextDueAt:    time.Date(2026, time.June, 15, 10, 10, 0, 0, time.UTC),
		Priority:     task.DefaultPriority,
		RetryPolicy:  task.DefaultRetryPolicy(),
	}
}
