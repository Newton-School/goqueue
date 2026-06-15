package workflow

import (
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

func newWorkflowID() (task.TaskID, error) {
	return task.NewTaskID()
}
