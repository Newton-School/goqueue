package backend

import (
	"fmt"

	"github.com/Newton-School/goqueue/task"
)

// ConsumerGroupRequest asks a backend to ensure a queue consumer group exists.
type ConsumerGroupRequest struct {
	Queue task.QueueName
	Group string
}

// Validate verifies the consumer group request.
func (r ConsumerGroupRequest) Validate() error {
	if err := task.ValidateQueueName(r.Queue.String()); err != nil {
		return err
	}
	if r.Group == "" {
		return fmt.Errorf("%w: consumer group is required", ErrInvalidBackendRequest)
	}

	return nil
}
