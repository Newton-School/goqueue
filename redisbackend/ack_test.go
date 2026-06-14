package redisbackend

import (
	"context"
	"errors"
	"testing"

	"github.com/Newton-School/goqueue/backend"
)

func TestAckRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	err := b.Ack(context.Background(), backend.AckRequest{
		Queue:    "default",
		Group:    "workers",
		StreamID: "1-0",
	})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("Ack error = %v, want ErrInvalidRedisOptions", err)
	}
}
