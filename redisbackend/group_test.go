package redisbackend

import (
	"context"
	"errors"
	"testing"

	"github.com/Newton-School/goqueue/backend"
)

func TestEnsureConsumerGroupRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	err := b.EnsureConsumerGroup(context.Background(), backend.ConsumerGroupRequest{Queue: "default", Group: "workers"})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("EnsureConsumerGroup error = %v, want ErrInvalidRedisOptions", err)
	}
}
