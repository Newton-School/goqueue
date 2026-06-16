package inspect

import (
	"context"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

// TaskResult fetches the latest persisted task result.
func (i *Inspector) TaskResult(ctx context.Context, taskID task.TaskID) (backend.TaskResultRecord, error) {
	if i == nil {
		return backend.TaskResultRecord{}, ErrNilInspector
	}
	if i.backend == nil {
		return backend.TaskResultRecord{}, ErrInspectorBackend
	}
	if err := task.ValidateTaskID(taskID.String()); err != nil {
		return backend.TaskResultRecord{}, err
	}

	return i.backend.GetTaskResult(ctx, taskID)
}
