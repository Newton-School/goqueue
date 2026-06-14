package redisbackend

import (
	"context"
	"testing"

	"github.com/Newton-School/goqueue/backend"
)

func TestIntegrationDeadLetterRoundTrip(t *testing.T) {
	options := redisIntegrationOptions(t)
	ctx := context.Background()
	b, err := New(options)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	defer b.Close()
	defer cleanupIntegrationNamespace(ctx, t, b)

	record, err := b.EnqueueDeadLetter(ctx, testDeadLetterRequest(t))
	if err != nil {
		t.Fatalf("EnqueueDeadLetter returned error: %v", err)
	}
	if record.StreamID == "" {
		t.Fatal("EnqueueDeadLetter returned empty stream id")
	}

	records, err := b.ReadDeadLetters(ctx, backend.ReadDeadLettersRequest{Queue: "default", Count: 10})
	if err != nil {
		t.Fatalf("ReadDeadLetters returned error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("dead letter records = %d, want 1", len(records))
	}
	if records[0].Message.ID != record.Message.ID {
		t.Fatalf("message id = %q, want %q", records[0].Message.ID, record.Message.ID)
	}
}

func TestIntegrationQueueStatsIncludesDeadLetters(t *testing.T) {
	options := redisIntegrationOptions(t)
	ctx := context.Background()
	b, err := New(options)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	defer b.Close()
	defer cleanupIntegrationNamespace(ctx, t, b)

	if _, err := b.EnqueueDeadLetter(ctx, testDeadLetterRequest(t)); err != nil {
		t.Fatalf("EnqueueDeadLetter returned error: %v", err)
	}

	stats, err := b.QueueStats(ctx, backend.QueueStatsRequest{Queue: "default"})
	if err != nil {
		t.Fatalf("QueueStats returned error: %v", err)
	}
	if stats.DeadLetterCount != 1 {
		t.Fatalf("DeadLetterCount = %d, want 1", stats.DeadLetterCount)
	}
}
