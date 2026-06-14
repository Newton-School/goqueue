package redisbackend

import (
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestDeadLetterCodecRoundTripsRecord(t *testing.T) {
	record := backend.DeadLetterRecord{
		StreamID:       "2-0",
		Message:        testTaskMessage(t),
		Reason:         task.FailureRetryExhausted,
		Error:          "no attempts left",
		SourceStreamID: "1-0",
		Group:          "workers",
		Consumer:       "pod-1",
		FailedAt:       time.Date(2026, time.June, 14, 11, 0, 0, 0, time.UTC),
	}

	values, err := (deadLetterCodec{}).encode(record)
	if err != nil {
		t.Fatalf("encode returned error: %v", err)
	}

	decoded, err := (deadLetterCodec{}).decode("2-0", values)
	if err != nil {
		t.Fatalf("decode returned error: %v", err)
	}

	if decoded.StreamID != record.StreamID {
		t.Fatalf("stream id = %q, want %q", decoded.StreamID, record.StreamID)
	}
	if decoded.Message.ID != record.Message.ID {
		t.Fatalf("message id = %q, want %q", decoded.Message.ID, record.Message.ID)
	}
	if decoded.Reason != record.Reason {
		t.Fatalf("reason = %q, want %q", decoded.Reason, record.Reason)
	}
	if decoded.Error != record.Error {
		t.Fatalf("error = %q, want %q", decoded.Error, record.Error)
	}
	if !decoded.FailedAt.Equal(record.FailedAt) {
		t.Fatalf("failed at = %v, want %v", decoded.FailedAt, record.FailedAt)
	}
}

func TestDeadLetterCodecRejectsMissingMessage(t *testing.T) {
	_, err := (deadLetterCodec{}).decode("2-0", map[string]any{"reason": string(task.FailureExecution)})
	if !errors.Is(err, ErrInvalidRedisMessage) {
		t.Fatalf("decode error = %v, want ErrInvalidRedisMessage", err)
	}
}
