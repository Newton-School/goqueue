package redisbackend

import (
	"errors"
	"testing"
)

func TestMessageCodecRoundTripsTaskMessage(t *testing.T) {
	codec := messageCodec{}
	message := testTaskMessage(t)

	encoded, err := codec.encode(message)
	if err != nil {
		t.Fatalf("encode returned error: %v", err)
	}

	decoded, err := codec.decode(encoded)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.ID != message.ID {
		t.Fatalf("decoded ID = %q, want %q", decoded.ID, message.ID)
	}
	if string(decoded.Payload) != string(message.Payload) {
		t.Fatalf("decoded payload = %q, want %q", decoded.Payload, message.Payload)
	}
}

func TestMessageCodecRejectsInvalidJSON(t *testing.T) {
	_, err := (messageCodec{}).decode([]byte("{"))
	if !errors.Is(err, ErrInvalidRedisMessage) {
		t.Fatalf("decode error = %v, want ErrInvalidRedisMessage", err)
	}
}

func TestMessageCodecRejectsDecodedMessageWithoutID(t *testing.T) {
	_, err := (messageCodec{}).decode([]byte(`{"name":"email.send","queue":"default"}`))
	if !errors.Is(err, ErrInvalidRedisMessage) {
		t.Fatalf("decode error = %v, want ErrInvalidRedisMessage", err)
	}
}
