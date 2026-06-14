package goqueue

import (
	"context"
	"testing"
)

func TestTaskHandlerFuncImplementsHandler(t *testing.T) {
	handler := TaskHandlerFunc(func(ctx HandlerContext, payload TaskPayload) (TaskResult, error) {
		if ctx.Context() == nil {
			t.Fatal("handler context should not be nil")
		}
		if payload.Args()[0] != "welcome" {
			t.Fatalf("payload arg = %v, want welcome", payload.Args()[0])
		}
		return SucceededResult("sent"), nil
	})

	envelope, err := NewTaskEnvelope(TaskEnvelopeInput{Name: "email.send", Queue: "default"})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}

	result, err := handler.HandleTask(
		NewHandlerContext(context.Background(), envelope),
		NewTaskPayload([]any{"welcome"}, nil),
	)
	if err != nil {
		t.Fatalf("HandleTask returned error: %v", err)
	}
	if result.State != TaskSucceeded {
		t.Fatalf("State = %s, want %s", result.State, TaskSucceeded)
	}
}
