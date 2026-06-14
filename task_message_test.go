package goqueue

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
