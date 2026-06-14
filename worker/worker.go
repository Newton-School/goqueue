package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

// Worker processes queue messages from a Redis stream-backed backend.
type Worker struct {
	backend                backend.QueueBackend
	registry               *task.TaskRegistry
	queue                  task.QueueName
	group                  string
	consumer               string
	codec                  task.PayloadCodec
	concurrency            int
	readBatch              int64
	block                  time.Duration
	moveDueEnabled         bool
	moveDueLimit           int64
	idleDelay              time.Duration
	deadLetterEnabled      bool
	pendingRecoveryEnabled bool
	pendingMinIdle         time.Duration
	pendingClaimBatch      int64
	pendingClaimInterval   time.Duration
	now                    func() time.Time
}

// NewWorker creates a worker for a queue, registry, and backend.
func NewWorker(queueBackend backend.QueueBackend, registry *task.TaskRegistry, opts ...WorkerOption) (*Worker, error) {
	if queueBackend == nil {
		return nil, ErrNilBackend
	}
	if registry == nil {
		return nil, ErrNilTaskRegistry
	}

	config := defaultWorkerConfig()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(&config); err != nil {
			return nil, err
		}
	}

	if err := task.ValidateQueueName(config.queue.String()); err != nil {
		return nil, err
	}
	if config.group == "" {
		return nil, fmt.Errorf("%w: consumer group is required", ErrInvalidWorkerOption)
	}
	if config.consumer == "" {
		return nil, fmt.Errorf("%w: consumer is required", ErrInvalidWorkerOption)
	}
	if config.codec == nil {
		return nil, ErrInvalidWorkerOption
	}
	if config.concurrency < 1 {
		return nil, fmt.Errorf("%w: concurrency must be at least 1", ErrInvalidWorkerOption)
	}
	if config.readBatch < 1 {
		return nil, fmt.Errorf("%w: read batch must be at least 1", ErrInvalidWorkerOption)
	}
	if config.block < 0 {
		return nil, fmt.Errorf("%w: block duration cannot be negative", ErrInvalidWorkerOption)
	}
	if config.moveDueLimit < 1 {
		return nil, fmt.Errorf("%w: move due limit must be at least 1", ErrInvalidWorkerOption)
	}
	if config.idleDelay < 0 {
		return nil, fmt.Errorf("%w: idle delay cannot be negative", ErrInvalidWorkerOption)
	}
	if config.pendingMinIdle < 0 {
		return nil, fmt.Errorf("%w: pending min idle cannot be negative", ErrInvalidWorkerOption)
	}
	if config.pendingClaimBatch < 1 {
		return nil, fmt.Errorf("%w: pending claim batch must be at least 1", ErrInvalidWorkerOption)
	}
	if config.pendingClaimInterval < 0 {
		return nil, fmt.Errorf("%w: pending claim interval cannot be negative", ErrInvalidWorkerOption)
	}
	if config.now == nil {
		config.now = time.Now().UTC
	}

	return &Worker{
		backend:                queueBackend,
		registry:               registry,
		queue:                  config.queue,
		group:                  config.group,
		consumer:               config.consumer,
		codec:                  config.codec,
		concurrency:            config.concurrency,
		readBatch:              config.readBatch,
		block:                  config.block,
		moveDueEnabled:         config.moveDueEnabled,
		moveDueLimit:           config.moveDueLimit,
		idleDelay:              config.idleDelay,
		deadLetterEnabled:      config.deadLetterEnabled,
		pendingRecoveryEnabled: config.pendingRecoveryEnabled,
		pendingMinIdle:         config.pendingMinIdle,
		pendingClaimBatch:      config.pendingClaimBatch,
		pendingClaimInterval:   config.pendingClaimInterval,
		now:                    config.now,
	}, nil
}

// Start runs the worker loop until context is canceled.
func (w *Worker) Start(ctx context.Context) error {
	if w == nil {
		return fmt.Errorf("%w: goqueue worker is nil", ErrNilWorker)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := w.backend.EnsureConsumerGroup(runCtx, backend.ConsumerGroupRequest{
		Queue: w.queue,
		Group: w.group,
	}); err != nil {
		return fmt.Errorf("goqueue worker: ensure consumer group: %w", err)
	}

	sem := make(chan struct{}, w.concurrency)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup

loop:
	for {
		select {
		case err := <-errCh:
			wg.Wait()
			return err
		default:
		}

		if runCtx.Err() != nil {
			break loop
		}

		if w.moveDueEnabled {
			if err := w.moveDueScheduled(runCtx); err != nil {
				return err
			}
		}

		readies, err := w.backend.ReadReady(runCtx, backend.ReadReadyRequest{
			Queue:    w.queue,
			Group:    w.group,
			Consumer: w.consumer,
			Count:    w.readBatch,
			Block:    w.block,
		})
		if err != nil {
			if runCtx.Err() != nil {
				break loop
			}
			return fmt.Errorf("goqueue worker: read ready queue: %w", err)
		}

		if len(readies) == 0 {
			if w.block == 0 && w.idleDelay > 0 {
				select {
				case <-runCtx.Done():
					break loop
				case <-time.After(w.idleDelay):
				}
			}
			continue
		}

		for _, ready := range readies {
			select {
			case sem <- struct{}{}:
			case <-runCtx.Done():
				break loop
			}

			wg.Add(1)
			go func(message backend.ReadyMessage) {
				defer wg.Done()
				defer func() {
					<-sem
				}()

				if err := w.processMessage(runCtx, message); err != nil {
					select {
					case errCh <- err:
						cancel()
					default:
					}
				}
			}(ready)
		}
	}

	wg.Wait()
	select {
	case err := <-errCh:
		return err
	default:
	}
	return nil
}

func (w *Worker) moveDueScheduled(ctx context.Context) error {
	if _, err := w.backend.MoveDueScheduled(ctx, backend.MoveDueScheduledRequest{
		Queue: w.queue,
		Now:   w.now(),
		Limit: w.moveDueLimit,
	}); err != nil {
		return fmt.Errorf("goqueue worker: move due scheduled: %w", err)
	}
	return nil
}

func (w *Worker) processMessage(ctx context.Context, message backend.ReadyMessage) error {
	envelope, err := task.TaskMessageToEnvelope(message.Message, w.codec)
	if err != nil {
		if err := w.ack(ctx, message.StreamID); err != nil {
			return err
		}

		return fmt.Errorf("goqueue worker: deserialize task message: %w", err)
	}

	if err := w.writeState(ctx, envelope.ID, task.TaskReceived, ""); err != nil {
		return err
	}

	expired, err := w.checkExpired(ctx, envelope, message)
	if err != nil {
		return err
	}
	if expired {
		return w.ack(ctx, message.StreamID)
	}

	handler, err := w.registry.Lookup(envelope.Name)
	if err != nil {
		result := task.FailedResult(err)
		if err := w.deadLetterTask(ctx, message.StreamID, envelope, message, task.FailureUnknownTask, result); err != nil {
			return err
		}
		return w.ack(ctx, message.StreamID)
	}

	if err := w.writeState(ctx, envelope.ID, task.TaskStarted, ""); err != nil {
		return err
	}

	handlerCtx := task.NewHandlerContext(ctx, envelope)
	result, err := handler.HandleTask(handlerCtx, envelope.Payload)
	if err != nil {
		result = task.FailedResult(err)
	}
	result = normalizeTaskResult(result)

	if shouldRetry(result.State, envelope) {
		if err := w.retryTask(ctx, message.StreamID, envelope, result); err != nil {
			return err
		}
		return nil
	}
	if isRetryExhausted(result.State, envelope) {
		if err := w.deadLetterTask(ctx, message.StreamID, envelope, message, task.FailureRetryExhausted, result); err != nil {
			return err
		}
		return w.ack(ctx, message.StreamID)
	}

	finalState := result.State
	if result.State == task.TaskRetrying {
		finalState = task.TaskFailed
	}
	if finalState != task.TaskSucceeded && finalState != task.TaskFailed && finalState != task.TaskRevoked &&
		finalState != task.TaskExpired && finalState != task.TaskDeadLettered {
		finalState = task.TaskFailed
	}

	if err := w.writeState(ctx, envelope.ID, finalState, result.Error); err != nil {
		return err
	}
	if err := w.saveResult(ctx, envelope.ID, result); err != nil {
		return err
	}

	return w.ack(ctx, message.StreamID)
}

func (w *Worker) retryTask(ctx context.Context, streamID string, envelope task.TaskEnvelope, result task.TaskResult) error {
	nextAttempt := envelope.Attempt + 1
	nextDelay := envelope.RetryPolicy.DelayForAttempt(nextAttempt)

	if !envelope.Timing.ExpiresAt.IsZero() && w.now().Add(nextDelay).After(envelope.Timing.ExpiresAt) {
		result = task.FailedResult(fmt.Errorf("task expired before retry"))
		if err := w.writeState(ctx, envelope.ID, task.TaskExpired, result.Error); err != nil {
			return err
		}
		if err := w.saveResult(ctx, envelope.ID, result); err != nil {
			return err
		}
		return w.ack(ctx, streamID)
	}

	retryEnvelope, err := task.NewTaskEnvelope(task.TaskEnvelopeInput{
		ID:       envelope.ID,
		Name:     envelope.Name,
		Queue:    envelope.Queue,
		Args:     envelope.Payload.Args(),
		Kwargs:   envelope.Payload.Kwargs(),
		Metadata: envelope.Metadata.Values(),
		Timing: task.TaskTiming{
			ETA:       w.now().Add(nextDelay),
			ExpiresAt: envelope.Timing.ExpiresAt,
		},
		Priority:    envelope.Priority,
		RetryPolicy: envelope.RetryPolicy,
		Attempt:     nextAttempt,
	})
	if err != nil {
		return fmt.Errorf("goqueue worker: build retry envelope: %w", err)
	}

	retryMessage, err := task.TaskEnvelopeToMessage(retryEnvelope, w.codec)
	if err != nil {
		return fmt.Errorf("goqueue worker: encode retry message: %w", err)
	}

	if err := w.writeState(ctx, envelope.ID, task.TaskRetrying, result.Error); err != nil {
		return err
	}
	if err := w.saveResult(ctx, envelope.ID, result); err != nil {
		return err
	}
	if _, err := w.backend.EnqueueScheduled(ctx, backend.EnqueueRequest{Message: retryMessage}); err != nil {
		if failedErr := w.writeState(ctx, envelope.ID, task.TaskFailed, fmt.Sprintf("retry schedule failed: %v", err)); failedErr != nil {
			return fmt.Errorf("goqueue worker: retry schedule failed: %w; state write failed: %v", err, failedErr)
		}
		return fmt.Errorf("goqueue worker: retry schedule failed: %w", err)
	}

	return w.ack(ctx, streamID)
}

func (w *Worker) checkExpired(ctx context.Context, envelope task.TaskEnvelope, message backend.ReadyMessage) (bool, error) {
	if envelope.Timing.ExpiresAt.IsZero() {
		return false, nil
	}

	if envelope.Timing.ExpiresAt.After(w.now()) {
		return false, nil
	}

	result := task.FailedResult(fmt.Errorf("task expired"))
	result, err := w.recordDeadLetter(ctx, message.StreamID, envelope, message, task.FailureExpired, result)
	if err != nil {
		return false, err
	}
	if err := w.writeState(ctx, envelope.ID, task.TaskExpired, result.Error); err != nil {
		return false, err
	}
	if err := w.saveResult(ctx, envelope.ID, result); err != nil {
		return false, err
	}

	return true, nil
}

func normalizeTaskResult(result task.TaskResult) task.TaskResult {
	if err := result.Validate(); err != nil || result.State == "" {
		return task.FailedResult(fmt.Errorf("invalid task result"))
	}

	if !isAllowedFinalState(result.State) {
		return task.FailedResult(fmt.Errorf("result state %q is not terminal", result.State))
	}
	return result
}

func isAllowedFinalState(state task.TaskState) bool {
	switch state {
	case task.TaskSucceeded, task.TaskFailed, task.TaskRevoked, task.TaskExpired, task.TaskDeadLettered, task.TaskRetrying:
		return true
	default:
		return false
	}
}

func shouldRetry(state task.TaskState, envelope task.TaskEnvelope) bool {
	if state != task.TaskFailed && state != task.TaskRetrying {
		return false
	}
	if envelope.Attempt+1 >= envelope.RetryPolicy.MaxAttempts {
		return false
	}

	return true
}

func isRetryExhausted(state task.TaskState, envelope task.TaskEnvelope) bool {
	if state != task.TaskFailed && state != task.TaskRetrying {
		return false
	}
	return envelope.Attempt+1 >= envelope.RetryPolicy.MaxAttempts
}

func (w *Worker) writeState(ctx context.Context, taskID task.TaskID, state task.TaskState, msg string) error {
	return w.backend.SetTaskState(ctx, backend.TaskStateRecord{
		TaskID: taskID,
		State:  state,
		Error:  msg,
	})
}

func (w *Worker) saveResult(ctx context.Context, taskID task.TaskID, result task.TaskResult) error {
	return w.backend.SaveTaskResult(ctx, backend.TaskResultRecord{
		TaskID: taskID,
		Result: result,
	})
}

func (w *Worker) ack(ctx context.Context, streamID string) error {
	return w.backend.Ack(ctx, backend.AckRequest{
		Queue:    w.queue,
		Group:    w.group,
		StreamID: streamID,
	})
}
