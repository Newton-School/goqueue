package producer

import (
	"context"
	"fmt"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

// AsyncResult references a published task.
type AsyncResult struct {
	taskID  task.TaskID
	backend backend.QueueBackend
}

// ID returns the task id for this async result.
func (r *AsyncResult) ID() task.TaskID {
	if r == nil {
		return ""
	}

	return r.taskID
}

// TaskState returns latest task lifecycle state.
func (r *AsyncResult) TaskState(ctx context.Context) (backend.TaskStateRecord, error) {
	if r == nil {
		return backend.TaskStateRecord{}, fmt.Errorf("goqueue: async result is nil")
	}

	if r.backend == nil {
		return backend.TaskStateRecord{}, fmt.Errorf("goqueue: async result backend is nil")
	}

	if r.taskID == "" {
		return backend.TaskStateRecord{}, fmt.Errorf("goqueue: async result task id is empty")
	}

	return r.backend.GetTaskState(ctx, r.taskID)
}

// TaskResult returns latest task result.
func (r *AsyncResult) TaskResult(ctx context.Context) (backend.TaskResultRecord, error) {
	if r == nil {
		return backend.TaskResultRecord{}, fmt.Errorf("goqueue: async result is nil")
	}

	if r.backend == nil {
		return backend.TaskResultRecord{}, fmt.Errorf("goqueue: async result backend is nil")
	}

	if r.taskID == "" {
		return backend.TaskResultRecord{}, fmt.Errorf("goqueue: async result task id is empty")
	}

	return r.backend.GetTaskResult(ctx, r.taskID)
}

// ForgetTaskResult removes any stored result data for this task.
func (r *AsyncResult) ForgetTaskResult(ctx context.Context) error {
	if r == nil {
		return fmt.Errorf("goqueue: async result is nil")
	}

	if r.backend == nil {
		return fmt.Errorf("goqueue: async result backend is nil")
	}

	if r.taskID == "" {
		return fmt.Errorf("goqueue: async result task id is empty")
	}

	return r.backend.ForgetTaskResult(ctx, r.taskID)
}
