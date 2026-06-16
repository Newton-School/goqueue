package producer

import (
	"context"
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

// Producer publishes task envelopes to backend queues.
type Producer struct {
	backend      backend.QueueBackend
	defaultQueue task.QueueName
	codec        task.PayloadCodec
	now          func() time.Time
}

// NewProducer creates a producer bound to a queue backend.
func NewProducer(queueBackend backend.QueueBackend, opts ...ProducerOption) (*Producer, error) {
	if queueBackend == nil {
		return nil, ErrNilBackend
	}

	config := defaultProducerConfig()
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if err := opt(&config); err != nil {
			return nil, err
		}
	}

	if err := task.ValidateQueueName(config.defaultQueue.String()); err != nil {
		return nil, err
	}

	if config.codec == nil {
		config.codec = task.JSONPayloadCodec{}
	}

	if config.now == nil {
		config.now = utcNow
	}

	return &Producer{
		backend:      queueBackend,
		defaultQueue: config.defaultQueue,
		codec:        config.codec,
		now:          config.now,
	}, nil
}

// ApplyAsync creates a task and publishes it to the queue.
func (p *Producer) ApplyAsync(
	ctx context.Context,
	name task.TaskName,
	args []any,
	kwargs map[string]any,
	options ...ApplyOption,
) (*AsyncResult, error) {
	if name == "" {
		return nil, ErrMissingTaskName
	}
	if err := task.ValidateTaskName(name.String()); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("producer: is nil")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	applyConfig := defaultApplyConfig()
	for _, applyOpt := range options {
		if applyOpt == nil {
			continue
		}

		if err := applyOpt(&applyConfig); err != nil {
			return nil, err
		}
	}

	queue := applyConfig.queue
	if queue == "" {
		queue = p.defaultQueue
	}

	if applyConfig.timing.ETA.IsZero() && applyConfig.countdown != nil {
		timing, err := task.TaskTimingFromCountdown(p.now(), *applyConfig.countdown)
		if err != nil {
			return nil, fmt.Errorf("producer: cannot resolve countdown: %w", err)
		}

		applyConfig.timing.ETA = timing.ETA
	}

	envelope, err := task.NewTaskEnvelope(task.TaskEnvelopeInput{
		ID:          applyConfig.id,
		Name:        name,
		Queue:       queue,
		Args:        args,
		Kwargs:      kwargs,
		Metadata:    applyConfig.metadata,
		Timing:      applyConfig.timing,
		Priority:    applyConfig.priority,
		RetryPolicy: applyConfig.retryPolicy,
		CreatedAt:   applyConfig.createdAt,
		Attempt:     applyConfig.Attempt,
	})
	if err != nil {
		return nil, fmt.Errorf("producer: build envelope: %w", err)
	}

	message, err := task.TaskEnvelopeToMessage(envelope, p.codec)
	if err != nil {
		return nil, fmt.Errorf("producer: serialize message: %w", err)
	}

	initialState := task.TaskPending
	if envelope.Timing.Scheduled() {
		initialState = task.TaskScheduled
	}

	enqueueStateErr := p.backend.SetTaskState(ctx, backend.TaskStateRecord{
		TaskID: envelope.ID,
		State:  initialState,
	})
	if enqueueStateErr != nil {
		return nil, fmt.Errorf("producer: set initial state: %w", enqueueStateErr)
	}

	req := backend.EnqueueRequest{Message: message}
	if envelope.Timing.Scheduled() {
		if _, err := p.backend.EnqueueScheduled(ctx, req); err != nil {
			return nil, p.markFailedIfPossible(ctx, envelope.ID, err, "enqueue scheduled")
		}

		return &AsyncResult{taskID: envelope.ID, backend: p.backend}, nil
	}

	if _, err := p.backend.EnqueueReady(ctx, req); err != nil {
		return nil, p.markFailedIfPossible(ctx, envelope.ID, err, "enqueue ready")
	}

	return &AsyncResult{taskID: envelope.ID, backend: p.backend}, nil
}

// NewAsyncResult creates a handle for an already known task id.
func NewAsyncResult(taskID task.TaskID, queueBackend backend.QueueBackend) *AsyncResult {
	return &AsyncResult{taskID: taskID, backend: queueBackend}
}

func (p *Producer) markFailedIfPossible(ctx context.Context, taskID task.TaskID, enqueueErr error, action string) error {
	stateErr := p.backend.SetTaskState(ctx, backend.TaskStateRecord{
		TaskID: taskID,
		State:  task.TaskFailed,
		Error:  enqueueErr.Error(),
	})

	if stateErr == nil {
		return fmt.Errorf("producer: %s: %w", action, enqueueErr)
	}

	return fmt.Errorf("producer: %s: %w; state write failed: %v", action, enqueueErr, stateErr)
}
