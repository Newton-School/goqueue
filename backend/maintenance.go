package backend

import (
	"fmt"

	"github.com/Newton-School/goqueue/task"
)

// PurgeQueueRequest asks a backend to delete queue-local artifacts.
type PurgeQueueRequest struct {
	Queue          task.QueueName
	DeleteMessages bool
	DeleteStates   bool
	DeleteResults  bool
}

// PurgeQueueResult reports what a purge operation removed.
type PurgeQueueResult struct {
	Queue            task.QueueName
	ReadyStream      int64
	ScheduledSet     int64
	DeadLetterStream int64
	TaskMessages     int64
	TaskStates       int64
	TaskResults      int64
}

// Validate verifies purge options.
func (r PurgeQueueRequest) Validate() error {
	if err := task.ValidateQueueName(r.Queue.String()); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidBackendRequest, err)
	}

	return nil
}
