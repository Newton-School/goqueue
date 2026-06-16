package inspect

import (
	"context"

	"github.com/Newton-School/goqueue/task"
)

// ForgetTaskResult clears any persisted task result value.
func (i *Inspector) ForgetTaskResult(ctx context.Context, taskID task.TaskID) error {
	if i == nil {
		return ErrNilInspector
	}
	if i.backend == nil {
		return ErrInspectorBackend
	}
	if err := task.ValidateTaskID(taskID.String()); err != nil {
		return err
	}

	return i.backend.ForgetTaskResult(ctx, taskID)
}
