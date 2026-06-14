package redisbackend

import (
	"context"
	"errors"
	"testing"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestSaveTaskResultRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	err := b.SaveTaskResult(context.Background(), backend.TaskResultRecord{
		TaskID: "4ac0a01f-1b16-4330-b3e7-e99826eacb1a",
		Result: task.SucceededResult("ok"),
	})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("SaveTaskResult error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestGetTaskResultRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.GetTaskResult(context.Background(), "4ac0a01f-1b16-4330-b3e7-e99826eacb1a")
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("GetTaskResult error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestForgetTaskResultRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	err := b.ForgetTaskResult(context.Background(), "4ac0a01f-1b16-4330-b3e7-e99826eacb1a")
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("ForgetTaskResult error = %v, want ErrInvalidRedisOptions", err)
	}
}
