package inspect

import (
	"context"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

// TaskState fetches latest known task state information.
func (i *Inspector) TaskState(ctx context.Context, taskID task.TaskID) (backend.TaskStateRecord, error) {
	if i == nil {
		return backend.TaskStateRecord{}, ErrNilInspector
	}
	if i.backend == nil {
		return backend.TaskStateRecord{}, ErrInspectorBackend
	}
	if err := task.ValidateTaskID(taskID.String()); err != nil {
		return backend.TaskStateRecord{}, err
	}

	return i.backend.GetTaskState(ctx, taskID)
}
