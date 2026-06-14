package backend

import (
	"fmt"
	"time"

	"github.com/Newton-School/goqueue/task"
)

// TaskResultRecord stores a task's terminal or latest execution result.
type TaskResultRecord struct {
	TaskID    task.TaskID
	Result    task.TaskResult
	UpdatedAt time.Time
	TTL       time.Duration
}

// Validate verifies the result record can be persisted.
func (r TaskResultRecord) Validate() error {
	if r.TaskID == "" {
		return fmt.Errorf("%w: task id is required", ErrInvalidBackendRequest)
	}
	if err := task.ValidateTaskID(r.TaskID.String()); err != nil {
		return err
	}
	if err := r.Result.Validate(); err != nil {
		return err
	}
	if r.TTL < 0 {
		return fmt.Errorf("%w: result ttl cannot be negative", ErrInvalidBackendRequest)
	}

	return nil
}
