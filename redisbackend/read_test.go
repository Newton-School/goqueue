package redisbackend

import (
	"context"
	"errors"
	"testing"

	"github.com/Newton-School/goqueue/backend"
)

func TestReadReadyRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.ReadReady(context.Background(), backend.ReadReadyRequest{
		Queue:    "default",
		Group:    "workers",
		Consumer: "pod-1",
	})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("ReadReady error = %v, want ErrInvalidRedisOptions", err)
	}
}
