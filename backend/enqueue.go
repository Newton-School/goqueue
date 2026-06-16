package backend

import (
	"fmt"

	"github.com/Newton-School/goqueue/task"
)

// EnqueueRequest asks a backend to persist and enqueue a task message.
type EnqueueRequest struct {
	Message task.TaskMessage
}

// EnqueueResponse describes where the task was placed.
type EnqueueResponse struct {
	TaskID    task.TaskID
	StreamID  string
	Scheduled bool
}

// Validate verifies that the enqueue request contains a usable message.
func (r EnqueueRequest) Validate() error {
	if r.Message.ID == "" {
		return fmt.Errorf("%w: task message id is required", ErrInvalidBackendRequest)
	}
	if err := task.ValidateTaskID(r.Message.ID); err != nil {
		return err
	}
	if r.Message.Name == "" {
		return fmt.Errorf("%w: task message name is required", ErrInvalidBackendRequest)
	}
	if err := task.ValidateTaskName(r.Message.Name); err != nil {
		return err
	}
	if r.Message.Queue == "" {
		return fmt.Errorf("%w: task message queue is required", ErrInvalidBackendRequest)
	}
	if err := task.ValidateQueueName(r.Message.Queue); err != nil {
		return err
	}
	if err := task.ValidatePriority(r.Message.Priority); err != nil {
		return err
	}
	if err := r.Message.RetryPolicy.Validate(); err != nil {
		return err
	}
	if err := r.Message.Timing.Validate(); err != nil {
		return err
	}
	if r.Message.Attempt < 0 {
		return task.ErrInvalidTaskAttempt
	}

	return nil
}
