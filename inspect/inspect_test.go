package inspect

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

var errInspectable = errors.New("inspect backend")

type fakeBackend struct {
	taskStateErr    error
	taskResultErr   error
	forgetResultErr error
	pingErr         error
	deadLetterReq   backend.ReadDeadLettersRequest
	deadLetterResp  []backend.DeadLetterRecord
	deadLetterErr   error
	statsReq        backend.QueueStatsRequest
	statsResp       backend.QueueStats
	statsErr        error
	stateResult     backend.TaskStateRecord
	resultResult    backend.TaskResultRecord
}

func (f *fakeBackend) EnqueueReady(_ context.Context, _ backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	return backend.EnqueueResponse{}, nil
}

func (f *fakeBackend) EnqueueScheduled(_ context.Context, _ backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	return backend.EnqueueResponse{}, nil
}

func (f *fakeBackend) MoveDueScheduled(_ context.Context, _ backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error) {
	return nil, nil
}

func (f *fakeBackend) EnsureConsumerGroup(_ context.Context, _ backend.ConsumerGroupRequest) error {
	return nil
}

func (f *fakeBackend) ReadReady(_ context.Context, _ backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
	return nil, nil
}

func (f *fakeBackend) ClaimStaleReady(_ context.Context, _ backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
	return nil, nil
}

func (f *fakeBackend) Ack(_ context.Context, _ backend.AckRequest) error { return nil }
func (f *fakeBackend) EnqueueDeadLetter(_ context.Context, _ backend.DeadLetterRequest) (backend.DeadLetterRecord, error) {
	return backend.DeadLetterRecord{}, nil
}

func (f *fakeBackend) ReadDeadLetters(_ context.Context, req backend.ReadDeadLettersRequest) ([]backend.DeadLetterRecord, error) {
	f.deadLetterReq = req
	return f.deadLetterResp, f.deadLetterErr
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

func (f *fakeBackend) SetTaskState(_ context.Context, _ backend.TaskStateRecord) error { return nil }
func (f *fakeBackend) GetTaskState(_ context.Context, taskID task.TaskID) (backend.TaskStateRecord, error) {
	if taskID == "" {
		return backend.TaskStateRecord{}, errInspectable
	}
	if f.taskStateErr != nil {
		return backend.TaskStateRecord{}, f.taskStateErr
	}
	return f.stateResult, nil
}

func (f *fakeBackend) SaveTaskResult(_ context.Context, _ backend.TaskResultRecord) error { return nil }
func (f *fakeBackend) GetTaskResult(_ context.Context, taskID task.TaskID) (backend.TaskResultRecord, error) {
	if taskID == "" {
		return backend.TaskResultRecord{}, errInspectable
	}
	if f.taskResultErr != nil {
		return backend.TaskResultRecord{}, f.taskResultErr
	}
	return f.resultResult, nil
}

func (f *fakeBackend) ForgetTaskResult(_ context.Context, _ task.TaskID) error {
	return f.forgetResultErr
}
func (f *fakeBackend) QueueStats(_ context.Context, req backend.QueueStatsRequest) (backend.QueueStats, error) {
	f.statsReq = req
	if f.statsErr != nil {
		return backend.QueueStats{}, f.statsErr
	}
	return f.statsResp, nil
}
func (f *fakeBackend) Ping(_ context.Context) error { return f.pingErr }
func (f *fakeBackend) Close() error                 { return nil }

func TestNewInspectorRequiresBackend(t *testing.T) {
	_, err := NewInspector(nil)
	if err == nil {
		t.Fatal("expected error from nil backend")
	}
}

func TestInspectorTaskStateReturnsBackendValue(t *testing.T) {
	backend := &fakeBackend{
		stateResult: backend.TaskStateRecord{
			TaskID:    task.TaskID("123e4567-e89b-42d3-a456-556642440000"),
			State:     task.TaskSucceeded,
			UpdatedAt: time.Date(2026, 6, 16, 10, 0, 0, 0, time.UTC),
		},
	}
	inspector, err := NewInspector(backend)
	if err != nil {
		t.Fatalf("NewInspector returned error: %v", err)
	}

	state, err := inspector.TaskState(context.Background(), backend.stateResult.TaskID)
	if err != nil {
		t.Fatalf("TaskState returned error: %v", err)
	}
	if state.State != task.TaskSucceeded {
		t.Fatalf("state = %q, want %q", state.State, task.TaskSucceeded)
	}
}

func TestInspectorTaskResultReturnsBackendValue(t *testing.T) {
	taskID := task.TaskID("123e4567-e89b-42d3-a456-556642440000")
	backend := &fakeBackend{
		resultResult: backend.TaskResultRecord{
			TaskID: taskID,
			Result: task.TaskResult{
				State: task.TaskSucceeded,
				Value: "ok",
			},
			UpdatedAt: time.Date(2026, 6, 16, 10, 0, 0, 0, time.UTC),
		},
	}
	inspector, err := NewInspector(backend)
	if err != nil {
		t.Fatalf("NewInspector returned error: %v", err)
	}

	result, err := inspector.TaskResult(context.Background(), taskID)
	if err != nil {
		t.Fatalf("TaskResult returned error: %v", err)
	}
	if result.Result.Value != "ok" {
		t.Fatalf("result value = %v, want ok", result.Result.Value)
	}
}

func TestInspectorReadDeadLettersForwardsQueueRequest(t *testing.T) {
	queue := task.QueueName("critical")
	backend := &fakeBackend{
		deadLetterResp: []backend.DeadLetterRecord{
			{
				Message: task.TaskMessage{
					ID:    "123e4567-e89b-42d3-a456-556642440111",
					Name:  "email.send",
					Queue: string(queue),
				},
				FailedAt: time.Date(2026, 6, 16, 10, 0, 0, 0, time.UTC),
			},
		},
	}
	inspector, err := NewInspector(backend)
	if err != nil {
		t.Fatalf("NewInspector returned error: %v", err)
	}

	records, err := inspector.ReadDeadLetters(context.Background(), queue, 2)
	if err != nil {
		t.Fatalf("ReadDeadLetters returned error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("dead letter count = %d, want 1", len(records))
	}
	if backend.deadLetterReq.Queue != queue {
		t.Fatalf("queue = %q, want %q", backend.deadLetterReq.Queue, queue)
	}
}

func TestInspectorQueueStatsForwardsQueueRequest(t *testing.T) {
	queue := task.QueueName("critical")
	backend := &fakeBackend{
		statsResp: backend.QueueStats{
			Queue:           queue,
			ReadyCount:      12,
			ScheduledCount:  3,
			DeadLetterCount: 1,
		},
	}
	inspector, err := NewInspector(backend)
	if err != nil {
		t.Fatalf("NewInspector returned error: %v", err)
	}

	stats, err := inspector.QueueStats(context.Background(), queue)
	if err != nil {
		t.Fatalf("QueueStats returned error: %v", err)
	}
	if stats.ReadyCount != 12 {
		t.Fatalf("ready count = %d, want 12", stats.ReadyCount)
	}
	if backend.statsReq.Queue != queue {
		t.Fatalf("queue = %q, want %q", backend.statsReq.Queue, queue)
	}
}

func TestInspectorSnapshotAggregatesStateAndResult(t *testing.T) {
	taskID := task.TaskID("123e4567-e89b-42d3-a456-556642440555")
	backend := &fakeBackend{
		stateResult: backend.TaskStateRecord{
			TaskID: taskID,
			State:  task.TaskSucceeded,
		},
		resultResult: backend.TaskResultRecord{
			TaskID: taskID,
			Result: task.TaskResult{
				State: task.TaskSucceeded,
				Value: map[string]any{"count": 7},
			},
			UpdatedAt: time.Now().UTC(),
		},
	}
	inspector, err := NewInspector(backend)
	if err != nil {
		t.Fatalf("NewInspector returned error: %v", err)
	}

	snapshot, err := inspector.TaskSnapshot(context.Background(), taskID)
	if err != nil {
		t.Fatalf("TaskSnapshot returned error: %v", err)
	}
	if !snapshot.StateFound {
		t.Fatal("snapshot.StateFound = false, want true")
	}
	if !snapshot.ResultFound {
		t.Fatal("snapshot.ResultFound = false, want true")
	}
	if snapshot.Result.State != task.TaskSucceeded {
		t.Fatalf("snapshot result state = %q, want %q", snapshot.Result.State, task.TaskSucceeded)
	}
}
