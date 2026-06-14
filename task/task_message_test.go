package task

import "testing"

func TestTaskEnvelopeToMessageSerializesPayload(t *testing.T) {
	envelope, err := NewTaskEnvelope(TaskEnvelopeInput{
		Name:   "email.send",
		Queue:  "default",
		Args:   []any{"welcome"},
		Kwargs: map[string]any{"user_id": "u_123"},
	})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}

	message, err := TaskEnvelopeToMessage(envelope, JSONPayloadCodec{})
	if err != nil {
		t.Fatalf("TaskEnvelopeToMessage returned error: %v", err)
	}

	if message.ID != envelope.ID.String() {
		t.Fatalf("message ID = %q, want %q", message.ID, envelope.ID)
	}
	if len(message.Payload) == 0 {
		t.Fatal("message payload should be serialized")
	}
}

func TestTaskMessageToEnvelopeDecodesPayload(t *testing.T) {
	envelope, err := NewTaskEnvelope(TaskEnvelopeInput{
		Name:   "email.send",
		Queue:  "default",
		Args:   []any{"welcome"},
		Kwargs: map[string]any{"user_id": "u_123"},
	})
	if err != nil {
		t.Fatalf("NewTaskEnvelope returned error: %v", err)
	}

	message, err := TaskEnvelopeToMessage(envelope, JSONPayloadCodec{})
	if err != nil {
		t.Fatalf("TaskEnvelopeToMessage returned error: %v", err)
	}

	decoded, err := TaskMessageToEnvelope(message, JSONPayloadCodec{})
	if err != nil {
		t.Fatalf("TaskMessageToEnvelope returned error: %v", err)
	}

	if got := decoded.Payload.Args()[0]; got != "welcome" {
		t.Fatalf("decoded arg = %v, want welcome", got)
	}
	if got := decoded.Payload.Kwargs()["user_id"]; got != "u_123" {
		t.Fatalf("decoded kwarg = %v, want u_123", got)
	}
}
