package scheduler

import (
	"fmt"
	"strings"
	"time"

	"github.com/Newton-School/goqueue/task"
)

type schedulerConfig struct {
	identity     string
	defaultQueue task.QueueName
	pollInterval time.Duration
	batchSize    int64
	lockTTL      time.Duration
	codec        task.PayloadCodec
	now          func() time.Time
}

// SchedulerOption customizes scheduler runtime behavior.
type SchedulerOption func(*schedulerConfig) error

func defaultSchedulerConfig() schedulerConfig {
	return schedulerConfig{
		defaultQueue: "default",
		pollInterval: time.Second,
		batchSize:    100,
		lockTTL:      30 * time.Second,
		codec:        task.JSONPayloadCodec{},
		now:          utcNow,
	}
}

func utcNow() time.Time {
	return time.Now().UTC()
}

// WithSchedulerIdentity sets the scheduler identity used in Redis leases.
func WithSchedulerIdentity(identity string) SchedulerOption {
	return func(config *schedulerConfig) error {
		identity = strings.TrimSpace(identity)
		if identity == "" {
			return fmt.Errorf("%w: identity is required", ErrInvalidSchedulerOption)
		}

		config.identity = identity
		return nil
	}
}

// WithSchedulerDefaultQueue sets the queue used when a periodic task omits one.
func WithSchedulerDefaultQueue(queue task.QueueName) SchedulerOption {
	return func(config *schedulerConfig) error {
		if err := task.ValidateQueueName(queue.String()); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidSchedulerOption, err)
		}

		config.defaultQueue = queue
		return nil
	}
}

// WithSchedulerPollInterval sets how often Start checks for due tasks.
func WithSchedulerPollInterval(interval time.Duration) SchedulerOption {
	return func(config *schedulerConfig) error {
		if interval <= 0 {
			return fmt.Errorf("%w: poll interval must be positive", ErrInvalidSchedulerOption)
		}

		config.pollInterval = interval
		return nil
	}
}

// WithSchedulerBatchSize sets the maximum number of due tasks claimed per poll.
func WithSchedulerBatchSize(size int64) SchedulerOption {
	return func(config *schedulerConfig) error {
		if size <= 0 {
			return fmt.Errorf("%w: batch size must be positive", ErrInvalidSchedulerOption)
		}

		config.batchSize = size
		return nil
	}
}

// WithSchedulerLockTTL sets how long Redis due-task leases live.
func WithSchedulerLockTTL(ttl time.Duration) SchedulerOption {
	return func(config *schedulerConfig) error {
		if ttl <= 0 {
			return fmt.Errorf("%w: lock ttl must be positive", ErrInvalidSchedulerOption)
		}

		config.lockTTL = ttl
		return nil
	}
}

// WithSchedulerCodec sets the payload codec used by the dispatch producer.
func WithSchedulerCodec(codec task.PayloadCodec) SchedulerOption {
	return func(config *schedulerConfig) error {
		if codec == nil {
			return fmt.Errorf("%w: codec is required", ErrInvalidSchedulerOption)
		}

		config.codec = codec
		return nil
	}
}

// WithSchedulerNow makes scheduler time deterministic for tests.
func WithSchedulerNow(now func() time.Time) SchedulerOption {
	return func(config *schedulerConfig) error {
		if now == nil {
			return fmt.Errorf("%w: now is required", ErrInvalidSchedulerOption)
		}

		config.now = now
		return nil
	}
}
