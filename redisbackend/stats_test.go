package redisbackend

import (
	"context"
	"errors"
	"testing"

	"github.com/Newton-School/goqueue/backend"
)

func TestQueueStatsRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.QueueStats(context.Background(), backend.QueueStatsRequest{Queue: "default"})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("QueueStats error = %v, want ErrInvalidRedisOptions", err)
	}
}
