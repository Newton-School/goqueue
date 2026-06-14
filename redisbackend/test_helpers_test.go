package redisbackend

import (
	"testing"
	"time"

	"github.com/Newton-School/goqueue/task"
)

func testTaskMessage(t *testing.T) task.TaskMessage {
	t.Helper()

	envelope, err := task.NewTaskEnvelope(task.TaskEnvelopeInput{
		Name:      "email.send",
		Queue:     "default",
		Args:      []any{"welcome"},
		CreatedAt: time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}

	message, err := task.TaskEnvelopeToMessage(envelope, task.JSONPayloadCodec{})
	if err != nil {
		t.Fatalf("TaskEnvelopeToMessage returned error: %v", err)
	}

	return message
}
