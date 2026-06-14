package redisbackend

import (
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestParseReadyStreamMessagesDecodesMessageField(t *testing.T) {
	message := testTaskMessage(t)
	encoded, err := (messageCodec{}).encode(message)
	if err != nil {
		t.Fatalf("encode returned error: %v", err)
	}

	ready, err := parseReadyStreamMessages([]redis.XStream{{
		Messages: []redis.XMessage{{
			ID:     "1-0",
			Values: map[string]any{"message": string(encoded)},
		}},
	}})
	if err != nil {
		t.Fatalf("parseReadyStreamMessages returned error: %v", err)
	}

	if len(ready) != 1 {
		t.Fatalf("len(ready) = %d, want 1", len(ready))
	}
	if ready[0].StreamID != "1-0" {
		t.Fatalf("StreamID = %q, want 1-0", ready[0].StreamID)
	}
}

func TestParseReadyMessagesDecodesClaimedMessages(t *testing.T) {
	message := testTaskMessage(t)
	encoded, err := (messageCodec{}).encode(message)
	if err != nil {
		t.Fatalf("encode returned error: %v", err)
	}

	ready, err := parseReadyMessages([]redis.XMessage{{
		ID:     "2-0",
		Values: map[string]any{"message": string(encoded)},
	}})
	if err != nil {
		t.Fatalf("parseReadyMessages returned error: %v", err)
	}
	if len(ready) != 1 {
		t.Fatalf("len(ready) = %d, want 1", len(ready))
	}
	if ready[0].StreamID != "2-0" {
		t.Fatalf("StreamID = %q, want 2-0", ready[0].StreamID)
	}
}
