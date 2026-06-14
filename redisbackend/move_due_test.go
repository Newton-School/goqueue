package redisbackend

import (
	"context"
	"errors"
	"testing"

	"github.com/Newton-School/goqueue/backend"
)

func TestMoveDueScheduledRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.MoveDueScheduled(context.Background(), backend.MoveDueScheduledRequest{Queue: "default"})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("MoveDueScheduled error = %v, want ErrInvalidRedisOptions", err)
	}
}
