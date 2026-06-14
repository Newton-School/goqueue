package redisbackend

import (
	"context"
	"errors"
	"testing"

	"github.com/Newton-School/goqueue/backend"
)

func TestEnqueueScheduledRejectsMissingETA(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}
	message := testTaskMessage(t)

	_, err := b.EnqueueScheduled(context.Background(), backend.EnqueueRequest{Message: message})
	if !errors.Is(err, ErrInvalidRedisMessage) {
		t.Fatalf("EnqueueScheduled error = %v, want ErrInvalidRedisMessage", err)
	}
}
