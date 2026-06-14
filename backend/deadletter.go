package backend

import (
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/task"
)

// DeadLetterRequest asks a backend to persist an unrecoverable task message.
type DeadLetterRequest struct {
	Message        task.TaskMessage
	Reason         task.FailureCategory
	Error          string
	SourceStreamID string
	Group          string
	Consumer       string
	FailedAt       time.Time
}

// DeadLetterRecord describes a task message stored in a dead-letter queue.
type DeadLetterRecord struct {
	StreamID       string
	Message        task.TaskMessage
	Reason         task.FailureCategory
	Error          string
	SourceStreamID string
	Group          string
	Consumer       string
	FailedAt       time.Time
}

// ReadDeadLettersRequest asks a backend to inspect dead-lettered messages.
type ReadDeadLettersRequest struct {
	Queue task.QueueName
	Count int64
}

// Validate verifies the dead-letter request can be stored safely.
func (r DeadLetterRequest) Validate() error {
	if err := (EnqueueRequest{Message: r.Message}).Validate(); err != nil {
		return err
	}
	if r.Reason == "" {
		return fmt.Errorf("%w: dead letter reason is required", ErrInvalidBackendRequest)
	}
	if r.SourceStreamID == "" {
		return fmt.Errorf("%w: source stream id is required", ErrInvalidBackendRequest)
	}

	return nil
}

// Validate verifies the dead-letter read request.
func (r ReadDeadLettersRequest) Validate() error {
	if err := task.ValidateQueueName(r.Queue.String()); err != nil {
		return err
	}
	if r.Count < 0 {
		return fmt.Errorf("%w: dead letter read count cannot be negative", ErrInvalidBackendRequest)
	}

	return nil
}
