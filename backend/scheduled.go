package backend

import (
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/task"
)

// MoveDueScheduledRequest asks a backend to move due scheduled tasks to ready.
type MoveDueScheduledRequest struct {
	Queue task.QueueName
	Now   time.Time
	Limit int64
}

// MovedScheduledMessage is a scheduled task moved to a ready stream.
type MovedScheduledMessage struct {
	StreamID string
	Message  task.TaskMessage
}

// Validate verifies the scheduled move request.
func (r MoveDueScheduledRequest) Validate() error {
	if err := task.ValidateQueueName(r.Queue.String()); err != nil {
		return err
	}
	if r.Limit < 0 {
		return fmt.Errorf("%w: limit cannot be negative", ErrInvalidBackendRequest)
	}

	return nil
}
