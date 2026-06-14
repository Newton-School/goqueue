package worker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

var errTask = errors.New("task failed")

func TestNewWorkerRequiresBackendAndRegistry(t *testing.T) {
	if _, err := NewWorker(nil, task.NewTaskRegistry()); err != ErrNilBackend {
		t.Fatalf("NewWorker error = %v, want ErrNilBackend", err)
	}

	if _, err := NewWorker(&fakeBackend{}, nil); err != ErrNilTaskRegistry {
		t.Fatalf("NewWorker error = %v, want ErrNilTaskRegistry", err)
	}
}

func TestNewWorkerRejectsInvalidOptions(t *testing.T) {
	_, err := NewWorker(&fakeBackend{}, task.NewTaskRegistry(), WithWorkerGroup(""))
	if err == nil {
		t.Fatal("NewWorker expected error")
	}
}

func TestWorkerExecutesSuccessfulTask(t *testing.T) {
	registry := task.NewTaskRegistry()
	if err := registry.Register("email.send", task.TaskHandlerFunc(func(_ task.HandlerContext, _ task.TaskPayload) (task.TaskResult, error) {
		return task.SucceededResult("done"), nil
	})); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	now := time.Date(2026, time.June, 14, 9, 0, 0, 0, time.UTC)
	message := readyMessage(t, task.JSONPayloadCodec{}, testEnvelopeInput{
		name:      "email.send",
		queue:     "billing",
		createdAt: now,
	})
	ackCh := make(chan struct{}, 1)

	backend := &fakeBackend{
		readReadyFn: makeReadOnce(message),
		ackFn: func(_ context.Context, _ backend.AckRequest) error {
			select {
			case ackCh <- struct{}{}:
			default:
			}
			return nil
		},
		moveDueFn: func(_ context.Context, _ backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error) {
			return nil, nil
		},
	}

	worker, err := NewWorker(
		backend,
		registry,
		WithWorkerGroup("workers"),
		WithWorkerConsumer("pod-1"),
		WithWorkerReadBatch(1),
		WithWorkerBlock(0),
		WithWorkerIdleDelay(1*time.Millisecond),
		WithWorkerNow(func() time.Time { return now }),
		WithWorkerMoveDueEnabled(false),
		WithWorkerQueue("billing"),
	)
	if err != nil {
		t.Fatalf("NewWorker returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- worker.Start(ctx)
	}()

	select {
	case <-ackCh:
	case <-time.After(2 * time.Second):
		t.Fatal("task was not acknowledged")
	}
	cancel()

	select {
	case gotErr := <-errCh:
		if gotErr != nil {
			t.Fatalf("Start returned error: %v", gotErr)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not return")
	}

	if len(backend.setStateRequests) < 3 {
		t.Fatalf("set state calls = %d, want at least 3", len(backend.setStateRequests))
	}
	lastState := backend.setStateRequests[len(backend.setStateRequests)-1]
	if lastState.State != task.TaskSucceeded {
		t.Fatalf("final state = %q, want %q", lastState.State, task.TaskSucceeded)
	}
}

func TestWorkerRetriesFailedTask(t *testing.T) {
	registry := task.NewTaskRegistry()
	if err := registry.Register("email.send", task.TaskHandlerFunc(func(_ task.HandlerContext, _ task.TaskPayload) (task.TaskResult, error) {
		return task.FailedResult(errTask), nil
	})); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	now := time.Date(2026, time.June, 14, 9, 0, 0, 0, time.UTC)
	message := readyMessage(t, task.JSONPayloadCodec{}, testEnvelopeInput{
		name:  "email.send",
		queue: "billing",
		retryPolicy: task.RetryPolicy{
			MaxAttempts: 2,
			Backoff:     10 * time.Second,
		},
		createdAt: now,
	})
	ackCh := make(chan struct{}, 1)

	backend := &fakeBackend{
		readReadyFn: makeReadOnce(message),
		ackFn: func(_ context.Context, _ backend.AckRequest) error {
			select {
			case ackCh <- struct{}{}:
			default:
			}
			return nil
		},
		enqueueScheduledFn: func(_ context.Context, req backend.EnqueueRequest) (backend.EnqueueResponse, error) {
			if req.Message.Attempt != 1 {
				t.Errorf("retry attempt = %d, want 1", req.Message.Attempt)
			}
			if req.Message.Timing.ETA != now.Add(10*time.Second) {
				t.Errorf("retry eta = %v, want %v", req.Message.Timing.ETA, now.Add(10*time.Second))
			}
			return backend.EnqueueResponse{}, nil
		},
		moveDueFn: func(_ context.Context, _ backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error) {
			return nil, nil
		},
	}

	worker, err := NewWorker(
		backend,
		registry,
		WithWorkerGroup("workers"),
		WithWorkerConsumer("pod-1"),
		WithWorkerReadBatch(1),
		WithWorkerBlock(0),
		WithWorkerIdleDelay(1*time.Millisecond),
		WithWorkerNow(func() time.Time { return now }),
		WithWorkerMoveDueEnabled(false),
		WithWorkerQueue("billing"),
	)
	if err != nil {
		t.Fatalf("NewWorker returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- worker.Start(ctx)
	}()

	select {
	case <-ackCh:
	case <-time.After(2 * time.Second):
		t.Fatal("task was not acknowledged")
	}
	cancel()

	select {
	case gotErr := <-errCh:
		if gotErr != nil {
			t.Fatalf("Start returned error: %v", gotErr)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not return")
	}

	if len(backend.enqueueScheduledRequests) != 1 {
		t.Fatalf("scheduled calls = %d, want 1", len(backend.enqueueScheduledRequests))
	}
	lastState := backend.setStateRequests[len(backend.setStateRequests)-1]
	if lastState.State != task.TaskRetrying {
		t.Fatalf("final state = %q, want %q", lastState.State, task.TaskRetrying)
	}
}

func TestWorkerDoesNotRetryWhenMaxAttemptsReached(t *testing.T) {
	registry := task.NewTaskRegistry()
	if err := registry.Register("email.send", task.TaskHandlerFunc(func(_ task.HandlerContext, _ task.TaskPayload) (task.TaskResult, error) {
		return task.FailedResult(errTask), nil
	})); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	now := time.Date(2026, time.June, 14, 9, 0, 0, 0, time.UTC)
	message := readyMessage(t, task.JSONPayloadCodec{}, testEnvelopeInput{
		name:  "email.send",
		queue: "billing",
		retryPolicy: task.RetryPolicy{
			MaxAttempts: 1,
			Backoff:     10 * time.Second,
		},
		createdAt: now,
	})
	ackCh := make(chan struct{}, 1)

	backend := &fakeBackend{
		readReadyFn: makeReadOnce(message),
		ackFn: func(_ context.Context, _ backend.AckRequest) error {
			select {
			case ackCh <- struct{}{}:
			default:
			}
			return nil
		},
		moveDueFn: func(_ context.Context, _ backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error) {
			return nil, nil
		},
	}

	worker, err := NewWorker(
		backend,
		registry,
		WithWorkerGroup("workers"),
		WithWorkerConsumer("pod-1"),
		WithWorkerReadBatch(1),
		WithWorkerBlock(0),
		WithWorkerIdleDelay(1*time.Millisecond),
		WithWorkerNow(func() time.Time { return now }),
		WithWorkerMoveDueEnabled(false),
		WithWorkerQueue("billing"),
	)
	if err != nil {
		t.Fatalf("NewWorker returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- worker.Start(ctx)
	}()

	select {
	case <-ackCh:
	case <-time.After(2 * time.Second):
		t.Fatal("task was not acknowledged")
	}
	cancel()

	select {
	case gotErr := <-errCh:
		if gotErr != nil {
			t.Fatalf("Start returned error: %v", gotErr)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not return")
	}

	if len(backend.enqueueScheduledRequests) != 0 {
		t.Fatalf("scheduled calls = %d, want 0", len(backend.enqueueScheduledRequests))
	}
	lastState := backend.setStateRequests[len(backend.setStateRequests)-1]
	if lastState.State != task.TaskFailed {
		t.Fatalf("final state = %q, want %q", lastState.State, task.TaskFailed)
	}
}

func TestWorkerSkipsExpiredTask(t *testing.T) {
	registry := task.NewTaskRegistry()
	if err := registry.Register("email.send", task.TaskHandlerFunc(func(_ task.HandlerContext, _ task.TaskPayload) (task.TaskResult, error) {
		t.Fatal("expired task handler should not execute")
		return task.SucceededResult("done"), nil
	})); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	now := time.Date(2026, time.June, 14, 9, 0, 0, 0, time.UTC)
	message := readyMessage(t, task.JSONPayloadCodec{}, testEnvelopeInput{
		name:      "email.send",
		queue:     "billing",
		expiresAt: now.Add(-time.Minute),
	})
	ackCh := make(chan struct{}, 1)

	backend := &fakeBackend{
		readReadyFn: makeReadOnce(message),
		ackFn: func(_ context.Context, _ backend.AckRequest) error {
			select {
			case ackCh <- struct{}{}:
			default:
			}
			return nil
		},
		moveDueFn: func(_ context.Context, _ backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error) {
			return nil, nil
		},
	}

	worker, err := NewWorker(
		backend,
		registry,
		WithWorkerGroup("workers"),
		WithWorkerConsumer("pod-1"),
		WithWorkerReadBatch(1),
		WithWorkerBlock(0),
		WithWorkerIdleDelay(1*time.Millisecond),
		WithWorkerNow(func() time.Time { return now }),
		WithWorkerMoveDueEnabled(false),
		WithWorkerQueue("billing"),
	)
	if err != nil {
		t.Fatalf("NewWorker returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- worker.Start(ctx)
	}()

	select {
	case <-ackCh:
	case <-time.After(2 * time.Second):
		t.Fatal("task was not acknowledged")
	}
	cancel()

	select {
	case gotErr := <-errCh:
		if gotErr != nil {
			t.Fatalf("Start returned error: %v", gotErr)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not return")
	}

	lastState := backend.setStateRequests[len(backend.setStateRequests)-1]
	if lastState.State != task.TaskExpired {
		t.Fatalf("final state = %q, want %q", lastState.State, task.TaskExpired)
	}
}

type fakeBackend struct {
	mu                       sync.Mutex
	ensureGroupRequests      []backend.ConsumerGroupRequest
	moveDueRequests          []backend.MoveDueScheduledRequest
	readReadyCalls           int
	setStateRequests         []backend.TaskStateRecord
	resultRequests           []backend.TaskResultRecord
	ackRequests              []backend.AckRequest
	enqueueScheduledRequests []backend.EnqueueRequest
	ensureConsumerGroupErr   error
	moveDueScheduledErr      error
	readReadyErr             error
	readReadyFn              func(context.Context, backend.ReadReadyRequest) ([]backend.ReadyMessage, error)
	ensureConsumerGroupFn    func(backend.ConsumerGroupRequest) error
	moveDueFn                func(context.Context, backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error)
	ackFn                    func(context.Context, backend.AckRequest) error
	setTaskStateErr          error
	saveTaskResultErr        error
	enqueueScheduledErr      error
	enqueueScheduledFn       func(context.Context, backend.EnqueueRequest) (backend.EnqueueResponse, error)
}

func (f *fakeBackend) EnqueueReady(_ context.Context, _ backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	return backend.EnqueueResponse{}, nil
}

func (f *fakeBackend) EnqueueScheduled(ctx context.Context, req backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	f.mu.Lock()
	f.enqueueScheduledRequests = append(f.enqueueScheduledRequests, req)
	f.mu.Unlock()

	if f.enqueueScheduledFn != nil {
		return f.enqueueScheduledFn(ctx, req)
	}

	if f.enqueueScheduledErr != nil {
		return backend.EnqueueResponse{}, f.enqueueScheduledErr
	}

	return backend.EnqueueResponse{}, nil
}

func (f *fakeBackend) MoveDueScheduled(ctx context.Context, req backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error) {
	f.mu.Lock()
	f.moveDueRequests = append(f.moveDueRequests, req)
	f.mu.Unlock()

	if f.moveDueFn != nil {
		return f.moveDueFn(ctx, req)
	}
	if f.moveDueScheduledErr != nil {
		return nil, f.moveDueScheduledErr
	}
	return nil, nil
}

func (f *fakeBackend) EnsureConsumerGroup(ctx context.Context, req backend.ConsumerGroupRequest) error {
	f.mu.Lock()
	f.ensureGroupRequests = append(f.ensureGroupRequests, req)
	f.mu.Unlock()

	if f.ensureConsumerGroupFn != nil {
		return f.ensureConsumerGroupFn(req)
	}

	return f.ensureConsumerGroupErr
}

func (f *fakeBackend) ReadReady(_ context.Context, req backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
	f.mu.Lock()
	f.readReadyCalls++
	f.mu.Unlock()

	return f.readReadyFn(context.Background(), req)
}

func (f *fakeBackend) Ack(ctx context.Context, req backend.AckRequest) error {
	f.mu.Lock()
	f.ackRequests = append(f.ackRequests, req)
	f.mu.Unlock()

	if f.ackFn != nil {
		return f.ackFn(ctx, req)
	}

	return nil
}

func (f *fakeBackend) SetTaskState(_ context.Context, record backend.TaskStateRecord) error {
	f.mu.Lock()
	f.setStateRequests = append(f.setStateRequests, record)
	f.mu.Unlock()
	return f.setTaskStateErr
}

func (f *fakeBackend) GetTaskState(_ context.Context, _ task.TaskID) (backend.TaskStateRecord, error) {
	return backend.TaskStateRecord{}, nil
}

func (f *fakeBackend) SaveTaskResult(_ context.Context, record backend.TaskResultRecord) error {
	f.mu.Lock()
	f.resultRequests = append(f.resultRequests, record)
	f.mu.Unlock()
	return f.saveTaskResultErr
}

func (f *fakeBackend) GetTaskResult(_ context.Context, _ task.TaskID) (backend.TaskResultRecord, error) {
	return backend.TaskResultRecord{}, nil
}

func (f *fakeBackend) ForgetTaskResult(_ context.Context, _ task.TaskID) error {
	return nil
}

func (f *fakeBackend) QueueStats(_ context.Context, _ backend.QueueStatsRequest) (backend.QueueStats, error) {
	return backend.QueueStats{}, nil
}

func (f *fakeBackend) Ping(_ context.Context) error {
	return nil
}

func (f *fakeBackend) Close() error {
	return nil
}

type testEnvelopeInput struct {
	name        string
	queue       string
	attempt     int
	createdAt   time.Time
	retryPolicy task.RetryPolicy
	expiresAt   time.Time
}

func readyMessage(t *testing.T, codec task.PayloadCodec, input testEnvelopeInput) backend.ReadyMessage {
	t.Helper()

	envelope, err := task.NewTaskEnvelope(task.TaskEnvelopeInput{
		Name:        task.TaskName(input.name),
		Queue:       task.QueueName(input.queue),
		CreatedAt:   input.createdAt,
		Attempt:     input.attempt,
		RetryPolicy: input.retryPolicy,
		Timing: task.TaskTiming{
			ExpiresAt: input.expiresAt,
		},
	})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}

	message, err := task.TaskEnvelopeToMessage(envelope, codec)
	if err != nil {
		t.Fatalf("TaskEnvelopeToMessage returned error: %v", err)
	}

	return backend.ReadyMessage{
		StreamID: "1-0",
		Message:  message,
	}
}

func makeReadOnce(message backend.ReadyMessage) func(context.Context, backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
	called := false

	return func(context.Context, backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
		if called {
			return []backend.ReadyMessage{}, nil
		}
		called = true

		return []backend.ReadyMessage{message}, nil
	}
}
