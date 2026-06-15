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

func TestNewWorkerCopiesReliabilityOptions(t *testing.T) {
	worker, err := NewWorker(
		&fakeBackend{},
		task.NewTaskRegistry(),
		WithWorkerDeadLetterEnabled(false),
		WithWorkerPendingRecoveryEnabled(true),
		WithWorkerPendingMinIdle(3*time.Minute),
		WithWorkerPendingClaimBatch(7),
		WithWorkerPendingClaimInterval(11*time.Second),
	)
	if err != nil {
		t.Fatalf("NewWorker returned error: %v", err)
	}

	if worker.deadLetterEnabled {
		t.Fatal("dead letter should be disabled")
	}
	if !worker.pendingRecoveryEnabled {
		t.Fatal("pending recovery should be enabled")
	}
	if worker.pendingMinIdle != 3*time.Minute {
		t.Fatalf("pending min idle = %v, want 3m", worker.pendingMinIdle)
	}
	if worker.pendingClaimBatch != 7 {
		t.Fatalf("pending claim batch = %d, want 7", worker.pendingClaimBatch)
	}
	if worker.pendingClaimInterval != 11*time.Second {
		t.Fatalf("pending claim interval = %v, want 11s", worker.pendingClaimInterval)
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
	lastResult := backend.resultRequests[len(backend.resultRequests)-1].Result
	if lastResult.Metadata[task.FailureMetadataCategoryKey] != string(task.FailureExecution) {
		t.Fatalf("failure category = %q, want execution", lastResult.Metadata[task.FailureMetadataCategoryKey])
	}
	if lastResult.Metadata[task.FailureMetadataRetryableKey] != "true" {
		t.Fatalf("retryable = %q, want true", lastResult.Metadata[task.FailureMetadataRetryableKey])
	}
	if lastResult.Metadata[task.FailureMetadataNextRetryAtKey] != now.Add(10*time.Second).Format(time.RFC3339Nano) {
		t.Fatalf("next retry = %q, want scheduled timestamp", lastResult.Metadata[task.FailureMetadataNextRetryAtKey])
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
	if len(backend.deadLetterRequests) != 1 {
		t.Fatalf("dead letter requests = %d, want 1", len(backend.deadLetterRequests))
	}
	if backend.deadLetterRequests[0].Reason != task.FailureRetryExhausted {
		t.Fatalf("dead letter reason = %q, want %q", backend.deadLetterRequests[0].Reason, task.FailureRetryExhausted)
	}
	lastState := backend.setStateRequests[len(backend.setStateRequests)-1]
	if lastState.State != task.TaskDeadLettered {
		t.Fatalf("final state = %q, want %q", lastState.State, task.TaskDeadLettered)
	}
	lastResult := backend.resultRequests[len(backend.resultRequests)-1].Result
	if lastResult.Metadata[task.FailureMetadataCategoryKey] != string(task.FailureRetryExhausted) {
		t.Fatalf("failure category = %q, want retry exhausted", lastResult.Metadata[task.FailureMetadataCategoryKey])
	}
	if lastResult.Metadata[task.FailureMetadataDeadLetteredKey] != "true" {
		t.Fatalf("dead lettered = %q, want true", lastResult.Metadata[task.FailureMetadataDeadLetteredKey])
	}
}

func TestWorkerDeadLettersRetryScheduleFailure(t *testing.T) {
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

	backend := &fakeBackend{
		readReadyFn:         makeReadOnce(message),
		enqueueScheduledErr: errTask,
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

	errCh := make(chan error, 1)
	go func() {
		errCh <- worker.Start(context.Background())
	}()

	select {
	case gotErr := <-errCh:
		if gotErr == nil {
			t.Fatal("Start expected retry schedule failure")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not return")
	}

	if len(backend.deadLetterRequests) != 1 {
		t.Fatalf("dead letter requests = %d, want 1", len(backend.deadLetterRequests))
	}
	if backend.deadLetterRequests[0].Reason != task.FailureRetryScheduleFailed {
		t.Fatalf("dead letter reason = %q, want %q", backend.deadLetterRequests[0].Reason, task.FailureRetryScheduleFailed)
	}
	if len(backend.ackRequests) != 0 {
		t.Fatalf("ack requests = %d, want 0", len(backend.ackRequests))
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
	if len(backend.deadLetterRequests) != 1 {
		t.Fatalf("dead letter requests = %d, want 1", len(backend.deadLetterRequests))
	}
	if backend.deadLetterRequests[0].Reason != task.FailureExpired {
		t.Fatalf("dead letter reason = %q, want %q", backend.deadLetterRequests[0].Reason, task.FailureExpired)
	}
	lastResult := backend.resultRequests[len(backend.resultRequests)-1].Result
	if lastResult.Metadata[task.FailureMetadataCategoryKey] != string(task.FailureExpired) {
		t.Fatalf("failure category = %q, want expired", lastResult.Metadata[task.FailureMetadataCategoryKey])
	}
	if lastResult.Metadata[task.FailureMetadataDeadLetteredKey] != "true" {
		t.Fatalf("dead lettered = %q, want true", lastResult.Metadata[task.FailureMetadataDeadLetteredKey])
	}
}

func TestWorkerDeadLettersUnknownTask(t *testing.T) {
	registry := task.NewTaskRegistry()
	now := time.Date(2026, time.June, 14, 9, 0, 0, 0, time.UTC)
	message := readyMessage(t, task.JSONPayloadCodec{}, testEnvelopeInput{
		name:      "email.missing",
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

	if len(backend.deadLetterRequests) != 1 {
		t.Fatalf("dead letter requests = %d, want 1", len(backend.deadLetterRequests))
	}
	if backend.deadLetterRequests[0].Reason != task.FailureUnknownTask {
		t.Fatalf("dead letter reason = %q, want %q", backend.deadLetterRequests[0].Reason, task.FailureUnknownTask)
	}
	lastState := backend.setStateRequests[len(backend.setStateRequests)-1]
	if lastState.State != task.TaskDeadLettered {
		t.Fatalf("final state = %q, want %q", lastState.State, task.TaskDeadLettered)
	}
	lastResult := backend.resultRequests[len(backend.resultRequests)-1].Result
	if lastResult.Metadata[task.FailureMetadataCategoryKey] != string(task.FailureUnknownTask) {
		t.Fatalf("failure category = %q, want unknown task", lastResult.Metadata[task.FailureMetadataCategoryKey])
	}
}

func TestWorkerMarksUnknownTaskFailedWhenDeadLetterDisabled(t *testing.T) {
	registry := task.NewTaskRegistry()
	now := time.Date(2026, time.June, 14, 9, 0, 0, 0, time.UTC)
	message := readyMessage(t, task.JSONPayloadCodec{}, testEnvelopeInput{
		name:      "email.missing",
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
		WithWorkerDeadLetterEnabled(false),
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

	if len(backend.deadLetterRequests) != 0 {
		t.Fatalf("dead letter requests = %d, want 0", len(backend.deadLetterRequests))
	}
	lastState := backend.setStateRequests[len(backend.setStateRequests)-1]
	if lastState.State != task.TaskFailed {
		t.Fatalf("final state = %q, want %q", lastState.State, task.TaskFailed)
	}
}

func TestWorkerDoesNotAckWhenDeadLetterFails(t *testing.T) {
	registry := task.NewTaskRegistry()
	now := time.Date(2026, time.June, 14, 9, 0, 0, 0, time.UTC)
	message := readyMessage(t, task.JSONPayloadCodec{}, testEnvelopeInput{
		name:      "email.missing",
		queue:     "billing",
		createdAt: now,
	})

	backend := &fakeBackend{
		readReadyFn:   makeReadOnce(message),
		deadLetterErr: errTask,
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
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- worker.Start(ctx)
	}()

	select {
	case gotErr := <-errCh:
		if gotErr == nil {
			t.Fatal("Start expected dead letter error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not return")
	}

	if len(backend.ackRequests) != 0 {
		t.Fatalf("ack requests = %d, want 0", len(backend.ackRequests))
	}
}

func TestWorkerDeadLettersMalformedPayload(t *testing.T) {
	registry := task.NewTaskRegistry()
	now := time.Date(2026, time.June, 14, 9, 0, 0, 0, time.UTC)
	message := backend.ReadyMessage{
		StreamID: "1-0",
		Message: task.TaskMessage{
			ID:          "4ac0a01f-1b16-4330-b3e7-e99826eacb1a",
			Name:        "email.send",
			Queue:       "billing",
			Payload:     []byte("{"),
			Priority:    task.DefaultPriority,
			RetryPolicy: task.DefaultRetryPolicy(),
			CreatedAt:   now,
		},
	}
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

	if len(backend.deadLetterRequests) != 1 {
		t.Fatalf("dead letter requests = %d, want 1", len(backend.deadLetterRequests))
	}
	if backend.deadLetterRequests[0].Reason != task.FailureMalformedMessage {
		t.Fatalf("dead letter reason = %q, want %q", backend.deadLetterRequests[0].Reason, task.FailureMalformedMessage)
	}
	lastResult := backend.resultRequests[len(backend.resultRequests)-1].Result
	if lastResult.Metadata[task.FailureMetadataCategoryKey] != string(task.FailureMalformedMessage) {
		t.Fatalf("failure category = %q, want malformed message", lastResult.Metadata[task.FailureMetadataCategoryKey])
	}
}

func TestWorkerProcessesClaimedPendingTask(t *testing.T) {
	registry := task.NewTaskRegistry()
	if err := registry.Register("email.send", task.TaskHandlerFunc(func(_ task.HandlerContext, _ task.TaskPayload) (task.TaskResult, error) {
		return task.SucceededResult("claimed"), nil
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
		claimStaleReadyFn: makeClaimOnce(message),
		readReadyFn: func(context.Context, backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
			return nil, nil
		},
		ackFn: func(_ context.Context, _ backend.AckRequest) error {
			select {
			case ackCh <- struct{}{}:
			default:
			}
			return nil
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
		WithWorkerPendingRecoveryEnabled(true),
		WithWorkerPendingMinIdle(2*time.Minute),
		WithWorkerPendingClaimBatch(3),
		WithWorkerPendingClaimInterval(0),
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
		t.Fatal("claimed task was not acknowledged")
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

	if len(backend.claimStaleReadyRequests) == 0 {
		t.Fatal("expected pending claim request")
	}
	request := backend.claimStaleReadyRequests[0]
	if request.MinIdle != 2*time.Minute {
		t.Fatalf("min idle = %v, want 2m", request.MinIdle)
	}
	if request.Count != 3 {
		t.Fatalf("count = %d, want 3", request.Count)
	}
}

func TestWorkerReturnsPendingRecoveryError(t *testing.T) {
	registry := task.NewTaskRegistry()
	backend := &fakeBackend{
		claimStaleReadyFn: func(context.Context, backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
			return nil, errTask
		},
		readReadyFn: func(context.Context, backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
			t.Fatal("ReadReady should not run after claim error")
			return nil, nil
		},
	}

	worker, err := NewWorker(
		backend,
		registry,
		WithWorkerGroup("workers"),
		WithWorkerConsumer("pod-1"),
		WithWorkerBlock(0),
		WithWorkerMoveDueEnabled(false),
		WithWorkerPendingRecoveryEnabled(true),
		WithWorkerPendingClaimInterval(0),
	)
	if err != nil {
		t.Fatalf("NewWorker returned error: %v", err)
	}

	err = worker.Start(context.Background())
	if err == nil {
		t.Fatal("Start expected pending recovery error")
	}
}

func TestWorkerDoesNotClaimPendingWhenRecoveryDisabled(t *testing.T) {
	registry := task.NewTaskRegistry()
	ctx, cancel := context.WithCancel(context.Background())
	backend := &fakeBackend{
		claimStaleReadyFn: func(context.Context, backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
			t.Fatal("ClaimStaleReady should not run")
			return nil, nil
		},
		readReadyFn: func(context.Context, backend.ReadReadyRequest) ([]backend.ReadyMessage, error) {
			cancel()
			return nil, nil
		},
	}

	worker, err := NewWorker(
		backend,
		registry,
		WithWorkerGroup("workers"),
		WithWorkerConsumer("pod-1"),
		WithWorkerBlock(0),
		WithWorkerIdleDelay(0),
		WithWorkerMoveDueEnabled(false),
	)
	if err != nil {
		t.Fatalf("NewWorker returned error: %v", err)
	}

	if err := worker.Start(ctx); err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
}

func TestWorkerAdvancesChainAfterSuccessfulTask(t *testing.T) {
	registry := task.NewTaskRegistry()
	if err := registry.Register("email.send", task.TaskHandlerFunc(func(_ task.HandlerContext, _ task.TaskPayload) (task.TaskResult, error) {
		return task.SucceededResult("done"), nil
	})); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	now := time.Date(2026, time.June, 15, 9, 0, 0, 0, time.UTC)
	message := readyMessage(t, task.JSONPayloadCodec{}, testEnvelopeInput{
		name:      "email.send",
		queue:     "billing",
		createdAt: now,
		metadata: map[string]string{
			"goqueue.workflow.kind":       "chain",
			"goqueue.workflow.chain_id":   "chain-1",
			"goqueue.workflow.chain_step": "0",
		},
	})
	next := backend.WorkflowSignatureRecord{
		Name:        "email.audit",
		Queue:       "billing",
		Priority:    task.DefaultPriority,
		RetryPolicy: task.DefaultRetryPolicy(),
	}
	ackCh := make(chan struct{}, 1)
	fake := &fakeBackend{
		readReadyFn: makeReadOnce(message),
		advanceChainResponse: backend.AdvanceWorkflowChainResponse{
			Advanced: true,
			Next:     &next,
		},
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
		fake,
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

	if len(fake.advanceChainRequests) != 1 {
		t.Fatalf("advance chain calls = %d, want 1", len(fake.advanceChainRequests))
	}
	if len(fake.enqueueReadyRequests) != 1 {
		t.Fatalf("enqueue ready calls = %d, want 1", len(fake.enqueueReadyRequests))
	}
	if fake.enqueueReadyRequests[0].Message.Name != "email.audit" {
		t.Fatalf("next task name = %q, want email.audit", fake.enqueueReadyRequests[0].Message.Name)
	}
}

func TestWorkerDoesNotAdvanceChainAfterFailedTask(t *testing.T) {
	registry := task.NewTaskRegistry()
	if err := registry.Register("email.send", task.TaskHandlerFunc(func(_ task.HandlerContext, _ task.TaskPayload) (task.TaskResult, error) {
		return task.TaskResult{}, errTask
	})); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	now := time.Date(2026, time.June, 15, 9, 0, 0, 0, time.UTC)
	message := readyMessage(t, task.JSONPayloadCodec{}, testEnvelopeInput{
		name:      "email.send",
		queue:     "billing",
		createdAt: now,
		metadata: map[string]string{
			"goqueue.workflow.kind":       "chain",
			"goqueue.workflow.chain_id":   "chain-1",
			"goqueue.workflow.chain_step": "0",
		},
	})
	ackCh := make(chan struct{}, 1)
	fake := &fakeBackend{
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
		fake,
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

	if len(fake.advanceChainRequests) != 0 {
		t.Fatalf("advance chain calls = %d, want 0", len(fake.advanceChainRequests))
	}
	if len(fake.enqueueReadyRequests) != 0 {
		t.Fatalf("enqueue ready calls = %d, want 0", len(fake.enqueueReadyRequests))
	}
}

func TestWorkerRecordsGroupProgressAfterTerminalTask(t *testing.T) {
	registry := task.NewTaskRegistry()
	if err := registry.Register("email.send", task.TaskHandlerFunc(func(_ task.HandlerContext, _ task.TaskPayload) (task.TaskResult, error) {
		return task.SucceededResult("done"), nil
	})); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	now := time.Date(2026, time.June, 15, 9, 0, 0, 0, time.UTC)
	message := readyMessage(t, task.JSONPayloadCodec{}, testEnvelopeInput{
		name:      "email.send",
		queue:     "billing",
		createdAt: now,
		metadata: map[string]string{
			"goqueue.workflow.kind":     "group",
			"goqueue.workflow.group_id": "group-1",
		},
	})
	ackCh := make(chan struct{}, 1)
	fake := &fakeBackend{
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
		fake,
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

	if len(fake.recordGroupRequests) != 1 {
		t.Fatalf("record group calls = %d, want 1", len(fake.recordGroupRequests))
	}
	if fake.recordGroupRequests[0].GroupID != "group-1" {
		t.Fatalf("group id = %q, want group-1", fake.recordGroupRequests[0].GroupID)
	}
	if fake.recordGroupRequests[0].State != task.TaskSucceeded {
		t.Fatalf("state = %q, want succeeded", fake.recordGroupRequests[0].State)
	}
}

func TestWorkerRecordsFailedGroupProgress(t *testing.T) {
	now := time.Date(2026, time.June, 15, 9, 0, 0, 0, time.UTC)
	envelope, err := task.NewTaskEnvelope(task.TaskEnvelopeInput{
		ID:        "11111111-1111-4111-8111-111111111111",
		Name:      "email.send",
		Queue:     "billing",
		Metadata:  map[string]string{"goqueue.workflow.kind": "group", "goqueue.workflow.group_id": "group-1"},
		CreatedAt: now,
	})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}
	fake := &fakeBackend{}
	worker := &Worker{backend: fake, now: func() time.Time { return now }}

	if err := worker.advanceWorkflow(context.Background(), envelope, task.TaskFailed); err != nil {
		t.Fatalf("advanceWorkflow returned error: %v", err)
	}

	if len(fake.advanceChainRequests) != 0 {
		t.Fatalf("advance chain calls = %d, want 0", len(fake.advanceChainRequests))
	}
	if len(fake.recordGroupRequests) != 1 {
		t.Fatalf("record group calls = %d, want 1", len(fake.recordGroupRequests))
	}
	if fake.recordGroupRequests[0].State != task.TaskFailed {
		t.Fatalf("state = %q, want failed", fake.recordGroupRequests[0].State)
	}
}

func TestWorkerDispatchesChordCallbackWhenGroupCompletes(t *testing.T) {
	registry := task.NewTaskRegistry()
	if err := registry.Register("email.send", task.TaskHandlerFunc(func(_ task.HandlerContext, _ task.TaskPayload) (task.TaskResult, error) {
		return task.SucceededResult("done"), nil
	})); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	now := time.Date(2026, time.June, 15, 9, 0, 0, 0, time.UTC)
	message := readyMessage(t, task.JSONPayloadCodec{}, testEnvelopeInput{
		name:      "email.send",
		queue:     "billing",
		createdAt: now,
		metadata: map[string]string{
			"goqueue.workflow.kind":     "chord",
			"goqueue.workflow.group_id": "group-1",
			"goqueue.workflow.chord_id": "group-1",
		},
	})
	callback := backend.WorkflowSignatureRecord{
		Name:        "email.complete",
		Queue:       "billing",
		Priority:    task.DefaultPriority,
		RetryPolicy: task.DefaultRetryPolicy(),
	}
	ackCh := make(chan struct{}, 1)
	fake := &fakeBackend{
		readReadyFn: makeReadOnce(message),
		recordGroupResponse: backend.WorkflowGroupProgress{
			GroupID:   "group-1",
			Total:     1,
			Completed: 1,
			Done:      true,
			Succeeded: true,
			Callback:  &callback,
		},
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
		fake,
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

	if len(fake.enqueueReadyRequests) != 1 {
		t.Fatalf("enqueue ready calls = %d, want 1", len(fake.enqueueReadyRequests))
	}
	callbackMessage := fake.enqueueReadyRequests[0].Message
	if callbackMessage.Name != "email.complete" {
		t.Fatalf("callback name = %q, want email.complete", callbackMessage.Name)
	}
	if callbackMessage.Metadata["goqueue.workflow.chord_callback"] != "true" {
		t.Fatalf("callback metadata = %q, want true", callbackMessage.Metadata["goqueue.workflow.chord_callback"])
	}
}

type fakeBackend struct {
	mu                       sync.Mutex
	ensureGroupRequests      []backend.ConsumerGroupRequest
	moveDueRequests          []backend.MoveDueScheduledRequest
	claimStaleReadyRequests  []backend.ClaimStaleReadyRequest
	readReadyCalls           int
	setStateRequests         []backend.TaskStateRecord
	resultRequests           []backend.TaskResultRecord
	ackRequests              []backend.AckRequest
	deadLetterRequests       []backend.DeadLetterRequest
	enqueueReadyRequests     []backend.EnqueueRequest
	advanceChainRequests     []backend.AdvanceWorkflowChainRequest
	recordGroupRequests      []backend.RecordWorkflowTaskCompletedRequest
	enqueueScheduledRequests []backend.EnqueueRequest
	ensureConsumerGroupErr   error
	moveDueScheduledErr      error
	readReadyErr             error
	readReadyFn              func(context.Context, backend.ReadReadyRequest) ([]backend.ReadyMessage, error)
	claimStaleReadyFn        func(context.Context, backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error)
	ensureConsumerGroupFn    func(backend.ConsumerGroupRequest) error
	moveDueFn                func(context.Context, backend.MoveDueScheduledRequest) ([]backend.MovedScheduledMessage, error)
	ackFn                    func(context.Context, backend.AckRequest) error
	setTaskStateErr          error
	saveTaskResultErr        error
	deadLetterErr            error
	enqueueScheduledErr      error
	advanceChainResponse     backend.AdvanceWorkflowChainResponse
	recordGroupResponse      backend.WorkflowGroupProgress
	enqueueScheduledFn       func(context.Context, backend.EnqueueRequest) (backend.EnqueueResponse, error)
}

func (f *fakeBackend) EnqueueReady(_ context.Context, req backend.EnqueueRequest) (backend.EnqueueResponse, error) {
	f.mu.Lock()
	f.enqueueReadyRequests = append(f.enqueueReadyRequests, req)
	f.mu.Unlock()
	return backend.EnqueueResponse{TaskID: task.TaskID(req.Message.ID), StreamID: "next-1"}, nil
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

func (f *fakeBackend) ClaimStaleReady(ctx context.Context, req backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
	f.mu.Lock()
	f.claimStaleReadyRequests = append(f.claimStaleReadyRequests, req)
	f.mu.Unlock()

	if f.claimStaleReadyFn != nil {
		return f.claimStaleReadyFn(ctx, req)
	}

	return nil, nil
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

func (f *fakeBackend) EnqueueDeadLetter(_ context.Context, req backend.DeadLetterRequest) (backend.DeadLetterRecord, error) {
	f.mu.Lock()
	f.deadLetterRequests = append(f.deadLetterRequests, req)
	f.mu.Unlock()

	if f.deadLetterErr != nil {
		return backend.DeadLetterRecord{}, f.deadLetterErr
	}

	return backend.DeadLetterRecord{
		StreamID: "dead-1",
		Message:  req.Message,
		Reason:   req.Reason,
		Error:    req.Error,
		FailedAt: req.FailedAt,
	}, nil
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

func (f *fakeBackend) AdvanceWorkflowChain(_ context.Context, req backend.AdvanceWorkflowChainRequest) (backend.AdvanceWorkflowChainResponse, error) {
	f.mu.Lock()
	f.advanceChainRequests = append(f.advanceChainRequests, req)
	f.mu.Unlock()
	return f.advanceChainResponse, nil
}

func (f *fakeBackend) SaveWorkflowGroup(_ context.Context, _ backend.WorkflowGroupRecord) error {
	return nil
}

func (f *fakeBackend) RecordWorkflowTaskCompleted(_ context.Context, req backend.RecordWorkflowTaskCompletedRequest) (backend.WorkflowGroupProgress, error) {
	f.mu.Lock()
	f.recordGroupRequests = append(f.recordGroupRequests, req)
	f.mu.Unlock()
	return f.recordGroupResponse, nil
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
	metadata    map[string]string
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
		Metadata:    input.metadata,
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

func makeClaimOnce(message backend.ReadyMessage) func(context.Context, backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
	called := false

	return func(context.Context, backend.ClaimStaleReadyRequest) ([]backend.ReadyMessage, error) {
		if called {
			return []backend.ReadyMessage{}, nil
		}
		called = true

		return []backend.ReadyMessage{message}, nil
	}
}
