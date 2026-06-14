package redisbackend

import (
	"context"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestIntegrationScheduledEnqueueAndMove(t *testing.T) {
	options := redisIntegrationOptions(t)
	ctx := context.Background()
	b, err := New(options)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	defer b.Close()
	defer cleanupIntegrationNamespace(ctx, t, b)

	envelope, err := task.NewTaskEnvelope(task.TaskEnvelopeInput{
		Name:   "email.send",
		Queue:  "default",
		Timing: task.TaskTiming{ETA: time.Now().UTC().Add(-time.Second)},
	})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}
	message, err := task.TaskEnvelopeToMessage(envelope, task.JSONPayloadCodec{})
	if err != nil {
		t.Fatalf("TaskEnvelopeToMessage returned error: %v", err)
	}

	if _, err := b.EnqueueScheduled(ctx, backend.EnqueueRequest{Message: message}); err != nil {
		t.Fatalf("EnqueueScheduled returned error: %v", err)
	}

	moved, err := b.MoveDueScheduled(ctx, backend.MoveDueScheduledRequest{Queue: "default", Limit: 10})
	if err != nil {
		t.Fatalf("MoveDueScheduled returned error: %v", err)
	}
	if len(moved) != 1 {
		t.Fatalf("len(moved) = %d, want 1", len(moved))
	}
	if moved[0].Message.ID != message.ID {
		t.Fatalf("moved ID = %q, want %q", moved[0].Message.ID, message.ID)
	}
}
