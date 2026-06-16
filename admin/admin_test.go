package admin

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

const (
	validTaskID = task.TaskID("123e4567-e89b-42d3-a456-556642440111")
)

var _ backend.QueueBackend = (*fakeControlBackend)(nil)

type fakeControlBackend struct {
	getTaskMessageFn    func(context.Context, task.TaskID) (task.TaskMessage, error)
	enqueueReadyFn      func(context.Context, backend.EnqueueRequest) (backend.EnqueueResponse, error)
	enqueueScheduledFn  func(context.Context, backend.EnqueueRequest) (backend.EnqueueResponse, error)
	setTaskStateFn      func(context.Context, backend.TaskStateRecord) error
	forgetTaskResultFn  func(context.Context, task.TaskID) error
	readDeadLetterFn    func(context.Context, task.QueueName, string) (backend.DeadLetterRecord, error)
	deleteDeadLettersFn func(context.Context, task.QueueName, ...string) (int64, error)
	purgeQueueFn        func(context.Context, backend.PurgeQueueRequest) (backend.PurgeQueueResult, error)

	getTaskMessageCalled   bool
	setTaskStateCalled     bool
	forgetTaskResultCalled bool
	enqueueReadyCalled     bool
	enqueueScheduledCalled bool
	readDeadLetterCalled   bool
	deleteDeadLettersArgs  struct {
		queue task.QueueName
		ids   []string
	}
	purgeQueueRequest backend.PurgeQueueRequest
}

func (f *fakeControlBackend) EnqueueReady(ctx context.Context, request backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	f.enqueueReadyCalled = true
	if f.enqueueReadyFn != nil {
		return f.enqueueReadyFn(ctx, request)
	}
	return backend.EnqueueResponse{TaskID: task.TaskID(request.Message.ID), StreamID: "stream-id"}, nil
}

func (f *fakeControlBackend) EnqueueScheduled(ctx context.Context, request backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	f.enqueueScheduledCalled = true
	if f.enqueueScheduledFn != nil {
		return f.enqueueScheduledFn(ctx, request)
	}
	return backend.EnqueueResponse{TaskID: task.TaskID(request.Message.ID), StreamID: "stream-id", Scheduled: true}, nil
}

func (f *fakeControlBackend) MoveDueScheduled(ctx context.Context, request backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error) {
	return nil, nil
}

func (f *fakeControlBackend) EnsureConsumerGroup(ctx context.Context, request backend.ConsumerGroupRequest) error {
	return nil
}

func (f *fakeControlBackend) ReadReady(ctx context.Context, request backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
	return nil, nil
}

func (f *fakeControlBackend) ClaimStaleReady(ctx context.Context, request backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
	return nil, nil
}

func (f *fakeControlBackend) Ack(ctx context.Context, request backend.AckRequest) error { return nil }

func (f *fakeControlBackend) EnqueueDeadLetter(ctx context.Context, request backend.DeadLetterRequest) (backend.DeadLetterRecord, error) {
	return backend.DeadLetterRecord{}, nil
}

func (f *fakeControlBackend) ReadDeadLetters(ctx context.Context, request backend.ReadDeadLettersRequest) ([]backend.DeadLetterRecord, error) {
	return nil, nil
}

func (f *fakeControlBackend) UpsertPeriodicTask(ctx context.Context, request backend.UpsertPeriodicTaskRequest) error {
	return nil
}

func (f *fakeControlBackend) DeletePeriodicTask(ctx context.Context, request backend.DeletePeriodicTaskRequest) error {
	return nil
}

func (f *fakeControlBackend) ListDuePeriodicTasks(ctx context.Context, request backend.ListDuePeriodicTasksRequest) ([]backend.DuePeriodicTask, error) {
	return nil, nil
}

func (f *fakeControlBackend) MarkPeriodicTaskDispatched(ctx context.Context, request backend.MarkPeriodicTaskDispatchedRequest) error {
	return nil
}

func (f *fakeControlBackend) SaveWorkflowChain(ctx context.Context, request backend.WorkflowChainRecord) error {
	return nil
}

func (f *fakeControlBackend) AdvanceWorkflowChain(ctx context.Context, request backend.AdvanceWorkflowChainRequest) (backend.AdvanceWorkflowChainResponse, error) {
	return backend.AdvanceWorkflowChainResponse{}, nil
}

func (f *fakeControlBackend) SaveWorkflowGroup(ctx context.Context, request backend.WorkflowGroupRecord) error {
	return nil
}

func (f *fakeControlBackend) RecordWorkflowTaskCompleted(ctx context.Context, request backend.RecordWorkflowTaskCompletedRequest) (backend.WorkflowGroupProgress, error) {
	return backend.WorkflowGroupProgress{}, nil
}

func (f *fakeControlBackend) SetTaskState(ctx context.Context, request backend.TaskStateRecord) error {
	f.setTaskStateCalled = true
	if f.setTaskStateFn != nil {
		return f.setTaskStateFn(ctx, request)
	}
	return nil
}

func (f *fakeControlBackend) GetTaskState(ctx context.Context, taskID task.TaskID) (backend.TaskStateRecord, error) {
	return backend.TaskStateRecord{}, nil
}

func (f *fakeControlBackend) SaveTaskResult(ctx context.Context, request backend.TaskResultRecord) error {
	return nil
}

func (f *fakeControlBackend) GetTaskResult(ctx context.Context, taskID task.TaskID) (backend.TaskResultRecord, error) {
	return backend.TaskResultRecord{}, nil
}

func (f *fakeControlBackend) ForgetTaskResult(ctx context.Context, taskID task.TaskID) error {
	f.forgetTaskResultCalled = true
	if f.forgetTaskResultFn != nil {
		return f.forgetTaskResultFn(ctx, taskID)
	}
	return nil
}

func (f *fakeControlBackend) QueueStats(ctx context.Context, request backend.QueueStatsRequest) (backend.QueueStats, error) {
	return backend.QueueStats{}, nil
}

func (f *fakeControlBackend) Ping(ctx context.Context) error { return nil }
func (f *fakeControlBackend) Close() error                   { return nil }

func (f *fakeControlBackend) GetTaskMessage(ctx context.Context, taskID task.TaskID) (task.TaskMessage, error) {
	f.getTaskMessageCalled = true
	if f.getTaskMessageFn != nil {
		return f.getTaskMessageFn(ctx, taskID)
	}
	return task.TaskMessage{}, nil
}

func (f *fakeControlBackend) ReadDeadLetter(ctx context.Context, queue task.QueueName, streamID string) (backend.DeadLetterRecord, error) {
	f.readDeadLetterCalled = true
	if f.readDeadLetterFn != nil {
		return f.readDeadLetterFn(ctx, queue, streamID)
	}
	return backend.DeadLetterRecord{}, nil
}

func (f *fakeControlBackend) DeleteDeadLetters(ctx context.Context, queue task.QueueName, streamIDs ...string) (int64, error) {
	f.deleteDeadLettersArgs = struct {
		queue task.QueueName
		ids   []string
	}{queue: queue, ids: append([]string{}, streamIDs...)}
	if f.deleteDeadLettersFn != nil {
		return f.deleteDeadLettersFn(ctx, queue, streamIDs...)
	}
	return int64(len(streamIDs)), nil
}

func (f *fakeControlBackend) PurgeQueue(ctx context.Context, request backend.PurgeQueueRequest) (backend.PurgeQueueResult, error) {
	f.purgeQueueRequest = request
	if f.purgeQueueFn != nil {
		return f.purgeQueueFn(ctx, request)
	}
	return backend.PurgeQueueResult{Queue: request.Queue, ReadyStream: 1}, nil
}

type fakeQueueOnlyBackend struct {
}

func (f *fakeQueueOnlyBackend) EnqueueReady(ctx context.Context, request backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	return backend.EnqueueResponse{}, nil
}

func (f *fakeQueueOnlyBackend) EnqueueScheduled(ctx context.Context, request backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	return backend.EnqueueResponse{}, nil
}

func (f *fakeQueueOnlyBackend) MoveDueScheduled(ctx context.Context, request backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error) {
	return nil, nil
}

func (f *fakeQueueOnlyBackend) EnsureConsumerGroup(ctx context.Context, request backend.ConsumerGroupRequest) error {
	return nil
}

func (f *fakeQueueOnlyBackend) ReadReady(ctx context.Context, request backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
	return nil, nil
}

func (f *fakeQueueOnlyBackend) ClaimStaleReady(ctx context.Context, request backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
	return nil, nil
}

func (f *fakeQueueOnlyBackend) Ack(ctx context.Context, request backend.AckRequest) error { return nil }

func (f *fakeQueueOnlyBackend) EnqueueDeadLetter(ctx context.Context, request backend.DeadLetterRequest) (backend.DeadLetterRecord, error) {
	return backend.DeadLetterRecord{}, nil
}

func (f *fakeQueueOnlyBackend) ReadDeadLetters(ctx context.Context, request backend.ReadDeadLettersRequest) ([]backend.DeadLetterRecord, error) {
	return nil, nil
}

func (f *fakeQueueOnlyBackend) UpsertPeriodicTask(ctx context.Context, request backend.UpsertPeriodicTaskRequest) error {
	return nil
}
func (f *fakeQueueOnlyBackend) DeletePeriodicTask(ctx context.Context, request backend.DeletePeriodicTaskRequest) error {
	return nil
}
func (f *fakeQueueOnlyBackend) ListDuePeriodicTasks(ctx context.Context, request backend.ListDuePeriodicTasksRequest) ([]backend.DuePeriodicTask, error) {
	return nil, nil
}
func (f *fakeQueueOnlyBackend) MarkPeriodicTaskDispatched(ctx context.Context, request backend.MarkPeriodicTaskDispatchedRequest) error {
	return nil
}
func (f *fakeQueueOnlyBackend) SaveWorkflowChain(ctx context.Context, request backend.WorkflowChainRecord) error {
	return nil
}
func (f *fakeQueueOnlyBackend) AdvanceWorkflowChain(ctx context.Context, request backend.AdvanceWorkflowChainRequest) (backend.AdvanceWorkflowChainResponse, error) {
	return backend.AdvanceWorkflowChainResponse{}, nil
}
func (f *fakeQueueOnlyBackend) SaveWorkflowGroup(ctx context.Context, request backend.WorkflowGroupRecord) error {
	return nil
}
func (f *fakeQueueOnlyBackend) RecordWorkflowTaskCompleted(ctx context.Context, request backend.RecordWorkflowTaskCompletedRequest) (backend.WorkflowGroupProgress, error) {
	return backend.WorkflowGroupProgress{}, nil
}
func (f *fakeQueueOnlyBackend) SetTaskState(ctx context.Context, request backend.TaskStateRecord) error {
	return nil
}
func (f *fakeQueueOnlyBackend) GetTaskState(ctx context.Context, taskID task.TaskID) (backend.TaskStateRecord, error) {
	return backend.TaskStateRecord{}, nil
}
func (f *fakeQueueOnlyBackend) SaveTaskResult(ctx context.Context, request backend.TaskResultRecord) error {
	return nil
}
func (f *fakeQueueOnlyBackend) GetTaskResult(ctx context.Context, taskID task.TaskID) (backend.TaskResultRecord, error) {
	return backend.TaskResultRecord{}, nil
}
func (f *fakeQueueOnlyBackend) ForgetTaskResult(ctx context.Context, taskID task.TaskID) error {
	return nil
}
func (f *fakeQueueOnlyBackend) QueueStats(ctx context.Context, request backend.QueueStatsRequest) (backend.QueueStats, error) {
	return backend.QueueStats{}, nil
}
func (f *fakeQueueOnlyBackend) Ping(ctx context.Context) error { return nil }
func (f *fakeQueueOnlyBackend) Close() error                   { return nil }

func testMessage(taskID task.TaskID) task.TaskMessage {
	return task.TaskMessage{
		ID:          taskID.String(),
		Name:        "email.send",
		Queue:       "default",
		Payload:     []byte(`{"v":[1],"kwargs":{}}`),
		Metadata:    map[string]string{"request_id": "req-1"},
		Timing:      task.TaskTiming{},
		Priority:    task.DefaultPriority,
		RetryPolicy: task.DefaultRetryPolicy(),
		CreatedAt:   time.Unix(1_700_000_000, 0).UTC(),
		Attempt:     2,
	}
}

func TestNewAdminRejectsNilBackend(t *testing.T) {
	_, err := NewAdmin(nil)
	if !errors.Is(err, ErrNilAdmin) {
		t.Fatalf("NewAdmin error = %v, want ErrNilAdmin", err)
	}
}

func TestNewAdminRejectsUnsupportedBackend(t *testing.T) {
	_, err := NewAdmin(&fakeQueueOnlyBackend{})
	if !errors.Is(err, ErrAdminBackend) {
		t.Fatalf("NewAdmin error = %v, want ErrAdminBackend", err)
	}
}

func TestNewAdminAcceptsControlBackend(t *testing.T) {
	backend := &fakeControlBackend{}
	admin, err := NewAdmin(backend)
	if err != nil {
		t.Fatalf("NewAdmin returned error: %v", err)
	}
	if admin == nil {
		t.Fatal("NewAdmin returned nil")
	}
}

func TestRetryTaskEnqueuesReadyMessageWithDefaults(t *testing.T) {
	backend := &fakeControlBackend{
		getTaskMessageFn: func(_ context.Context, id task.TaskID) (task.TaskMessage, error) {
			if id != validTaskID {
				return task.TaskMessage{}, fmt.Errorf("unexpected id")
			}
			return testMessage(id), nil
		},
	}
	admin, err := NewAdmin(backend)
	if err != nil {
		t.Fatalf("NewAdmin returned error: %v", err)
	}

	result, err := admin.RetryTask(context.Background(), validTaskID, RetryTaskOptions{})
	if err != nil {
		t.Fatalf("RetryTask returned error: %v", err)
	}

	if result.TaskID != validTaskID {
		t.Fatalf("result.TaskID = %q, want %q", result.TaskID, validTaskID)
	}
	if !backend.getTaskMessageCalled {
		t.Fatal("GetTaskMessage was not called")
	}
	if !backend.enqueueReadyCalled {
		t.Fatal("EnqueueReady was not called")
	}
}

func TestRetryTaskSupportsQueueOverrideAndCountDown(t *testing.T) {
	retryAt := time.Unix(1_700_000_500, 0).UTC()
	backend := &fakeControlBackend{
		getTaskMessageFn: func(_ context.Context, id task.TaskID) (task.TaskMessage, error) {
			return testMessage(id), nil
		},
		enqueueReadyFn: func(_ context.Context, request backend.EnqueueRequest) (backend.EnqueueResponse, error) {
			if request.Message.Queue != "critical" {
				t.Fatalf("requeued queue = %q, want %q", request.Message.Queue, "critical")
			}
			if request.Message.Attempt != 0 {
				t.Fatalf("requeued attempt = %d, want 0", request.Message.Attempt)
			}
			if !request.Message.Timing.ETA.Equal(retryAt) {
				t.Fatalf("requeued eta = %s, want %s", request.Message.Timing.ETA.Format(time.RFC3339), retryAt.Format(time.RFC3339))
			}
			return backend.EnqueueResponse{TaskID: validTaskID, StreamID: "stream-id"}, nil
		},
	}
	admin, err := NewAdmin(backend)
	if err != nil {
		t.Fatalf("NewAdmin returned error: %v", err)
	}

	result, err := admin.RetryTask(context.Background(), validTaskID, RetryTaskOptions{Queue: "critical", ScheduledAt: retryAt, PreserveAttempt: false})
	if err != nil {
		t.Fatalf("RetryTask returned error: %v", err)
	}
	if result.Queue != "critical" {
		t.Fatalf("result.Queue = %q, want %q", result.Queue, "critical")
	}
	if result.OriginalQueue != "default" {
		t.Fatalf("result.OriginalQueue = %q, want %q", result.OriginalQueue, "default")
	}
	if !result.ScheduledAt.Equal(retryAt) {
		t.Fatalf("result.ScheduledAt = %s, want %s", result.ScheduledAt.Format(time.RFC3339), retryAt.Format(time.RFC3339))
	}
}

func TestRetryTaskReturnsErrorForInvalidOptions(t *testing.T) {
	backend := &fakeControlBackend{}
	admin, err := NewAdmin(backend)
	if err != nil {
		t.Fatalf("NewAdmin returned error: %v", err)
	}

	_, err = admin.RetryTask(context.Background(), validTaskID, RetryTaskOptions{CountDown: -1})
	if err == nil {
		t.Fatal("RetryTask expected error for negative countdown")
	}
}

func TestRevokeTaskWritesRevokedState(t *testing.T) {
	var captured backend.TaskStateRecord
	backend := &fakeControlBackend{
		setTaskStateFn: func(_ context.Context, state backend.TaskStateRecord) error {
			captured = state
			return nil
		},
	}
	admin, err := NewAdmin(backend)
	if err != nil {
		t.Fatalf("NewAdmin returned error: %v", err)
	}

	result, err := admin.RevokeTask(context.Background(), validTaskID, "operator requested")
	if err != nil {
		t.Fatalf("RevokeTask returned error: %v", err)
	}
	if result.TaskID != validTaskID {
		t.Fatalf("result.TaskID = %q, want %q", result.TaskID, validTaskID)
	}
	if result.State != task.TaskRevoked {
		t.Fatalf("result.State = %q, want %q", result.State, task.TaskRevoked)
	}
	if captured.TaskID != validTaskID {
		t.Fatalf("captured.TaskID = %q, want %q", captured.TaskID, validTaskID)
	}
	if captured.State != task.TaskRevoked {
		t.Fatalf("captured.State = %q, want %q", captured.State, task.TaskRevoked)
	}
	if captured.Error != "operator requested" {
		t.Fatalf("captured.Error = %q, want %q", captured.Error, "operator requested")
	}
}

func TestReplayDeadLetterReplaysEntryToSameQueueByDefault(t *testing.T) {
	var enqueued bool
	backend := &fakeControlBackend{
		readDeadLetterFn: func(_ context.Context, queue task.QueueName, streamID string) (backend.DeadLetterRecord, error) {
			if queue != "default" {
				t.Fatalf("queue = %q, want %q", queue, "default")
			}
			if streamID != "1-0" {
				t.Fatalf("streamID = %q, want %q", streamID, "1-0")
			}
			return backend.DeadLetterRecord{
				Message: task.TaskMessage{
					ID:       validTaskID.String(),
					Name:     "email.send",
					Queue:    "default",
					Payload:  []byte(`{"v":[],"kwargs":{}}`),
					Metadata: map[string]string{},
				},
				Reason:         task.FailureRetryExhausted,
				SourceStreamID: "1-0",
				Consumer:       "worker",
				Group:          "q",
				FailedAt:       time.Unix(1_700_000_100, 0).UTC(),
			}, nil
		},
		enqueueReadyFn: func(_ context.Context, request backend.EnqueueRequest) (backend.EnqueueResponse, error) {
			enqueued = true
			if request.Message.Queue != "default" {
				t.Fatalf("requeued queue = %q, want %q", request.Message.Queue, "default")
			}
			if request.Message.Attempt != 0 {
				t.Fatalf("requeued attempt = %d, want 0", request.Message.Attempt)
			}
			return backend.EnqueueResponse{TaskID: validTaskID, StreamID: "x"}, nil
		},
	}
	admin, err := NewAdmin(backend)
	if err != nil {
		t.Fatalf("NewAdmin returned error: %v", err)
	}

	result, err := admin.ReplayDeadLetter(context.Background(), "default", "1-0", ReplayDeadLetterOptions{})
	if err != nil {
		t.Fatalf("ReplayDeadLetter returned error: %v", err)
	}
	if !enqueued {
		t.Fatal("replay did not enqueue dead-letter message")
	}
	if result.Destination != "default" {
		t.Fatalf("result.Destination = %q, want %q", result.Destination, "default")
	}
	if result.StreamID != "1-0" {
		t.Fatalf("result.StreamID = %q, want %q", result.StreamID, "1-0")
	}
}

func TestReplayDeadLetterCanDeleteSourceWhenRequested(t *testing.T) {
	backend := &fakeControlBackend{
		readDeadLetterFn: func(_ context.Context, queue task.QueueName, streamID string) (backend.DeadLetterRecord, error) {
			return backend.DeadLetterRecord{
				Message: task.TaskMessage{
					ID:      validTaskID.String(),
					Name:    "email.send",
					Queue:   "default",
					Payload: []byte(`{"v":[],"kwargs":{}}`),
				},
			}, nil
		},
		deleteDeadLettersFn: func(_ context.Context, _ task.QueueName, ids ...string) (int64, error) {
			if len(ids) != 1 || ids[0] != "1-0" {
				return 0, fmt.Errorf("unexpected delete args")
			}
			return 1, nil
		},
	}
	admin, err := NewAdmin(backend)
	if err != nil {
		t.Fatalf("NewAdmin returned error: %v", err)
	}

	result, err := admin.ReplayDeadLetter(context.Background(), "default", "1-0", ReplayDeadLetterOptions{DeleteSource: true})
	if err != nil {
		t.Fatalf("ReplayDeadLetter returned error: %v", err)
	}
	if !result.SourceDeleted {
		t.Fatal("result.SourceDeleted = false, want true")
	}
}

func TestDeleteDeadLettersRejectsEmptyInput(t *testing.T) {
	admin, err := NewAdmin(&fakeControlBackend{})
	if err != nil {
		t.Fatalf("NewAdmin returned error: %v", err)
	}

	if _, err := admin.DeleteDeadLetters(context.Background(), "default", ""); err == nil {
		t.Fatal("DeleteDeadLetters expected error for empty id")
	}
}

func TestPurgeQueueForwardsRequest(t *testing.T) {
	backend := &fakeControlBackend{}
	admin, err := NewAdmin(backend)
	if err != nil {
		t.Fatalf("NewAdmin returned error: %v", err)
	}

	_, err = admin.PurgeQueue(context.Background(), PurgeQueueOptions{Queue: "critical", DeleteMessage: true, DeleteResult: true})
	if err != nil {
		t.Fatalf("PurgeQueue returned error: %v", err)
	}

	if backend.purgeQueueRequest.Queue != "critical" {
		t.Fatalf("backend.purgeQueueRequest.Queue = %q, want %q", backend.purgeQueueRequest.Queue, "critical")
	}
	if !backend.purgeQueueRequest.DeleteMessages {
		t.Fatal("backend request DeleteMessages = false, want true")
	}
	if !backend.purgeQueueRequest.DeleteResults {
		t.Fatal("backend request DeleteResults = false, want true")
	}
	if backend.purgeQueueRequest.DeleteStates {
		t.Fatal("backend request DeleteStates = true, want false")
	}
}

func TestValidateRetryTaskOptionsRejectsBothScheduleInputs(t *testing.T) {
	if err := validateRetryTaskOptions(RetryTaskOptions{ScheduledAt: time.Now(), CountDown: time.Second}); err == nil {
		t.Fatal("expected error for both scheduled_at and countdown")
	}
}

func TestRetryTaskScheduledWhenTaskMessageWasScheduled(t *testing.T) {
	backend := &fakeControlBackend{
		getTaskMessageFn: func(_ context.Context, taskID task.TaskID) (task.TaskMessage, error) {
			message := testMessage(taskID)
			message.Timing = task.TaskTiming{ETA: time.Unix(1_700_000_250, 0).UTC()}
			return message, nil
		},
	}
	admin, err := NewAdmin(backend)
	if err != nil {
		t.Fatalf("NewAdmin returned error: %v", err)
	}

	options := RetryTaskOptions{CountDown: 0}
	_, err = admin.RetryTask(context.Background(), validTaskID, options)
	if err != nil {
		t.Fatalf("RetryTask returned error: %v", err)
	}
	if !backend.enqueueScheduledCalled {
		t.Fatal("RetryTask expected enqueue scheduled when original timing had ETA")
	}
}
