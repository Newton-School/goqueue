package scheduler

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/task"
)

func TestPeriodicTaskValidateRequiresName(t *testing.T) {
	definition := validPeriodicTask()
	definition.Name = ""

	if err := definition.Validate(); !errors.Is(err, ErrInvalidPeriodicTask) {
		t.Fatalf("Validate error = %v, want ErrInvalidPeriodicTask", err)
	}
}

func TestPeriodicTaskValidateRequiresTaskName(t *testing.T) {
	definition := validPeriodicTask()
	definition.TaskName = ""

	if err := definition.Validate(); !errors.Is(err, task.ErrInvalidTaskName) {
		t.Fatalf("Validate error = %v, want ErrInvalidTaskName", err)
	}
}

func TestPeriodicTaskValidateRequiresValidQueueWhenSet(t *testing.T) {
	definition := validPeriodicTask()
	definition.Queue = "invalid queue"

	if err := definition.Validate(); !errors.Is(err, task.ErrInvalidQueueName) {
		t.Fatalf("Validate error = %v, want ErrInvalidQueueName", err)
	}
}

func TestPeriodicTaskValidateRequiresSchedule(t *testing.T) {
	definition := validPeriodicTask()
	definition.Schedule = IntervalSchedule{}

	if err := definition.Validate(); !errors.Is(err, ErrInvalidSchedule) {
		t.Fatalf("Validate error = %v, want ErrInvalidSchedule", err)
	}
}

func TestPeriodicTaskNormalizeAppliesSafeDefaults(t *testing.T) {
	definition := validPeriodicTask()
	definition.Queue = ""
	definition.Priority = 0
	definition.RetryPolicy = task.RetryPolicy{}

	normalized, err := definition.Normalize("critical")
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}

	if normalized.Queue != "critical" {
		t.Fatalf("Queue = %q, want critical", normalized.Queue)
	}
	if normalized.Priority != task.DefaultPriority {
		t.Fatalf("Priority = %d, want default", normalized.Priority)
	}
	if normalized.RetryPolicy != task.DefaultRetryPolicy() {
		t.Fatalf("RetryPolicy = %+v, want default", normalized.RetryPolicy)
	}
}

func TestPeriodicTaskNormalizeCopiesMutableFields(t *testing.T) {
	definition := validPeriodicTask()

	normalized, err := definition.Normalize("default")
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}

	definition.Args[0] = "u_999"
	definition.Kwargs["template"] = "changed"
	definition.Metadata["source"] = "changed"

	if normalized.Args[0] != "u_123" {
		t.Fatalf("Args[0] = %v, want copied value", normalized.Args[0])
	}
	if normalized.Kwargs["template"] != "welcome" {
		t.Fatalf("Kwargs template = %v, want copied value", normalized.Kwargs["template"])
	}
	if normalized.Metadata["source"] != "scheduler" {
		t.Fatalf("Metadata source = %v, want copied value", normalized.Metadata["source"])
	}
}

func TestPeriodicTaskFirstDueDefaultsToOneIntervalAfterNow(t *testing.T) {
	now := time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC)
	definition := validPeriodicTask()

	next := definition.FirstDueAfter(now)

	want := now.Add(10 * time.Minute)
	if !next.Equal(want) {
		t.Fatalf("FirstDueAfter = %v, want %v", next, want)
	}
}

func TestPeriodicTaskFirstDueUsesFutureStartAt(t *testing.T) {
	now := time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC)
	startAt := now.Add(2 * time.Minute)
	definition := validPeriodicTask()
	definition.StartAt = startAt

	next := definition.FirstDueAfter(now)

	if !next.Equal(startAt) {
		t.Fatalf("FirstDueAfter = %v, want %v", next, startAt)
	}
}

func TestPeriodicTaskFirstDueUsesNowWhenStartAtHasPassed(t *testing.T) {
	now := time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC)
	definition := validPeriodicTask()
	definition.StartAt = now.Add(-time.Minute)

	next := definition.FirstDueAfter(now)

	if !next.Equal(now) {
		t.Fatalf("FirstDueAfter = %v, want %v", next, now)
	}
}

func TestPeriodicTaskNextDueAfterUsesSchedule(t *testing.T) {
	now := time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC)
	definition := validPeriodicTask()

	next := definition.NextDueAfter(now)

	want := now.Add(10 * time.Minute)
	if !next.Equal(want) {
		t.Fatalf("NextDueAfter = %v, want %v", next, want)
	}
}

func validPeriodicTask() PeriodicTask {
	return PeriodicTask{
		Name:        "welcome-email",
		TaskName:    "email.send",
		Queue:       "default",
		Args:        []any{"u_123"},
		Kwargs:      map[string]any{"template": "welcome"},
		Metadata:    map[string]string{"source": "scheduler"},
		Schedule:    Every(10 * time.Minute),
		Priority:    task.DefaultPriority,
		RetryPolicy: task.DefaultRetryPolicy(),
	}
}
