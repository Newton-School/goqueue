package backend

import (
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/task"
)

// ReadReadyRequest asks a backend to read ready messages for a consumer group.
type ReadReadyRequest struct {
	Queue    task.QueueName
	Group    string
	Consumer string
	Count    int64
	Block    time.Duration
}

// ReadyMessage is a task message read from a ready queue.
type ReadyMessage struct {
	StreamID string
	Message  task.TaskMessage
}

// AckRequest asks a backend to acknowledge a stream message.
type AckRequest struct {
	Queue    task.QueueName
	Group    string
	StreamID string
}

// Validate verifies read request fields.
func (r ReadReadyRequest) Validate() error {
	if err := task.ValidateQueueName(r.Queue.String()); err != nil {
		return err
	}
	if r.Group == "" {
		return fmt.Errorf("%w: consumer group is required", ErrInvalidBackendRequest)
	}
	if r.Consumer == "" {
		return fmt.Errorf("%w: consumer name is required", ErrInvalidBackendRequest)
	}
	if r.Count < 0 {
		return fmt.Errorf("%w: read count cannot be negative", ErrInvalidBackendRequest)
	}
	if r.Block < 0 {
		return fmt.Errorf("%w: read block duration cannot be negative", ErrInvalidBackendRequest)
	}

	return nil
}

// Validate verifies ack request fields.
func (r AckRequest) Validate() error {
	if err := task.ValidateQueueName(r.Queue.String()); err != nil {
		return err
	}
	if r.Group == "" {
		return fmt.Errorf("%w: consumer group is required", ErrInvalidBackendRequest)
	}
	if r.StreamID == "" {
		return fmt.Errorf("%w: stream id is required", ErrInvalidBackendRequest)
	}

	return nil
}
