package goqueue

import (
	"errors"
	"testing"
)

func TestJSONPayloadCodecRoundTripsPayload(t *testing.T) {
	codec := JSONPayloadCodec{}
	payload := NewTaskPayload([]any{"welcome"}, map[string]any{"user_id": "u_123"})

	encoded, err := codec.EncodePayload(payload)
	if err != nil {
		t.Fatalf("EncodePayload returned error: %v", err)
	}

	decoded, err := codec.DecodePayload(encoded)
	if err != nil {
		t.Fatalf("DecodePayload returned error: %v", err)
	}

	if got := decoded.Args()[0]; got != "welcome" {
		t.Fatalf("decoded arg = %v, want welcome", got)
	}
	if got := decoded.Kwargs()["user_id"]; got != "u_123" {
		t.Fatalf("decoded kwarg = %v, want u_123", got)
	}
}

func TestJSONPayloadCodecRejectsInvalidJSON(t *testing.T) {
	_, err := (JSONPayloadCodec{}).DecodePayload([]byte("{"))
	if !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("DecodePayload error = %v, want ErrInvalidPayload", err)
	}
}
