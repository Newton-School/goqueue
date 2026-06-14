package redisbackend

import (
	"context"
	"testing"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestIntegrationStateAndResultStorage(t *testing.T) {
	options := redisIntegrationOptions(t)
	ctx := context.Background()
	b, err := New(options)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	defer b.Close()
	defer cleanupIntegrationNamespace(ctx, t, b)

	taskID := task.TaskID("4ac0a01f-1b16-4330-b3e7-e99826eacb1a")
	if err := b.SetTaskState(ctx, backend.TaskStateRecord{TaskID: taskID, State: task.TaskStarted}); err != nil {
		t.Fatalf("SetTaskState returned error: %v", err)
	}
	state, err := b.GetTaskState(ctx, taskID)
	if err != nil {
		t.Fatalf("GetTaskState returned error: %v", err)
	}
	if state.State != task.TaskStarted {
		t.Fatalf("state = %s, want %s", state.State, task.TaskStarted)
	}

	if err := b.SaveTaskResult(ctx, backend.TaskResultRecord{TaskID: taskID, Result: task.SucceededResult("ok")}); err != nil {
		t.Fatalf("SaveTaskResult returned error: %v", err)
	}
	result, err := b.GetTaskResult(ctx, taskID)
	if err != nil {
		t.Fatalf("GetTaskResult returned error: %v", err)
	}
	if result.Result.State != task.TaskSucceeded {
		t.Fatalf("result state = %s, want %s", result.Result.State, task.TaskSucceeded)
	}
}
