package redisbackend

import (
	"context"
	"testing"
	"time"

	"github.com/Newton-School/goqueue/backend"
)

func TestIntegrationReadyEnqueueReadAck(t *testing.T) {
	options := redisIntegrationOptions(t)
	ctx := context.Background()
	b, err := New(options)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	defer b.Close()
	defer cleanupIntegrationNamespace(ctx, t, b)

	message := testTaskMessage(t)
	enqueueResponse, err := b.EnqueueReady(ctx, backend.EnqueueRequest{Message: message})
	if err != nil {
		t.Fatalf("EnqueueReady returned error: %v", err)
	}
	if enqueueResponse.StreamID == "" {
		t.Fatal("EnqueueReady returned empty stream ID")
	}

	if err := b.EnsureConsumerGroup(ctx, backend.ConsumerGroupRequest{Queue: "default", Group: "workers"}); err != nil {
		t.Fatalf("EnsureConsumerGroup returned error: %v", err)
	}

	ready, err := b.ReadReady(ctx, backend.ReadReadyRequest{
		Queue:    "default",
		Group:    "workers",
		Consumer: "pod-1",
		Count:    1,
		Block:    time.Second,
	})
	if err != nil {
		t.Fatalf("ReadReady returned error: %v", err)
	}
	if len(ready) != 1 {
		t.Fatalf("len(ready) = %d, want 1", len(ready))
	}

	if err := b.Ack(ctx, backend.AckRequest{Queue: "default", Group: "workers", StreamID: ready[0].StreamID}); err != nil {
		t.Fatalf("Ack returned error: %v", err)
	}
}
