package worker

import (
	"fmt"
	"strings"
	"time"

	"github.com/Newton-School/goqueue/task"
)

// WorkerConfig holds worker defaults.
type WorkerConfig struct {
	queue          task.QueueName
	group          string
	consumer       string
	codec          task.PayloadCodec
	concurrency    int
	readBatch      int64
	block          time.Duration
	moveDueEnabled bool
	moveDueLimit   int64
	idleDelay      time.Duration
	now            func() time.Time
}

// WorkerOption customizes worker behavior.
type WorkerOption func(*WorkerConfig) error

func defaultWorkerConfig() WorkerConfig {
	return WorkerConfig{
		queue:          "default",
		group:          "goqueue",
		consumer:       "worker",
		codec:          task.JSONPayloadCodec{},
		concurrency:    1,
		readBatch:      1,
		block:          250 * time.Millisecond,
		moveDueEnabled: true,
		moveDueLimit:   100,
		idleDelay:      50 * time.Millisecond,
		now:            time.Now().UTC,
	}
}

// WithWorkerQueue sets the queue the worker reads from.
func WithWorkerQueue(queue task.QueueName) WorkerOption {
	return func(config *WorkerConfig) error {
		if err := task.ValidateQueueName(queue.String()); err != nil {
			return err
		}

		config.queue = queue
		return nil
	}
}

// WithWorkerGroup sets the consumer group name.
func WithWorkerGroup(group string) WorkerOption {
	return func(config *WorkerConfig) error {
		group = strings.TrimSpace(group)
		if group == "" {
			return fmt.Errorf("%w: consumer group is required", ErrInvalidWorkerOption)
		}

		config.group = group
		return nil
	}
}

// WithWorkerConsumer sets the consumer name for Redis consumer-group reads.
func WithWorkerConsumer(consumer string) WorkerOption {
	return func(config *WorkerConfig) error {
		consumer = strings.TrimSpace(consumer)
		if consumer == "" {
			return fmt.Errorf("%w: consumer name is required", ErrInvalidWorkerOption)
		}

		config.consumer = consumer
		return nil
	}
}

// WithWorkerCodec sets the payload codec used while decoding messages.
func WithWorkerCodec(codec task.PayloadCodec) WorkerOption {
	return func(config *WorkerConfig) error {
		if codec == nil {
			return fmt.Errorf("%w: codec is required", ErrInvalidWorkerOption)
		}

		config.codec = codec
		return nil
	}
}

// WithWorkerConcurrency sets the number of goroutines that process tasks.
func WithWorkerConcurrency(concurrency int) WorkerOption {
	return func(config *WorkerConfig) error {
		if concurrency < 1 {
			return fmt.Errorf("%w: concurrency must be at least 1", ErrInvalidWorkerOption)
		}

		config.concurrency = concurrency
		return nil
	}
}

// WithWorkerReadBatch sets how many messages can be read per poll.
func WithWorkerReadBatch(readBatch int64) WorkerOption {
	return func(config *WorkerConfig) error {
		if readBatch < 1 {
			return fmt.Errorf("%w: read batch must be at least 1", ErrInvalidWorkerOption)
		}

		config.readBatch = readBatch
		return nil
	}
}

// WithWorkerBlock sets the XREADGROUP block duration.
func WithWorkerBlock(block time.Duration) WorkerOption {
	return func(config *WorkerConfig) error {
		if block < 0 {
			return fmt.Errorf("%w: block duration cannot be negative", ErrInvalidWorkerOption)
		}

		config.block = block
		return nil
	}
}

// WithWorkerMoveDueLimit sets max number of scheduled messages moved per loop.
func WithWorkerMoveDueLimit(limit int64) WorkerOption {
	return func(config *WorkerConfig) error {
		if limit < 1 {
			return fmt.Errorf("%w: move due limit must be at least 1", ErrInvalidWorkerOption)
		}

		config.moveDueLimit = limit
		return nil
	}
}

// WithWorkerMoveDueEnabled toggles scheduled-task migration into ready queues.
func WithWorkerMoveDueEnabled(enabled bool) WorkerOption {
	return func(config *WorkerConfig) error {
		config.moveDueEnabled = enabled
		return nil
	}
}

// WithWorkerIdleDelay sets delay between immediate polls when no message is returned.
func WithWorkerIdleDelay(delay time.Duration) WorkerOption {
	return func(config *WorkerConfig) error {
		if delay < 0 {
			return fmt.Errorf("%w: idle delay cannot be negative", ErrInvalidWorkerOption)
		}

		config.idleDelay = delay
		return nil
	}
}

// WithWorkerNow sets deterministic time in tests.
func WithWorkerNow(now func() time.Time) WorkerOption {
	return func(config *WorkerConfig) error {
		if now == nil {
			return fmt.Errorf("%w: now function is required", ErrInvalidWorkerOption)
		}

		config.now = now
		return nil
	}
}
