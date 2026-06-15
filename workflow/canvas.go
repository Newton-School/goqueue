package workflow

import (
	"context"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/producer"
	"github.com/Newton-School/goqueue/task"
)

// Canvas dispatches workflow primitives through a queue backend.
type Canvas struct {
	backend      backend.QueueBackend
	producer     *producer.Producer
	defaultQueue task.QueueName
	codec        task.PayloadCodec
	now          func() time.Time
}

// NewCanvas creates a workflow canvas bound to a backend.
func NewCanvas(queueBackend backend.QueueBackend, opts ...CanvasOption) (*Canvas, error) {
	if queueBackend == nil {
		return nil, ErrNilBackend
	}

	config := defaultCanvasConfig()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(&config); err != nil {
			return nil, err
		}
	}

	dispatchProducer, err := producer.NewProducer(
		queueBackend,
		producer.WithProducerDefaultQueue(config.defaultQueue),
		producer.WithProducerCodec(config.codec),
		producer.WithProducerNow(config.now),
	)
	if err != nil {
		return nil, err
	}

	return &Canvas{
		backend:      queueBackend,
		producer:     dispatchProducer,
		defaultQueue: config.defaultQueue,
		codec:        config.codec,
		now:          config.now,
	}, nil
}

// ApplySignature dispatches one signature through the producer path.
func (c *Canvas) ApplySignature(ctx context.Context, signature Signature) (*producer.AsyncResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	normalized, err := signature.Normalize(c.defaultQueue)
	if err != nil {
		return nil, err
	}

	return c.producer.ApplyAsync(
		ctx,
		normalized.Name,
		copyAnySlice(normalized.Args),
		copyAnyMap(normalized.Kwargs),
		producer.WithApplyQueue(normalized.Queue),
		producer.WithApplyMetadata(copyStringMap(normalized.Metadata)),
		producer.WithApplyPriority(normalized.Priority),
		producer.WithApplyRetryPolicy(normalized.RetryPolicy),
		producer.WithApplyETA(normalized.Timing.ETA),
		producer.WithApplyExpiresAt(normalized.Timing.ExpiresAt),
		producer.WithApplyCreatedAt(c.now()),
	)
}

func newWorkflowID() (task.TaskID, error) {
	return task.NewTaskID()
}
