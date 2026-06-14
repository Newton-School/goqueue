package redisbackend

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
)

func TestClaimStaleReadyRejectsNilClient(t *testing.T) {
	b := &Backend{options: NewOptions("redis://localhost:6379/0"), keys: newKeyBuilder("goqueue")}

	_, err := b.ClaimStaleReady(context.Background(), backend.ClaimStaleReadyRequest{
		Queue:    "default",
		Group:    "workers",
		Consumer: "pod-2",
		MinIdle:  time.Minute,
	})
	if !errors.Is(err, ErrInvalidRedisOptions) {
		t.Fatalf("ClaimStaleReady error = %v, want ErrInvalidRedisOptions", err)
	}
}
