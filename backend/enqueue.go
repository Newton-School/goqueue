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
	if r.Message.Name == "" {
		return fmt.Errorf("%w: task message name is required", ErrInvalidBackendRequest)
	}
	if r.Message.Queue == "" {
		return fmt.Errorf("%w: task message queue is required", ErrInvalidBackendRequest)
	}

	return nil
}
