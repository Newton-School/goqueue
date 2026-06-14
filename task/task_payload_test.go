package task

import "testing"

func TestNewTaskPayloadCopiesArgsAndKwargs(t *testing.T) {
	args := []any{"welcome", 10}
	kwargs := map[string]any{"user_id": "u_123"}

	payload := NewTaskPayload(args, kwargs)
	args[0] = "mutated"
	kwargs["user_id"] = "u_456"

	if got := payload.Args()[0]; got != "welcome" {
		t.Fatalf("Args()[0] = %v, want original value", got)
	}
	if got := payload.Kwargs()["user_id"]; got != "u_123" {
		t.Fatalf("Kwargs()[user_id] = %v, want original value", got)
	}
}

func TestTaskPayloadAccessorsReturnCopies(t *testing.T) {
	payload := NewTaskPayload([]any{"welcome"}, map[string]any{"user_id": "u_123"})

	args := payload.Args()
	args[0] = "mutated"

	kwargs := payload.Kwargs()
	kwargs["user_id"] = "u_456"

	if got := payload.Args()[0]; got != "welcome" {
		t.Fatalf("Args()[0] = %v, want original value", got)
	}
	if got := payload.Kwargs()["user_id"]; got != "u_123" {
		t.Fatalf("Kwargs()[user_id] = %v, want original value", got)
	}
}
