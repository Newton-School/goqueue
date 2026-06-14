package redisbackend

import (
	"context"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
)

func TestIntegrationClaimStaleReady(t *testing.T) {
	options := redisIntegrationOptions(t)
	ctx := context.Background()
	b, err := New(options)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	defer b.Close()
	defer cleanupIntegrationNamespace(ctx, t, b)

	if _, err := b.EnqueueReady(ctx, backend.EnqueueRequest{Message: testTaskMessage(t)}); err != nil {
		t.Fatalf("EnqueueReady returned error: %v", err)
	}
	if err := b.EnsureConsumerGroup(ctx, backend.ConsumerGroupRequest{Queue: "default", Group: "workers"}); err != nil {
		t.Fatalf("EnsureConsumerGroup returned error: %v", err)
	}
	if _, err := b.ReadReady(ctx, backend.ReadReadyRequest{
		Queue:    "default",
		Group:    "workers",
		Consumer: "pod-1",
		Count:    1,
		Block:    time.Second,
	}); err != nil {
		t.Fatalf("ReadReady returned error: %v", err)
	}

	claimed, err := b.ClaimStaleReady(ctx, backend.ClaimStaleReadyRequest{
		Queue:    "default",
		Group:    "workers",
		Consumer: "pod-2",
		MinIdle:  0,
		Count:    1,
		StartID:  "0-0",
	})
	if err != nil {
		t.Fatalf("ClaimStaleReady returned error: %v", err)
	}
	if len(claimed) != 1 {
		t.Fatalf("claimed = %d, want 1", len(claimed))
	}
}
