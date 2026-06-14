package redisbackend

import (
	"context"
	"errors"
	"testing"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestSetTaskStateRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	err := b.SetTaskState(context.Background(), backend.TaskStateRecord{
		TaskID: "4ac0a01f-1b16-4330-b3e7-e99826eacb1a",
		State:  task.TaskStarted,
	})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("SetTaskState error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestGetTaskStateRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.GetTaskState(context.Background(), "4ac0a01f-1b16-4330-b3e7-e99826eacb1a")
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("GetTaskState error = %v, want ErrInvalidRedisOptions", err)
	}
}
