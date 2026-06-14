package backend

import (
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/task"
)

// ClaimStaleReadyRequest asks a backend to claim idle pending ready messages.
type ClaimStaleReadyRequest struct {
	Queue    task.QueueName
	Group    string
	Consumer string
	MinIdle  time.Duration
	Count    int64
	StartID  string
}

// Validate verifies pending-claim request fields.
func (r ClaimStaleReadyRequest) Validate() error {
	if err := task.ValidateQueueName(r.Queue.String()); err != nil {
		return err
	}
	if r.Group == "" {
		return fmt.Errorf("%w: consumer group is required", ErrInvalidBackendRequest)
	}
	if r.Consumer == "" {
		return fmt.Errorf("%w: consumer name is required", ErrInvalidBackendRequest)
	}
	if r.MinIdle < 0 {
		return fmt.Errorf("%w: min idle cannot be negative", ErrInvalidBackendRequest)
	}
	if r.Count < 0 {
		return fmt.Errorf("%w: claim count cannot be negative", ErrInvalidBackendRequest)
	}

	return nil
}
