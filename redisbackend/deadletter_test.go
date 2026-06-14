package redisbackend

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
	"github.com/Newton-School/goqueue/task"
)

func TestEnqueueDeadLetterRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.EnqueueDeadLetter(context.Background(), testDeadLetterRequest(t))
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("EnqueueDeadLetter error = %v, want ErrInvalidRedisOptions", err)
	}
}

func TestReadDeadLettersRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.ReadDeadLetters(context.Background(), backend.ReadDeadLettersRequest{Queue: "default"})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("ReadDeadLetters error = %v, want ErrInvalidRedisOptions", err)
	}
}

func testDeadLetterRequest(t *testing.T) backend.DeadLetterRequest {
	t.Helper()

	return backend.DeadLetterRequest{
		Message:        testTaskMessage(t),
		Reason:         task.FailureRetryExhausted,
		Error:          "no attempts left",
		SourceStreamID: "1-0",
		Group:          "workers",
		Consumer:       "pod-1",
		FailedAt:       time.Date(2026, time.June, 14, 12, 30, 0, 0, time.UTC),
	}
}
