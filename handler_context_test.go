package goqueue

import (
	"context"
	"testing"
)

func TestNewHandlerContextExposesTaskDetails(t *testing.T) {
	envelope, err := NewTaskEnvelope(TaskEnvelopeInput{
		Name:     "email.send",
		Queue:    "default",
		Metadata: map[string]string{"trace_id": "trace-1"},
		Attempt:  2,
	})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}

	handlerContext := NewHandlerContext(context.Background(), envelope)

	if handlerContext.TaskID() != envelope.ID {
		t.Fatalf("TaskID = %q, want %q", handlerContext.TaskID(), envelope.ID)
	}
	if handlerContext.TaskName() != envelope.Name {
		t.Fatalf("TaskName = %q, want %q", handlerContext.TaskName(), envelope.Name)
	}
	if handlerContext.Queue() != envelope.Queue {
		t.Fatalf("Queue = %q, want %q", handlerContext.Queue(), envelope.Queue)
	}
	if handlerContext.Attempt() != 2 {
		t.Fatalf("Attempt = %d, want 2", handlerContext.Attempt())
	}
	if handlerContext.Metadata()["trace_id"] != "trace-1" {
		t.Fatalf("Metadata trace_id = %q, want trace-1", handlerContext.Metadata()["trace_id"])
	}
}

func TestNewHandlerContextUsesBackgroundForNilContext(t *testing.T) {
	envelope, err := NewTaskEnvelope(TaskEnvelopeInput{Name: "email.send", Queue: "default"})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}

	handlerContext := NewHandlerContext(nil, envelope)
	if handlerContext.Context() == nil {
		t.Fatal("Context should not be nil")
	}
}
