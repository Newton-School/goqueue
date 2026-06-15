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

// Identity returns the scheduler identity used for Redis leases.
func (s *Scheduler) Identity() string {
	if s == nil {
		return ""
	}

	return s.identity
}

// DefaultQueue returns the scheduler default queue.
func (s *Scheduler) DefaultQueue() task.QueueName {
	if s == nil {
		return ""
	}

	return s.defaultQueue
}

// PollInterval returns how often Start polls for due definitions.
func (s *Scheduler) PollInterval() time.Duration {
	if s == nil {
		return 0
	}

	return s.pollInterval
}

// BatchSize returns the due definition claim limit.
func (s *Scheduler) BatchSize() int64 {
	if s == nil {
		return 0
	}

	return s.batchSize
}

// LockTTL returns the due definition lease duration.
func (s *Scheduler) LockTTL() time.Duration {
	if s == nil {
		return 0
	}

	return s.lockTTL
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

// PollOnce leases due periodic tasks, dispatches them, and advances their schedules.
func (s *Scheduler) PollOnce(ctx context.Context) (int, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	now := s.now().UTC()
	dueTasks, err := s.backend.ListDuePeriodicTasks(ctx, backend.ListDuePeriodicTasksRequest{
		Now:         now,
		Limit:       s.batchSize,
		SchedulerID: s.identity,
		LockTTL:     s.lockTTL,
	})
	if err != nil {
		return 0, err
	}

	dispatched := 0
	for _, due := range dueTasks {
		if err := due.Validate(); err != nil {
			return dispatched, err
		}

		definition, err := periodicTaskFromBackendRecord(due.Record)
		if err != nil {
			return dispatched, err
		}

		result, err := s.producer.ApplyAsync(
			ctx,
			definition.TaskName,
			copyAnySlice(definition.Args),
			copyAnyMap(definition.Kwargs),
			producer.WithApplyQueue(definition.Queue),
			producer.WithApplyMetadata(periodicDispatchMetadata(
				definition.Metadata,
				due.Record.Name,
				due.Record.NextDueAt.UTC().Format(time.RFC3339Nano),
			)),
			producer.WithApplyPriority(definition.Priority),
			producer.WithApplyRetryPolicy(definition.RetryPolicy),
			producer.WithApplyCreatedAt(now),
		)
		if err != nil {
			return dispatched, err
		}

		if err := s.backend.MarkPeriodicTaskDispatched(ctx, backend.MarkPeriodicTaskDispatchedRequest{
			Name:             due.Record.Name,
			LockToken:        due.LockToken,
			DispatchedTaskID: result.ID(),
			DispatchedAt:     now,
			NextDueAt:        definition.NextDueAfter(now),
		}); err != nil {
			return dispatched, err
		}

		dispatched++
	}

	return dispatched, nil
}

// Start runs the scheduler loop until ctx is canceled.
func (s *Scheduler) Start(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil
	}

	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		if _, err := s.PollOnce(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func newSchedulerIdentity() (string, error) {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", fmt.Errorf("%w: generate scheduler identity: %v", ErrInvalidSchedulerOption, err)
	}

	return "scheduler-" + hex.EncodeToString(bytes[:]), nil
}
