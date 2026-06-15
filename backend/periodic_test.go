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

func TestUpsertPeriodicTaskRequestValidateUsesRecordValidation(t *testing.T) {
	request := UpsertPeriodicTaskRequest{Record: validPeriodicTaskRecord()}

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}

	request.Record.Name = ""
	if err := request.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestListDuePeriodicTasksRequestValidateAcceptsCompleteRequest(t *testing.T) {
	request := ListDuePeriodicTasksRequest{
		Now:         time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC),
		Limit:       10,
		SchedulerID: "scheduler-1",
		LockTTL:     time.Minute,
	}

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestListDuePeriodicTasksRequestValidateRequiresNow(t *testing.T) {
	request := ListDuePeriodicTasksRequest{
		Limit:       10,
		SchedulerID: "scheduler-1",
		LockTTL:     time.Minute,
	}

	if err := request.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestListDuePeriodicTasksRequestValidateRequiresPositiveLimit(t *testing.T) {
	request := ListDuePeriodicTasksRequest{
		Now:         time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC),
		Limit:       0,
		SchedulerID: "scheduler-1",
		LockTTL:     time.Minute,
	}

	if err := request.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestListDuePeriodicTasksRequestValidateRequiresSchedulerID(t *testing.T) {
	request := ListDuePeriodicTasksRequest{
		Now:     time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC),
		Limit:   10,
		LockTTL: time.Minute,
	}

	if err := request.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestListDuePeriodicTasksRequestValidateRequiresPositiveLockTTL(t *testing.T) {
	request := ListDuePeriodicTasksRequest{
		Now:         time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC),
		Limit:       10,
		SchedulerID: "scheduler-1",
		LockTTL:     0,
	}

	if err := request.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestDuePeriodicTaskValidateAcceptsCompleteRecord(t *testing.T) {
	due := DuePeriodicTask{
		Record:      validPeriodicTaskRecord(),
		LockToken:   "token",
		LockedUntil: time.Date(2026, time.June, 15, 10, 1, 0, 0, time.UTC),
	}

	if err := due.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestDuePeriodicTaskValidateRequiresLockToken(t *testing.T) {
	due := DuePeriodicTask{
		Record:      validPeriodicTaskRecord(),
		LockedUntil: time.Date(2026, time.June, 15, 10, 1, 0, 0, time.UTC),
	}

	if err := due.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestDuePeriodicTaskValidateRequiresLockedUntil(t *testing.T) {
	due := DuePeriodicTask{
		Record:    validPeriodicTaskRecord(),
		LockToken: "token",
	}

	if err := due.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestMarkPeriodicTaskDispatchedRequestValidateAcceptsCompleteRequest(t *testing.T) {
	request := validMarkPeriodicTaskDispatchedRequest()

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestMarkPeriodicTaskDispatchedRequestValidateRequiresName(t *testing.T) {
	request := validMarkPeriodicTaskDispatchedRequest()
	request.Name = ""

	if err := request.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestMarkPeriodicTaskDispatchedRequestValidateRequiresLockToken(t *testing.T) {
	request := validMarkPeriodicTaskDispatchedRequest()
	request.LockToken = ""

	if err := request.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestMarkPeriodicTaskDispatchedRequestValidateRequiresTaskID(t *testing.T) {
	request := validMarkPeriodicTaskDispatchedRequest()
	request.DispatchedTaskID = ""

	if err := request.Validate(); !errors.Is(err, task.ErrInvalidTaskID) {
		t.Fatalf("Validate error = %v, want ErrInvalidTaskID", err)
	}
}

func TestMarkPeriodicTaskDispatchedRequestValidateRequiresDispatchedAt(t *testing.T) {
	request := validMarkPeriodicTaskDispatchedRequest()
	request.DispatchedAt = time.Time{}

	if err := request.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func TestMarkPeriodicTaskDispatchedRequestValidateRequiresNextDueAt(t *testing.T) {
	request := validMarkPeriodicTaskDispatchedRequest()
	request.NextDueAt = time.Time{}

	if err := request.Validate(); !errors.Is(err, ErrInvalidBackendRequest) {
		t.Fatalf("Validate error = %v, want ErrInvalidBackendRequest", err)
	}
}

func validMarkPeriodicTaskDispatchedRequest() MarkPeriodicTaskDispatchedRequest {
	return MarkPeriodicTaskDispatchedRequest{
		Name:             "welcome-email",
		LockToken:        "token",
		DispatchedTaskID: "11111111-1111-4111-8111-111111111111",
		DispatchedAt:     time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC),
		NextDueAt:        time.Date(2026, time.June, 15, 10, 10, 0, 0, time.UTC),
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
