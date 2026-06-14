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
