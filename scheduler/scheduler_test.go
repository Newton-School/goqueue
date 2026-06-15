package scheduler

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestNewSchedulerRequiresBackend(t *testing.T) {
	_, err := NewScheduler(nil)
	if !errors.Is(err, ErrNilBackend) {
		t.Fatalf("NewScheduler error = %v, want ErrNilBackend", err)
	}
}

func TestNewSchedulerAppliesOptions(t *testing.T) {
	now := time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC)

	scheduler, err := NewScheduler(
		&fakeBackend{},
		WithSchedulerIdentity("scheduler-1"),
		WithSchedulerDefaultQueue("critical"),
		WithSchedulerBatchSize(12),
		WithSchedulerLockTTL(time.Minute),
		WithSchedulerPollInterval(2*time.Second),
		WithSchedulerNow(func() time.Time { return now }),
	)
	if err != nil {
		t.Fatalf("NewScheduler returned error: %v", err)
	}

	if scheduler.identity != "scheduler-1" {
		t.Fatalf("identity = %q, want scheduler-1", scheduler.identity)
	}
	if scheduler.defaultQueue != "critical" {
		t.Fatalf("defaultQueue = %q, want critical", scheduler.defaultQueue)
	}
	if scheduler.batchSize != 12 {
		t.Fatalf("batchSize = %d, want 12", scheduler.batchSize)
	}
	if scheduler.lockTTL != time.Minute {
		t.Fatalf("lockTTL = %v, want 1m", scheduler.lockTTL)
	}
	if scheduler.pollInterval != 2*time.Second {
		t.Fatalf("pollInterval = %v, want 2s", scheduler.pollInterval)
	}
}

func TestNewSchedulerGeneratesIdentity(t *testing.T) {
	scheduler, err := NewScheduler(&fakeBackend{})
	if err != nil {
		t.Fatalf("NewScheduler returned error: %v", err)
	}

	if !strings.HasPrefix(scheduler.identity, "scheduler-") {
		t.Fatalf("identity = %q, want scheduler prefix", scheduler.identity)
	}
}

func TestSchedulerRegisterPeriodicTaskUpsertsBackendRecord(t *testing.T) {
	now := time.Date(2026, time.June, 15, 10, 0, 0, 0, time.UTC)
	backend := &fakeBackend{}
	scheduler, err := NewScheduler(
		backend,
		WithSchedulerIdentity("scheduler-1"),
		WithSchedulerDefaultQueue("critical"),
		WithSchedulerNow(func() time.Time { return now }),
	)
	if err != nil {
		t.Fatalf("NewScheduler returned error: %v", err)
	}

	definition := validPeriodicTask()
	definition.Queue = ""
	if err := scheduler.RegisterPeriodicTask(context.Background(), definition); err != nil {
		t.Fatalf("RegisterPeriodicTask returned error: %v", err)
	}

	if len(backend.upsertRequests) != 1 {
		t.Fatalf("upsert calls = %d, want 1", len(backend.upsertRequests))
	}
	record := backend.upsertRequests[0].Record
	if record.Name != definition.Name.String() {
		t.Fatalf("record name = %q, want %q", record.Name, definition.Name)
	}
	if record.Queue != "critical" {
		t.Fatalf("record queue = %q, want critical", record.Queue)
	}
	if !record.NextDueAt.Equal(now.Add(10 * time.Minute)) {
		t.Fatalf("next due = %v, want interval after now", record.NextDueAt)
	}
}

func TestSchedulerDeletePeriodicTaskDeletesBackendRecord(t *testing.T) {
	backend := &fakeBackend{}
	scheduler, err := NewScheduler(backend, WithSchedulerIdentity("scheduler-1"))
	if err != nil {
		t.Fatalf("NewScheduler returned error: %v", err)
	}

	if err := scheduler.DeletePeriodicTask(context.Background(), "welcome-email"); err != nil {
		t.Fatalf("DeletePeriodicTask returned error: %v", err)
	}

	if len(backend.deleteRequests) != 1 {
		t.Fatalf("delete calls = %d, want 1", len(backend.deleteRequests))
	}
	if backend.deleteRequests[0].Name != "welcome-email" {
		t.Fatalf("delete name = %q, want welcome-email", backend.deleteRequests[0].Name)
	}
}

type fakeBackend struct {
	mu             sync.Mutex
	upsertRequests []backend.UpsertPeriodicTaskRequest
	deleteRequests []backend.DeletePeriodicTaskRequest
}

func (f *fakeBackend) EnqueueReady(context.Context, backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	return backend.EnqueueResponse{}, nil
}
func (f *fakeBackend) EnqueueScheduled(context.Context, backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	return backend.EnqueueResponse{}, nil
}
func (f *fakeBackend) MoveDueScheduled(context.Context, backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error) {
	return nil, nil
}
func (f *fakeBackend) EnsureConsumerGroup(context.Context, backend.ConsumerGroupRequest) error {
	return nil
}
func (f *fakeBackend) ReadReady(context.Context, backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
	return nil, nil
}
func (f *fakeBackend) ClaimStaleReady(context.Context, backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
	return nil, nil
}
func (f *fakeBackend) Ack(context.Context, backend.AckRequest) error { return nil }
func (f *fakeBackend) EnqueueDeadLetter(context.Context, backend.DeadLetterRequest) (backend.DeadLetterRecord, error) {
	return backend.DeadLetterRecord{}, nil
}
func (f *fakeBackend) ReadDeadLetters(context.Context, backend.ReadDeadLettersRequest) ([]backend.DeadLetterRecord, error) {
	return nil, nil
}
func (f *fakeBackend) UpsertPeriodicTask(_ context.Context, request backend.UpsertPeriodicTaskRequest) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.upsertRequests = append(f.upsertRequests, request)
	return nil
}
func (f *fakeBackend) DeletePeriodicTask(_ context.Context, request backend.DeletePeriodicTaskRequest) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.deleteRequests = append(f.deleteRequests, request)
	return nil
}
func (f *fakeBackend) ListDuePeriodicTasks(context.Context, backend.ListDuePeriodicTasksRequest) ([]backend.DuePeriodicTask, error) {
	return nil, nil
}
func (f *fakeBackend) MarkPeriodicTaskDispatched(context.Context, backend.MarkPeriodicTaskDispatchedRequest) error {
	return nil
}
func (f *fakeBackend) SetTaskState(context.Context, backend.TaskStateRecord) error {
	return nil
}
func (f *fakeBackend) GetTaskState(context.Context, task.TaskID) (backend.TaskStateRecord, error) {
	return backend.TaskStateRecord{}, nil
}
func (f *fakeBackend) SaveTaskResult(context.Context, backend.TaskResultRecord) error {
	return nil
}
func (f *fakeBackend) GetTaskResult(context.Context, task.TaskID) (backend.TaskResultRecord, error) {
	return backend.TaskResultRecord{}, nil
}
func (f *fakeBackend) ForgetTaskResult(context.Context, task.TaskID) error { return nil }
func (f *fakeBackend) QueueStats(context.Context, backend.QueueStatsRequest) (backend.QueueStats, error) {
	return backend.QueueStats{}, nil
}
func (f *fakeBackend) Ping(context.Context) error { return nil }
func (f *fakeBackend) Close() error               { return nil }
