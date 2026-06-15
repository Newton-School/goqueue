package producer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

var errBackend = errors.New("backend failure")

type fakeBackend struct {
	setStateRequests     []backend.TaskStateRecord
	enqueueReadyRequests []backend.EnqueueRequest
	enqueueScheduledReq  []backend.EnqueueRequest
	setStateErr          error
	enqueueReadyErr      error
	enqueueScheduledErr  error
	getStateResult       backend.TaskStateRecord
	getResultResult      backend.TaskResultRecord
	forgetResultErr      error
	getStateErr          error
	getResultErr         error
	queueStatsErr        error
	ensureGroupErr       error
	ackErr               error
	readErr              error
	closeErr             error
	pingErr              error
	closeCalled          bool
	ackCalled            bool
	readCalled           bool
	moveScheduledCalled  bool
	ensureGroupCalled    bool
	queueStatsCalled     bool
	statsRequest         backend.QueueStatsRequest
}

func (f *fakeBackend) EnqueueReady(_ context.Context, req backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	f.enqueueReadyRequests = append(f.enqueueReadyRequests, req)
	return backend.EnqueueResponse{TaskID: task.TaskID(req.Message.ID)}, f.enqueueReadyErr
}

func (f *fakeBackend) EnqueueScheduled(_ context.Context, req backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	f.enqueueScheduledReq = append(f.enqueueScheduledReq, req)
	return backend.EnqueueResponse{TaskID: task.TaskID(req.Message.ID), Scheduled: true}, f.enqueueScheduledErr
}

func (f *fakeBackend) MoveDueScheduled(_ context.Context, _ backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error) {
	f.moveScheduledCalled = true
	return nil, nil
}

func (f *fakeBackend) EnsureConsumerGroup(_ context.Context, _ backend.ConsumerGroupRequest) error {
	f.ensureGroupCalled = true
	return f.ensureGroupErr
}

func (f *fakeBackend) ReadReady(_ context.Context, _ backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
	f.readCalled = true
	return nil, f.readErr
}

func (f *fakeBackend) ClaimStaleReady(_ context.Context, _ backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
	return nil, nil
}

func (f *fakeBackend) Ack(_ context.Context, _ backend.AckRequest) error {
	f.ackCalled = true
	return f.ackErr
}

func (f *fakeBackend) EnqueueDeadLetter(_ context.Context, _ backend.DeadLetterRequest) (backend.DeadLetterRecord, error) {
	return backend.DeadLetterRecord{}, nil
}

func (f *fakeBackend) ReadDeadLetters(_ context.Context, _ backend.ReadDeadLettersRequest) ([]backend.DeadLetterRecord, error) {
	return nil, nil
}

func (f *fakeBackend) UpsertPeriodicTask(_ context.Context, _ backend.UpsertPeriodicTaskRequest) error {
	return nil
}

func (f *fakeBackend) DeletePeriodicTask(_ context.Context, _ backend.DeletePeriodicTaskRequest) error {
	return nil
}

func (f *fakeBackend) ListDuePeriodicTasks(_ context.Context, _ backend.ListDuePeriodicTasksRequest) ([]backend.DuePeriodicTask, error) {
	return nil, nil
}

func (f *fakeBackend) MarkPeriodicTaskDispatched(_ context.Context, _ backend.MarkPeriodicTaskDispatchedRequest) error {
	return nil
}

func (f *fakeBackend) SaveWorkflowChain(_ context.Context, _ backend.WorkflowChainRecord) error {
	return nil
}

func (f *fakeBackend) AdvanceWorkflowChain(_ context.Context, _ backend.AdvanceWorkflowChainRequest) (backend.AdvanceWorkflowChainResponse, error) {
	return backend.AdvanceWorkflowChainResponse{}, nil
}

func (f *fakeBackend) SaveWorkflowGroup(_ context.Context, _ backend.WorkflowGroupRecord) error {
	return nil
}

func (f *fakeBackend) RecordWorkflowTaskCompleted(_ context.Context, _ backend.RecordWorkflowTaskCompletedRequest) (backend.WorkflowGroupProgress, error) {
	return backend.WorkflowGroupProgress{}, nil
}

func (f *fakeBackend) SetTaskState(_ context.Context, record backend.TaskStateRecord) error {
	f.setStateRequests = append(f.setStateRequests, record)
	return f.setStateErr
}

func (f *fakeBackend) GetTaskState(_ context.Context, _ task.TaskID) (backend.TaskStateRecord, error) {
	return f.getStateResult, f.getStateErr
}

func (f *fakeBackend) SaveTaskResult(_ context.Context, _ backend.TaskResultRecord) error {
	return nil
}

func (f *fakeBackend) GetTaskResult(_ context.Context, _ task.TaskID) (backend.TaskResultRecord, error) {
	return f.getResultResult, f.getResultErr
}

func (f *fakeBackend) ForgetTaskResult(_ context.Context, _ task.TaskID) error {
	return f.forgetResultErr
}

func (f *fakeBackend) QueueStats(_ context.Context, req backend.QueueStatsRequest) (backend.QueueStats, error) {
	f.queueStatsCalled = true
	f.statsRequest = req
	return backend.QueueStats{Queue: req.Queue}, f.queueStatsErr
}

func (f *fakeBackend) Ping(_ context.Context) error {
	return f.pingErr
}

func (f *fakeBackend) Close() error {
	f.closeCalled = true
	return f.closeErr
}

func TestNewProducerRequiresBackend(t *testing.T) {
	if _, err := NewProducer(nil); err != ErrNilBackend {
		t.Fatalf("NewProducer error = %v, want ErrNilBackend", err)
	}
}

func TestProducerApplyAsyncEnqueuesReadyTask(t *testing.T) {
	backend := &fakeBackend{}
	producer, err := NewProducer(backend)
	if err != nil {
		t.Fatalf("NewProducer returned error: %v", err)
	}

	result, err := producer.ApplyAsync(context.Background(), "email.send", []any{"welcome"}, map[string]any{"u": 123})
	if err != nil {
		t.Fatalf("ApplyAsync returned error: %v", err)
	}

	if result == nil || result.ID() == "" {
		t.Fatal("ApplyAsync returned empty async result id")
	}
	if len(backend.enqueueReadyRequests) != 1 {
		t.Fatalf("ready calls = %d, want 1", len(backend.enqueueReadyRequests))
	}
	if len(backend.setStateRequests) != 1 {
		t.Fatalf("set state calls = %d, want 1", len(backend.setStateRequests))
	}
	if backend.setStateRequests[0].State != task.TaskPending {
		t.Fatalf("initial state = %q, want %q", backend.setStateRequests[0].State, task.TaskPending)
	}
}

func TestProducerApplyAsyncEnqueuesScheduledTask(t *testing.T) {
	backend := &fakeBackend{}
	fixedNow := func() time.Time { return time.Date(2026, time.June, 14, 10, 0, 0, 0, time.UTC) }
	producer, err := NewProducer(
		backend,
		WithProducerNow(fixedNow),
	)
	if err != nil {
		t.Fatalf("NewProducer returned error: %v", err)
	}

	_, err = producer.ApplyAsync(
		context.Background(),
		"email.send",
		nil,
		nil,
		WithApplyCountDown(10*time.Second),
	)
	if err != nil {
		t.Fatalf("ApplyAsync returned error: %v", err)
	}

	if len(backend.enqueueScheduledReq) != 1 {
		t.Fatalf("scheduled calls = %d, want 1", len(backend.enqueueScheduledReq))
	}
	if len(backend.setStateRequests) != 1 {
		t.Fatalf("set state calls = %d, want 1", len(backend.setStateRequests))
	}
	if backend.setStateRequests[0].State != task.TaskScheduled {
		t.Fatalf("initial state = %q, want %q", backend.setStateRequests[0].State, task.TaskScheduled)
	}
	if backend.enqueueScheduledReq[0].Message.Timing.ETA != fixedNow().Add(10*time.Second) {
		t.Fatalf("scheduled eta = %v, want %v", backend.enqueueScheduledReq[0].Message.Timing.ETA, fixedNow().Add(10*time.Second))
	}
}

func TestProducerApplyAsyncCanOverrideQueue(t *testing.T) {
	backend := &fakeBackend{}
	producer, err := NewProducer(backend)
	if err != nil {
		t.Fatalf("NewProducer returned error: %v", err)
	}

	_, err = producer.ApplyAsync(context.Background(), "email.send", nil, nil, WithApplyQueue("billing"))
	if err != nil {
		t.Fatalf("ApplyAsync returned error: %v", err)
	}

	if len(backend.enqueueReadyRequests) != 1 {
		t.Fatalf("ready calls = %d, want 1", len(backend.enqueueReadyRequests))
	}
	if got := backend.enqueueReadyRequests[0].Message.Queue; got != "billing" {
		t.Fatalf("queued task queue = %q, want billing", got)
	}
}

func TestProducerApplyAsyncUsesProvidedTaskID(t *testing.T) {
	backend := &fakeBackend{}
	producer, err := NewProducer(backend)
	if err != nil {
		t.Fatalf("NewProducer returned error: %v", err)
	}

	id := task.TaskID("12345678-1234-4123-8234-1234567890ab")
	result, err := producer.ApplyAsync(context.Background(), "email.send", nil, nil, WithApplyTaskID(id))
	if err != nil {
		t.Fatalf("ApplyAsync returned error: %v", err)
	}
	if result.ID() != id {
		t.Fatalf("result id = %q, want %q", result.ID(), id)
	}
}

func TestProducerApplyAsyncMarksFailureStateOnEnqueueError(t *testing.T) {
	backend := &fakeBackend{enqueueReadyErr: errBackend}
	producer, err := NewProducer(backend)
	if err != nil {
		t.Fatalf("NewProducer error: %v", err)
	}

	_, err = producer.ApplyAsync(context.Background(), "email.send", nil, nil)
	if err == nil {
		t.Fatal("ApplyAsync expected error")
	}
	if len(backend.setStateRequests) != 2 {
		t.Fatalf("set state calls = %d, want 2", len(backend.setStateRequests))
	}
	if backend.setStateRequests[1].State != task.TaskFailed {
		t.Fatalf("failure state = %q, want %q", backend.setStateRequests[1].State, task.TaskFailed)
	}
}

func TestAsyncResultQueriesBackendStateAndResult(t *testing.T) {
	backend := &fakeBackend{getStateResult: backend.TaskStateRecord{TaskID: task.TaskID("123e4567-e89b-42d3-a456-556642440000"), State: task.TaskSucceeded}}
	result := &AsyncResult{taskID: backend.getStateResult.TaskID, backend: backend}

	state, err := result.TaskState(context.Background())
	if err != nil {
		t.Fatalf("TaskState returned error: %v", err)
	}
	if state.TaskID != backend.getStateResult.TaskID {
		t.Fatalf("TaskState task id = %q, want %q", state.TaskID, backend.getStateResult.TaskID)
	}
}

func TestAsyncResultForgetRequiresTaskID(t *testing.T) {
	result := &AsyncResult{backend: &fakeBackend{}}
	if err := result.ForgetTaskResult(context.Background()); err == nil {
		t.Fatal("ForgetTaskResult expected error")
	}
}
