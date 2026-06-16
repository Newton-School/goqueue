package workflow

import (
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/task"
)

type canvasConfig struct {
	defaultQueue task.QueueName
	codec        task.PayloadCodec
	now          func() time.Time
}

// CanvasOption customizes canvas workflow dispatch.
type CanvasOption func(*canvasConfig) error

func defaultCanvasConfig() canvasConfig {
	return canvasConfig{
		defaultQueue: "default",
		codec:        task.JSONPayloadCodec{},
		now:          utcNow,
	}
}

func utcNow() time.Time {
	return time.Now().UTC()
}

// WithCanvasDefaultQueue sets the queue used when a signature omits one.
func WithCanvasDefaultQueue(queue task.QueueName) CanvasOption {
	return func(config *canvasConfig) error {
		if err := task.ValidateQueueName(queue.String()); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidWorkflow, err)
		}

		config.defaultQueue = queue
		return nil
	}
}

// WithCanvasCodec sets the payload codec used for dispatched signatures.
func WithCanvasCodec(codec task.PayloadCodec) CanvasOption {
	return func(config *canvasConfig) error {
		if codec == nil {
			return fmt.Errorf("%w: codec is required", ErrInvalidWorkflow)
		}

		config.codec = codec
		return nil
	}
}

// WithCanvasNow makes workflow dispatch deterministic for tests.
func WithCanvasNow(now func() time.Time) CanvasOption {
	return func(config *canvasConfig) error {
		if now == nil {
			return fmt.Errorf("%w: now function is required", ErrInvalidWorkflow)
		}

		config.now = now
		return nil
	}
}
