package backend

import (
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/task"
)

// TaskStateRecord stores the latest lifecycle state for a task.
type TaskStateRecord struct {
	TaskID    task.TaskID
	State     task.TaskState
	Error     string
	UpdatedAt time.Time
	TTL       time.Duration
}

// Validate verifies the state record can be persisted.
func (r TaskStateRecord) Validate() error {
	if r.TaskID == "" {
		return fmt.Errorf("%w: task id is required", ErrInvalidBackendRequest)
	}
	if err := task.ValidateTaskID(r.TaskID.String()); err != nil {
		return err
	}
	if err := task.ValidateTaskState(r.State); err != nil {
		return err
	}
	if r.TTL < 0 {
		return fmt.Errorf("%w: state ttl cannot be negative", ErrInvalidBackendRequest)
	}

	return nil
}
