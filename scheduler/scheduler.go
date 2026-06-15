package scheduler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/producer"
	"github.com/Newton-School/goqueue/task"
)

// Scheduler registers periodic definitions and dispatches due task instances.
type Scheduler struct {
	backend      backend.QueueBackend
	producer     *producer.Producer
	identity     string
	defaultQueue task.QueueName
	pollInterval time.Duration
	batchSize    int64
	lockTTL      time.Duration
	codec        task.PayloadCodec
	now          func() time.Time
}

// NewScheduler creates a scheduler bound to a backend.
func NewScheduler(queueBackend backend.QueueBackend, opts ...SchedulerOption) (*Scheduler, error) {
	if queueBackend == nil {
		return nil, ErrNilBackend
	}

	config := defaultSchedulerConfig()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(&config); err != nil {
			return nil, err
		}
	}
	if config.identity == "" {
		identity, err := newSchedulerIdentity()
		if err != nil {
			return nil, err
		}
		config.identity = identity
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

	return &Scheduler{
		backend:      queueBackend,
		producer:     dispatchProducer,
		identity:     config.identity,
		defaultQueue: config.defaultQueue,
		pollInterval: config.pollInterval,
		batchSize:    config.batchSize,
		lockTTL:      config.lockTTL,
		codec:        config.codec,
		now:          config.now,
	}, nil
}

// RegisterPeriodicTask stores or updates a periodic task definition.
func (s *Scheduler) RegisterPeriodicTask(ctx context.Context, definition PeriodicTask) error {
	if ctx == nil {
		ctx = context.Background()
	}
	record, err := definition.toBackendRecord(s.defaultQueue, s.now())
	if err != nil {
		return err
	}

	return s.backend.UpsertPeriodicTask(ctx, backend.UpsertPeriodicTaskRequest{Record: record})
}

// DeletePeriodicTask removes a periodic task definition.
func (s *Scheduler) DeletePeriodicTask(ctx context.Context, name PeriodicTaskName) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return s.backend.DeletePeriodicTask(ctx, backend.DeletePeriodicTaskRequest{Name: name.String()})
}

func newSchedulerIdentity() (string, error) {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", fmt.Errorf("%w: generate scheduler identity: %v", ErrInvalidSchedulerOption, err)
	}

	return "scheduler-" + hex.EncodeToString(bytes[:]), nil
}
