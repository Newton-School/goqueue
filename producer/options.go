package producer

import (
	"time"

	"github.com/Newton-School/goqueue/task"
)

// ProducerConfig controls producer defaults.
type ProducerConfig struct {
	defaultQueue task.QueueName
	codec        task.PayloadCodec
	now          func() time.Time
}

// ProducerOption customizes producer defaults.
type ProducerOption func(*ProducerConfig) error

// ApplyConfig controls a single task dispatch.
type ApplyConfig struct {
	queue       task.QueueName
	id          task.TaskID
	metadata    map[string]string
	priority    task.Priority
	retryPolicy task.RetryPolicy
	timing      task.TaskTiming
	countdown   *time.Duration

	createdAt time.Time
	Attempt   int
}

// ApplyOption customizes a single task dispatch.
type ApplyOption func(*ApplyConfig) error

func defaultProducerConfig() ProducerConfig {
	return ProducerConfig{
		defaultQueue: "default",
		codec:        task.JSONPayloadCodec{},
		now:          time.Now().UTC,
	}
}

func defaultApplyConfig() ApplyConfig {
	return ApplyConfig{
		priority:    task.DefaultPriority,
		retryPolicy: task.DefaultRetryPolicy(),
	}
}

// WithProducerDefaultQueue sets a default queue for apply calls.
func WithProducerDefaultQueue(queue task.QueueName) ProducerOption {
	return func(config *ProducerConfig) error {
		if err := task.ValidateQueueName(queue.String()); err != nil {
			return err
		}

		config.defaultQueue = queue
		return nil
	}
}

// WithProducerCodec sets the payload codec used by the producer.
func WithProducerCodec(codec task.PayloadCodec) ProducerOption {
	return func(config *ProducerConfig) error {
		if codec == nil {
			return ErrMissingApplyOption
		}

		config.codec = codec
		return nil
	}
}

// WithProducerNow makes time resolution deterministic for tests.
func WithProducerNow(now func() time.Time) ProducerOption {
	return func(config *ProducerConfig) error {
		if now == nil {
			return ErrMissingApplyOption
		}

		config.now = now
		return nil
	}
}

// WithApplyQueue overrides the queue for a single task.
func WithApplyQueue(queue task.QueueName) ApplyOption {
	return func(config *ApplyConfig) error {
		if err := task.ValidateQueueName(queue.String()); err != nil {
			return err
		}

		config.queue = queue
		return nil
	}
}

// WithApplyTaskID sets a specific task id.
func WithApplyTaskID(id task.TaskID) ApplyOption {
	return func(config *ApplyConfig) error {
		if err := task.ValidateTaskID(id.String()); err != nil {
			return err
		}

		config.id = id
		return nil
	}
}

// WithApplyMetadata sets task metadata.
func WithApplyMetadata(metadata map[string]string) ApplyOption {
	return func(config *ApplyConfig) error {
		config.metadata = metadata
		return nil
	}
}

// WithApplyPriority sets task priority.
func WithApplyPriority(priority task.Priority) ApplyOption {
	return func(config *ApplyConfig) error {
		if err := task.ValidatePriority(priority); err != nil {
			return err
		}

		config.priority = priority
		return nil
	}
}

// WithApplyRetryPolicy sets task retry settings.
func WithApplyRetryPolicy(policy task.RetryPolicy) ApplyOption {
	return func(config *ApplyConfig) error {
		if err := policy.Validate(); err != nil {
			return err
		}

		config.retryPolicy = policy
		return nil
	}
}

// WithApplyCountDown schedules the task after a delay.
func WithApplyCountDown(countdown time.Duration) ApplyOption {
	return func(config *ApplyConfig) error {
		if countdown < 0 {
			return ErrMissingApplyOption
		}

		config.countdown = &countdown
		return nil
	}
}

// WithApplyETA schedules the task at a specific execution time.
func WithApplyETA(eta time.Time) ApplyOption {
	return func(config *ApplyConfig) error {
		config.timing.ETA = eta
		return nil
	}
}

// WithApplyExpiresAt sets task expiration timestamp.
func WithApplyExpiresAt(expiresAt time.Time) ApplyOption {
	return func(config *ApplyConfig) error {
		config.timing.ExpiresAt = expiresAt
		return nil
	}
}

// WithApplyAttempt sets the attempt count for retry replays.
func WithApplyAttempt(attempt int) ApplyOption {
	return func(config *ApplyConfig) error {
		if attempt < 0 {
			return ErrMissingApplyOption
		}

		config.Attempt = attempt
		return nil
	}
}

// WithApplyCreatedAt sets task creation time.
func WithApplyCreatedAt(createdAt time.Time) ApplyOption {
	return func(config *ApplyConfig) error {
		config.createdAt = createdAt
		return nil
	}
}
