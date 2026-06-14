package redisbackend

import (
	"testing"

	"github.com/Newton-School/goqueue/backend"
)

func TestParseMovedScheduledMessagesDecodesPairs(t *testing.T) {
	message := testTaskMessage(t)
	encoded, err := (messageCodec{}).encode(message)
	if err != nil {
		t.Fatalf("encode returned error: %v", err)
	}

	moved, err := parseMovedScheduledMessages([]any{"1-0", string(encoded)})
	if err != nil {
		t.Fatalf("parseMovedScheduledMessages returned error: %v", err)
	}

	if len(moved) != 1 {
		t.Fatalf("len(moved) = %d, want 1", len(moved))
	}
	if moved[0].StreamID != "1-0" {
		t.Fatalf("StreamID = %q, want 1-0", moved[0].StreamID)
	}
	if moved[0].Message.ID != message.ID {
		t.Fatalf("Message.ID = %q, want %q", moved[0].Message.ID, message.ID)
	}
}

func TestParseMovedScheduledMessagesReturnsBackendType(t *testing.T) {
	var _ []backend.MovedScheduledMessage
}
