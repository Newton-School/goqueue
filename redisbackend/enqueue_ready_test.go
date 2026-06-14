package redisbackend

import (
	"context"
	"errors"
	"testing"

	"github.com/Newton-School/goqueue/backend"
)

func TestEnqueueReadyRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.EnqueueReady(context.Background(), backend.EnqueueRequest{Message: testTaskMessage(t)})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("EnqueueReady error = %v, want ErrInvalidRedisOptions", err)
	}
}
