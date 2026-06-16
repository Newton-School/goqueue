package inspect

import (
	"context"
	"time"

	"github.com/Newton-School/goqueue/task"
)

// TaskInspection aggregates state and result for a task id.
type TaskInspection struct {
	TaskID      task.TaskID `json:"task_id"`
	State       TaskState   `json:"state"`
	Result      TaskResult  `json:"result"`
	CheckedAt   time.Time   `json:"checked_at"`
	StateFound  bool        `json:"state_found"`
	ResultFound bool        `json:"result_found"`
}

// TaskState wraps backend state storage with normalized JSON-safe shape.
type TaskState struct {
	TaskID    task.TaskID    `json:"task_id"`
	State     task.TaskState `json:"state"`
	Error     string         `json:"error"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// TaskResult wraps backend result storage with normalized JSON-safe shape.
type TaskResult struct {
	TaskID    task.TaskID       `json:"task_id"`
	State     task.TaskState    `json:"state"`
	Value     any               `json:"value"`
	Error     string            `json:"error"`
	Metadata  map[string]string `json:"metadata"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// TaskSnapshot returns combined task state and result data.
func (i *Inspector) TaskSnapshot(ctx context.Context, taskID task.TaskID) (TaskInspection, error) {
	if i == nil {
		return TaskInspection{}, ErrNilInspector
	}

	snapshot := TaskInspection{
		TaskID:    taskID,
		CheckedAt: time.Now().UTC(),
	}

	state, stateErr := i.TaskState(ctx, taskID)
	if stateErr == nil {
		snapshot.StateFound = true
		snapshot.State = TaskState{
			TaskID:    state.TaskID,
			State:     state.State,
			Error:     state.Error,
			UpdatedAt: state.UpdatedAt,
		}
	}

	result, resultErr := i.TaskResult(ctx, taskID)
	if resultErr == nil {
		snapshot.ResultFound = true
		snapshot.Result = TaskResult{
			TaskID:    result.TaskID,
			State:     result.Result.State,
			Value:     result.Result.Value,
			Error:     result.Result.Error,
			Metadata:  result.Result.Metadata,
			UpdatedAt: result.UpdatedAt,
		}
	}

	if !snapshot.StateFound && !snapshot.ResultFound {
		if stateErr != nil {
			return TaskInspection{}, stateErr
		}
	}

	return snapshot, nil
}
